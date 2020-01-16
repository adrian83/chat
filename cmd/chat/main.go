package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/adrian83/chat/pkg/client"
	"github.com/adrian83/chat/pkg/config"
	"github.com/adrian83/chat/pkg/connection"
	"github.com/adrian83/chat/pkg/db"
	"github.com/adrian83/chat/pkg/exchange"
	"github.com/adrian83/chat/pkg/handler"
	"github.com/adrian83/chat/pkg/rooms"
	"github.com/adrian83/chat/pkg/user"

	"github.com/adrian83/go-redis-session"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func initLogger() {
	logger.SetFormatter(&logger.JSONFormatter{})
	logger.SetOutput(os.Stdout)
}

func initConfig() *config.Config {
	configPath := os.Args[1]
	logger.Infof("Reading configuration from: %v", configPath)
	appConfig, err := config.ReadConfig(configPath)
	if err != nil {
		logger.Errorf("Error while reading configuration! Error: %v", err)
		panic(err)
	}
	return appConfig
}

func initRethink(dbConfig *config.DatabaseConfig) *db.RethinkDB {
	tables := map[string]string{dbConfig.UsersTableName: dbConfig.UsersTablePKName}

	rethink := db.NewRethinkDB(dbConfig.Host, dbConfig.Port, dbConfig.DBName, tables)

	if err := rethink.Connect(); err != nil {
		logger.Errorf("Error while creating RethinkDB session! Error: %v", err)
		panic(err)
	}

	if err := rethink.Setup(); err != nil {
		logger.Errorf("Error while creating RethinkDB database and/or tables! Error: %v", err)
		panic(err)
	}

	logger.Info("RethinkDB session created")

	return rethink
}

func initSession(sessionConfig *config.SessionConfig) (session.Store, func()) {

	options := &redis.Options{
		Addr:     fmt.Sprintf("%v:%v", sessionConfig.Host, sessionConfig.Port),
		Password: sessionConfig.Password,
		DB:       sessionConfig.DB,
	}

	client := redis.NewClient(options)

	sessionStore := session.NewStore(client, handler.SessionValidFor)

	closeFunc := func() {
		if err1 := sessionStore.Close(); err1 != nil {
			logger.Errorf("Error while closing SessionStore session! Error: %v", err1)
		}
	}
	logger.Info("SessionStore created.")
	return sessionStore, closeFunc
}

func main() {

	// initialize logger
	initLogger()

	// read config
	appConfig := initConfig()

	// init database
	dbConfig := &appConfig.Database
	rethink := initRethink(dbConfig)
	defer rethink.Close()

	// init session
	sessionConfig := &appConfig.Session
	sessionStore, close := initSession(sessionConfig)
	defer close()

	// create chat rooms
	chatRooms := rooms.NewRooms()

	// ---------------------------------------
	// useful structures
	// ---------------------------------------

	staticsConfig := appConfig.Statics

	userTable := rethink.GetTable(appConfig.Database.UsersTableName)
	userService := user.NewUserService(userTable)

	templateRepository := handler.NewTemplateRepository(staticsConfig)

	loginHandler := handler.NewLoginHandler(templateRepository, userService, sessionStore)
	logoutHandler := handler.NewLogoutHandler(templateRepository, sessionStore)
	registerHandler := handler.NewRegisterHandler(templateRepository, userService)
	indexHandler := handler.NewIndexHandler(templateRepository, sessionStore)
	conversationHandler := handler.NewConversationHandler(templateRepository, sessionStore)

	// ---------------------------------------
	// routing
	// ---------------------------------------
	mux := mux.NewRouter()
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticsConfig.Path))))

	mux.HandleFunc("/", indexHandler.ShowIndexPage)

	mux.HandleFunc("/login", loginHandler.ShowLoginPage).Methods("GET")
	mux.HandleFunc("/login", loginHandler.LoginUser).Methods("POST")

	mux.HandleFunc("/logout", logoutHandler.Logout).Methods("GET")

	mux.HandleFunc("/register", registerHandler.ShowRegisterPage).Methods("GET")
	mux.HandleFunc("/register", registerHandler.RegisterUser).Methods("POST")

	mux.HandleFunc("/conversation", conversationHandler.ShowConversationPage).Methods("GET")

	mux.Handle("/talk", websocket.Handler(connect(sessionStore, chatRooms)))

	// ---------------------------------------
	// http server
	// ---------------------------------------

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, os.Kill)

	serverAddress := appConfig.Server.Host + ":" + strconv.Itoa(appConfig.Server.Port)
	logger.Infof("Starting server on: %v", serverAddress)
	server := &http.Server{Addr: serverAddress, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Errorf("Server error! Error: %v", err)
		}
	}()

	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Warnf("Error while stopping server. Error: %v", err)
	}

	logger.Info("Server stopped.")
}

func connect(sessionStore session.Store, chatRooms *rooms.DefaultRooms) func(*websocket.Conn) {
	logger.Infof("New connection")
	return func(wsc *websocket.Conn) {

		sessionID, err := handler.ReadSessionIDFromCookie(wsc.Request())
		if err != nil {
			logger.Error(err)
			return
		}

		session, err := sessionStore.Find(sessionID)
		if err != nil {
			logger.Errorf("Error while getting user session. Error: %v", err)
			return
		}

		user := new(user.User)
		err = session.Get("user", user)
		if err != nil {
			logger.Errorf("Error while getting user data from session. Error: %v", err)
			return
		}

		wsConn := connection.NewWebSocketConn(wsc)
		client := client.NewClient(sessionID, user, chatRooms, wsConn)

		chatRooms.AddClientToRoom(exchange.MainRoomName(), client)

		logger.Infof("New connection received from %v, %v", client, user)

		client.Start()
	}
}
