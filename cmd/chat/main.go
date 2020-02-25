package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/adrian83/chat/pkg/config"
	"github.com/adrian83/chat/pkg/db"
	"github.com/adrian83/chat/pkg/exchange"
	"github.com/adrian83/chat/pkg/handler"
	"github.com/adrian83/chat/pkg/user"

	session "github.com/adrian83/go-redis-session"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

const (
	configPrefix = "chat"
)

func initLogger() {
	logger.SetFormatter(&logger.JSONFormatter{})
	logger.SetOutput(os.Stdout)
}

func initConfig() *config.Config {

	logger.Info("Reading configuration from environment variables")
	appConfig, err := config.ReadConfig(configPrefix)

	if err != nil {
		logger.Errorf("Error while reading configuration! Error: %v", err)
		panic(err)
	}

	return appConfig
}

func initRethink(config *config.Config) *db.RethinkDB {

	rethink := db.NewRethinkDB(config.DatabaseHost, config.DatabasePort, config.DatabaseName)

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

func initSession(config *config.Config) (session.Store, func()) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.SessionDbHost, config.SessionDbPort),
		Password: config.SessionDbPassword,
		DB:       config.SessionDbName,
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
	rethink := initRethink(appConfig)
	defer rethink.Close()

	// init session
	sessionStore, closeFnc := initSession(appConfig)
	defer closeFnc()

	// create chat rooms
	chatRooms := exchange.NewRooms()

	// ---------------------------------------
	// useful structures
	// ---------------------------------------

	userTable := rethink.GetUserTable()
	userService := user.NewUserService(userTable)

	templateRepository := handler.NewTemplateRepository(appConfig.StaticsPath)

	loginHandler := handler.NewLoginHandler(templateRepository, userService, sessionStore)
	logoutHandler := handler.NewLogoutHandler(templateRepository, sessionStore)
	registerHandler := handler.NewRegisterHandler(templateRepository, userService)
	indexHandler := handler.NewIndexHandler(templateRepository, sessionStore)
	conversationHandler := handler.NewConversationHandler(templateRepository, sessionStore)

	// ---------------------------------------
	// routing
	// ---------------------------------------
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(appConfig.StaticsPath))))

	router.HandleFunc("/", indexHandler.ShowIndexPage)

	router.HandleFunc("/login", loginHandler.ShowLoginPage).Methods("GET")
	router.HandleFunc("/login", loginHandler.LoginUser).Methods("POST")

	router.HandleFunc("/logout", logoutHandler.Logout).Methods("GET")

	router.HandleFunc("/register", registerHandler.ShowRegisterPage).Methods("GET")
	router.HandleFunc("/register", registerHandler.RegisterUser).Methods("POST")

	router.HandleFunc("/conversation", conversationHandler.ShowConversationPage).Methods("GET")

	router.Handle("/talk", websocket.Handler(connect(sessionStore, chatRooms)))

	// ---------------------------------------
	// http server
	// ---------------------------------------

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	serverAddress := appConfig.ServerHost + ":" + strconv.Itoa(appConfig.ServerPort)
	logger.Infof("Starting server on: %v", serverAddress)
	server := &http.Server{Addr: serverAddress, Handler: router}

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

func connect(sessionStore session.Store, chatRooms *exchange.Rooms) func(*websocket.Conn) {
	logger.Infof("New connection")

	return func(wsc *websocket.Conn) {
		sessionID, err := handler.ReadSessionIDFromCookie(wsc.Request())
		if err != nil {
			logger.Error(err)
			return
		}

		sess, err := sessionStore.Find(sessionID)
		if err != nil {
			logger.Errorf("Error while getting user session. Error: %v", err)
			return
		}

		var user user.User
		if err = sess.Get("user", &user); err != nil {
			logger.Errorf("Error while getting user data from session. Error: %v", err)
			return
		}

		router := exchange.NewRouter()

		wsConn := exchange.NewWebSocketConn(wsc)
		client := exchange.NewClient(sessionID, &user, chatRooms, wsConn, router)

		router.RegisterRoute(exchange.NewRoute(exchange.MsgUserJoinedRoomMT, exchange.NewAddClientToRoomHandler(chatRooms, client)))
		router.RegisterRoute(exchange.NewRoute(exchange.MsgTextMsgMT, exchange.NewSendMsgToRoomHandler(chatRooms)))
		router.RegisterRoute(exchange.NewRoute(exchange.MsgCreateRoomMT, exchange.NewCreateRoomHandler(chatRooms, client)))
		router.RegisterRoute(exchange.NewRoute(exchange.MsgUserLeftRoomMT, exchange.NewRemoveClientFromRoomHandler(chatRooms, client)))
		router.RegisterRoute(exchange.NewRoute(exchange.MsgLogoutMT, exchange.NewLogoutHandler(client)))

		chatRooms.AddClientToRoom(exchange.MainRoomName(), client)

		logger.Infof("New connection received from %v, %v", client, &user)

		client.Start()
	}
}
