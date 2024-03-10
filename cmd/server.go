package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/cyclex/planet-ban/pkg"
	"github.com/cyclex/planet-ban/repository/mongo"
	"github.com/cyclex/planet-ban/repository/postgre"
	"github.com/cyclex/planet-ban/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	_HttpDelivery "github.com/cyclex/planet-ban/delivery/http"
)

func run_server(server, config string, debug bool) (err error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[run_server] panic occurred")
		}
	}()

	// Create a context that can be cancelled when a shutdown signal is received
	c, cancel := context.WithCancel(context.Background())

	// Handle SIGINT (Ctrl+C) to initiate a graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v\n", sig)
		cancel() // Cancel the context to initiate graceful shutdown
	}()

	// load config
	cfg, err := pkg.LoadServiceConfig(config)
	if err != nil {
		err = errors.Wrap(err, "[run_server]")
		return
	}

	appLog = pkg.New("app", debug)
	authLog = pkg.New("authchatbot", debug)

	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbUser := cfg.Database.User
	dbPass := cfg.Database.Pass
	dbName := cfg.Database.Name
	dbSsl := cfg.Database.Ssl
	dbTimeout := cfg.Database.Timeout

	if dbSsl == "" {
		dbSsl = "disable"
	}

	if dbTimeout <= 0 {
		dbTimeout = 5
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s connect_timeout=%d", dbHost, dbPort, dbUser, dbName, dbPass, dbSsl, dbTimeout)
	fmt.Println(dsn)
	conn, err := ConnectDB("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	queueHost := cfg.Queue.Host
	queuePort := cfg.Queue.Port
	queueName := cfg.Queue.Name
	expired := cfg.Queue.Expired
	if expired < 1 {
		expired = 24
	}
	dsn = fmt.Sprintf("mongodb://%s:%d", queueHost, queuePort)
	queue, err := ConnectQueue(dsn, c)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = queue.Disconnect(c); err != nil {
			log.Fatal(err)
		}
	}()

	timeoutCtx := time.Duration(30) * time.Second
	namespace := cfg.Chatbot.Namespace
	parameterNamespace := cfg.Chatbot.ParameterNamespace
	urlSendMsg := cfg.Chatbot.Host

	ordersQueue := mongo.NewmongoRepository(c, queue.Database(queueName), queueName, time.Duration(expired))
	model := postgre.NewPostgreRepository(c, conn)
	ordersUcase := usecase.NewOrdersUcase(ordersQueue, timeoutCtx)
	chatUcase := usecase.NewChatUcase(model, urlSendMsg, "", namespace, parameterNamespace, ordersQueue)
	cmsUcase := usecase.NewCmsUcase(model, timeoutCtx, urlSendMsg, namespace, parameterNamespace)

	e := echo.New()
	_HttpDelivery.NewCmsHandler(e, cmsUcase)

	InitCron(ordersUcase, chatUcase, cmsUcase, timeoutCtx)

	RefreshToken(&processingAuth, chatUcase, c)

	go func() {
		if err := e.Start(server); err != nil {
			log.Fatalf("[run_server] %s", err)
		}
	}()

	// Wait for the context to be cancelled (e.g., by receiving SIGINT)
	<-c.Done()

	log.Println("Shutting down gracefully...")

	// Create a context with a timeout for shutdown
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// Shutdown the server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped.")

	return nil
}

func run_webhook(server, config string, debug bool) (err error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[run_webhook] panic occurred")
		}
	}()

	// Create a context that can be cancelled when a shutdown signal is received
	c, cancel := context.WithCancel(context.Background())

	// Handle SIGINT (Ctrl+C) to initiate a graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v\n", sig)
		cancel() // Cancel the context to initiate graceful shutdown
	}()

	// load config
	cfg, err := pkg.LoadServiceConfig(config)
	if err != nil {
		err = errors.Wrap(err, "[run_webhook]")
		return
	}

	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbUser := cfg.Database.User
	dbPass := cfg.Database.Pass
	dbName := cfg.Database.Name
	dbSsl := cfg.Database.Ssl
	dbTimeout := cfg.Database.Timeout

	if dbSsl == "" {
		dbSsl = "disable"
	}

	if dbTimeout <= 0 {
		dbTimeout = 5
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s connect_timeout=%d", dbHost, dbPort, dbUser, dbName, dbPass, dbSsl, dbTimeout)
	conn, err := ConnectDB("postgre", dsn)
	if err != nil {
		log.Fatal(err)
	}

	queueHost := cfg.Queue.Host
	queuePort := cfg.Queue.Port
	queueName := cfg.Queue.Name
	expired := cfg.Queue.Expired
	if expired < 1 {
		expired = 24
	}
	dsn = fmt.Sprintf("mongodb://%s:%d", queueHost, queuePort)
	queue, err := ConnectQueue(dsn, c)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = queue.Disconnect(c); err != nil {
			log.Fatal(err)
		}
	}()

	timeoutCtx := time.Duration(30) * time.Second
	namespace := cfg.Chatbot.Namespace
	parameterNamespace := cfg.Chatbot.ParameterNamespace
	urlSendMsg := cfg.Chatbot.Host

	ordersQueue := mongo.NewmongoRepository(c, queue.Database(queueName), queueName, time.Duration(expired))
	model := postgre.NewPostgreRepository(c, conn)
	ordersUcase := usecase.NewOrdersUcase(ordersQueue, timeoutCtx)
	chatUcase := usecase.NewChatUcase(model, urlSendMsg, "", namespace, parameterNamespace, ordersQueue)

	e := echo.New()
	_HttpDelivery.NewOrderHandler(e, ordersUcase, chatUcase, debug)

	go func() {
		if err := e.Start(server); err != nil {
			log.Fatalf("[run_webhook] %s", err)
		}
	}()

	// Wait for the context to be cancelled (e.g., by receiving SIGINT)
	<-c.Done()

	log.Println("Shutting down gracefully...")

	// Create a context with a timeout for shutdown
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// Shutdown the server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped.")

	return nil
}
