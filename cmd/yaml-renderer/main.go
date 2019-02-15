package main

import (
	"log"
)

var (
	buildtime    string
	gitcommit    string
	appversion   string
)

func main() {
	log.Printf("appversion: %s", appversion)
	log.Printf("gitcommit:  %s", gitcommit)
	log.Printf("buildtime:  %s", buildtime)
}
