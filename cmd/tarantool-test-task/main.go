package main

import (
	"log"

	"vk_tarantool_test_task/internal/pkg/app"
)

func main() {
	log.Print("Start application")
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
