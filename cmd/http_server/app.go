package main

import "github.com/gofiber/fiber/v2"

type App struct {
	*fiber.App
}

func newApp(httpSrv *fiber.App) *App {
	return &App{App: httpSrv}
}
