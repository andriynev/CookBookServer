package main

import (
	"fmt"
	"food/src/api/config"
	"food/src/api/database"
	"food/src/api/handler"
	"food/src/api/jwt_auth"
	"food/src/api/server"
	"os"
	"os/signal"
)

func main() {

	cfg := config.GetConfig()
	_, err := database.InitDB(cfg.DSN)
	if err != nil {
		panic(err)
	}
	jwt_auth.Setup(cfg.JwtKey)

	addr := fmt.Sprintf(":%s", cfg.Port)
	httpHandler := handler.SetupHandler()
	srv := server.NewServer(addr, httpHandler)
	srv.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	srv.Stop()
}
