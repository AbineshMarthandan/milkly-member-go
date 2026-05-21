package di

import (
	"milkly-member/config"
	"milkly-member/controller"
	"milkly-member/repository"
	"milkly-member/service"

	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()

	// Register config
	Container.Provide(config.NewConfig)

	// Register repositories
	Container.Provide(repository.NewMemberRepository)

	// Register services
	Container.Provide(service.NewMemberService)

	// Register controllers
	Container.Provide(controller.NewMemberController)
	Container.Provide(controller.NewControllers)
}