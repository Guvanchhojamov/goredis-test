package main

import (
	"log"
	"redis-task/database"
	"redis-task/handler"
	"redis-task/server"
)

func main() {

	_, err := database.NewRedisDB()
	if err != nil {
		log.Fatalf("Error loading redis Client %v", err.Error())
	}
	_, err = database.NewPostgresDB()
	if err != nil {
		log.Fatalf("error connecting postgres DB %v", err.Error())
	}
	// run server
	srv := new(server.Server)
	handlers := new(handler.Handler)
	err = srv.Run(":8085", handlers.InitRoutes())
	if err != nil {
		log.Fatalf("server err: %s", err.Error())
		return
	}

}
