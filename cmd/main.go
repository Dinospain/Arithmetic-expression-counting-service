package main

import (
	"github.com/Dinospain/Arithmetic-expression-counting-service/internal/application"
)

func main() {
	app := application.New()
	// app.Run()
	app.RunServer()
}
