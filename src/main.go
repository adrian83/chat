package main

import (
	"net/http"
	"os"
	"strconv"

	config "chatconfig"
	"db"
	"handler"
	"logger"
	chatsession "session"
	"ws"

	redisSession "github.com/adrian83/go-redis-session"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var (
	channels = ws.NewChannels()
)

func main() {

	// ---------------------------------------
	// application config
	// ---------------------------------------
	configPath := os.Args[1]
	logger.Infof("Main", "main", "Reading configuration from: %v", configPath)
	appConfig, err := config.ReadConfig(configPath)
	if err != nil {
		logger.Errorf("Main", "main", "Error while reading configuration! Error: %v", err)
		panic(err)
	}

	// ---------------------------------------
	// database
	// ---------------------------------------
	databaseConfig := &db.Config{
		Host:             appConfig.Database.Host,
		Port:             appConfig.Database.Port,
		DBName:           appConfig.Database.DBName,
		UsersTableName:   appConfig.Database.UsersTableName,
		UsersTablePKName: appConfig.Database.UsersTablePKName,
	}

	database, err := db.New(databaseConfig)
	if err != nil {
		logger.Errorf("Main", "main", "Error while creating RethinkDB session! Error: %v", err)
		panic(err)
	}
	defer func() {
		if err1 := database.Close(); err1 != nil {
			logger.Errorf("Main", "main", "Error while closing RethinkDB session! Error: %v", err1)
		}
	}()

	logger.Info("Main", "main", "RethinkDB session created")

	// ---------------------------------------
	// database setup
	// ---------------------------------------
	if err = database.Setup(); err != nil {
		logger.Errorf("Main", "main", "Error during RethinkDB database setup! Error: %v", err)
		panic(err)
	}
	logger.Info("Main", "main", "RethinkDB database initialized")

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
		logger.Errorf("Main", "main", "Error while creating SessionStore. Error: %v", err)
		return
	}
	defer func() {
		if err1 := sessionStore.Close(); err1 != nil {
			logger.Errorf("Main", "main", "Error while closing SessionStore session! Error: %v", err1)
		}
	}()
	logger.Info("Main", "main", "SessionStore created.")

	// ---------------------------------------
	// useful structures
	// ---------------------------------------
	simpleSession := chatsession.New(sessionStore)

	userRepository := db.NewUserRepository(database)

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

	mux.Handle("/talk", websocket.Handler(connect(simpleSession, channels)))

	// ---------------------------------------
	// http server
	// ---------------------------------------
	serverAddress := appConfig.Server.Host + ":" + strconv.Itoa(appConfig.Server.Port)
	logger.Infof("Main", "main", "Starting server on: %v", serverAddress)
	server := &http.Server{Addr: serverAddress, Handler: mux}
	if err2 := server.ListenAndServe(); err2 != nil {
		logger.Errorf("Main", "main", "Error while starting server! Error: %v", err2)
	}

}

func connect(simpleSession *chatsession.Session, channels ws.Channels) func(*websocket.Conn) {
	logger.Infof("Main", "Connect", "New connection")
	return func(wsc *websocket.Conn) {

		sessionID := chatsession.FindSessionID(wsc.Request())
		if sessionID == "" {
			logger.Errorf("Main", "Connect", "Error while getting sessionID from WebSocket.")
			return
		}

		user, err2 := simpleSession.FindUserData(sessionID)
		if err2 != nil {
			logger.Errorf("Main", "Connect", "Error while getting user data from session. Error: %v", err2)
			return
		}

		client := ws.NewClient(sessionID, user, ws.NewConnection(wsc), channels)

		if err := channels.AddClient(client); err != nil {
			logger.Infof("Main", "Connect", "Error while sending message: %v", err)
		}

		logger.Infof("Main", "Connect", "New connection received from %v, %v", client, user)

		client.Start()
	}
}
