package main

import (
	"bank-chat/cmd/bank-history-cards/app"
	"bank-chat/pkg/core/chat"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jafarsirojov/mux/pkg/mux"
	"log"
	"net"
	"net/http"
)

var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)

//-host 0.0.0.0 -port 9013 -dsn postgres://user:pass@localhost:5509/app
const ENV_PORT = "PORT"
const ENV_DSN = "DATABASE_URL"
const ENV_HOST = "HOST"

func main() {
	flag.Parse()
	envPort, ok := FromFlagOrEnv(*port, ENV_PORT)
	if !ok {
		log.Println("can't port")
		return
	}
	envDsn, ok := FromFlagOrEnv(*dsn, ENV_DSN)
	if !ok {
		log.Println("can't dsn")
		return
	}
	envHost, ok := FromFlagOrEnv(*host, ENV_HOST)
	if !ok {
		log.Println("can't host")
		return
	}
	addr := net.JoinHostPort(envHost, envPort)
	log.Println("starting server!")
	log.Printf("host = %s, port = %s\n", envHost, envPort)

	pool, err := pgxpool.Connect(
		context.Background(),
		envDsn,
	)
	if err != nil {
		panic(err)
	}
	usersSvc := chat.NewService(pool)
	usersSvc.Start()
	exactMux := mux.NewExactMux()
	server := app.NewMainServer(exactMux, usersSvc)
	exactMux.GET("/api/chat/message/all",
		server.HandleGetMessageAll,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	exactMux.GET("/api/chat/message/{recipient_id}",
		server.HandleGetMessaegeByRecipientID,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	exactMux.POST("/api/chat/message",
		server.HandlePostAddMassage,
		jwtMiddleware,
		requestIdier,
		logger,
	)
	panic(http.ListenAndServe(addr, server))
}
