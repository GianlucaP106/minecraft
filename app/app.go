package app

import (
	"log"
)

func Start() {
	log.Println("Starting game...")
	g := Game{}
	g.Init()
	g.Run()
}
