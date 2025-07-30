package main

import (
	"archiveFiles/config"
	"archiveFiles/internal/app"
	"fmt"
	"log"
)

func main() {
	fmt.Println("start app")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
