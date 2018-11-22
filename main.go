package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/adrian83/chat/chat/config"
	"github.com/adrian83/chat/chat/db"
	"github.com/adrian83/chat/chat/handler"
	"github.com/adrian83/chat/chat/session"
	"github.com/adrian83/chat/chat/ws/client"
	"github.com/adrian83/chat/chat/ws/connection"
	"github.com/adrian83/chat/chat/ws/room"
	"github.com/adrian83/chat/chat/ws/rooms"

	redisSession "github.com/adrian83/go-redis-session"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var (
	chatRooms = rooms.NewRooms()
)

func main() {

	logger.SetFormatter(&logger.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// ---------------------------------------
	// application config
	// ---------------------------------------
	configPath := os.Args[1]
	logger.Infof("Reading configuration from: %v", configPath)
	appConfig, err := config.ReadConfig(configPath)
	if err != nil {
		logger.Errorf("Error while reading configuration! Error: %v", err)
		panic(err)
	}

	// ---------------------------------------
	// database
	// ---------------------------------------

	rethink := db.NewRethinkDB(appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.DBName, map[string]string{appConfig.Database.UsersTableName: appConfig.Database.UsersTablePKName})
	err = rethink.Connect()
	if err != nil {
		logger.Errorf("Error while creating RethingDB session! Error: %v", err)
		panic(err)
	}

	defer rethink.Close()

	err = rethink.Setup()
	if err != nil {
		logger.Errorf("Error while creating RethinkDB database and/or tables! Error: %v", err)
		panic(err)
	}

	

	logger.Info("RethinkDB session created")

	// ---------------------------------------
	// session
	// ---------------------------------------
	sessionStoreConfig := redisSession.Config{
		DB:       appConfig.Session.DB,
		Password: appConfig.Session.Password,
		Host:     appConfig.Session.Host,
		Port:     appConfig.Session.Port,
		IDLength: appConfig.Session.IDLength,
	}

	sessionStore, err := redisSession.NewStore(sessionStoreConfig)
	if err != nil {
		logger.Errorf("Error while creating SessionStore. Error: %v", err)
		return
	}
	defer func() {
		if err1 := sessionStore.Close(); err1 != nil {
			logger.Errorf("Error while closing SessionStore session! Error: %v", err1)
		}
	}()
	logger.Info("SessionStore created.")

	// ---------------------------------------
	// useful structures
	// ---------------------------------------
	simpleSession := session.New(sessionStore)

	userTable := rethink.GetTable(appConfig.Database.UsersTableName)
	userRepository := db.NewUserRepository(userTable)

	loginHandler := handler.NewLoginHandler(userRepository, simpleSession)
	logoutHandler := handler.NewLogoutHandler(simpleSession)
	registerHandler := handler.NewRegisterHandler(userRepository)
	indexHandler := handler.NewIndexHandler(simpleSession)
	conversationHandler := handler.NewConversationHandler(simpleSession)

	// ---------------------------------------
	// routing
	// ---------------------------------------
	mux := mux.NewRouter()
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	mux.HandleFunc("/", indexHandler.ShowIndexPage)

	mux.HandleFunc("/login", loginHandler.ShowLoginPage).Methods("GET")
	mux.HandleFunc("/login", loginHandler.LoginUser).Methods("POST")

	mux.HandleFunc("/logout", logoutHandler.Logout).Methods("GET")

	mux.HandleFunc("/register", registerHandler.ShowRegisterPage).Methods("GET")
	mux.HandleFunc("/register", registerHandler.RegisterUser).Methods("POST")

	mux.HandleFunc("/conversation", conversationHandler.ShowConversationPage).Methods("GET")

	mux.Handle("/talk", websocket.Handler(connect(simpleSession, chatRooms)))

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

	if err = server.Shutdown(ctx); err != nil {
		logger.Warnf("Error while stopping server. Error: %v", err)
	}

	logger.Info("Server stopped.")
}

func connect(simpleSession *session.Session, chatRooms *rooms.DefaultRooms) func(*websocket.Conn) {
	logger.Infof("New connection")
	return func(wsc *websocket.Conn) {

		sessionID := session.FindSessionID(wsc.Request())
		if sessionID == "" {
			logger.Errorf("Error while getting sessionID from WebSocket.")
			return
		}

		user, err := simpleSession.FindUserData(sessionID)
		if err != nil {
			logger.Errorf("Error while getting user data from session. Error: %v", err)
			return
		}

		wsConn := connection.NewWebSocketConn(wsc)
		client := client.NewClient(sessionID, &user, chatRooms, wsConn)

		chatRooms.AddClientToRoom(room.MainRoomName(), client)

		logger.Infof("New connection received from %v, %v", client, user)

		client.Start()
	}
}
