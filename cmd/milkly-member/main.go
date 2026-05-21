package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"milkly-member/controller"
	"milkly-member/di"
	"milkly-member/config"

	_ "milkly-member/docs"
)

//	@title			Milkly Member Service API
//	@version		1.0
//	@description	API for Milkly Member Management
//	@basePath		/api/v1

//	@contact.name	Milkly Team
func main() {
	err := di.Container.Invoke(func(cfg *config.Config, controllers *controller.Controllers) error {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		// Start HTTP server with controllers
		go controllers.Listen()

		log.Println("Milkly Member Service started successfully")

		<-ch
		log.Println("Shutting down application and releasing resources")

		return nil
	})

	if err != nil {
		log.Panicf("Application failed to start: %s", err)
	}
}