package main

import "github.com/kleo-53/web_calc_go/internal/application"

func main() {
	app := application.New()
	app.RunServer()
}