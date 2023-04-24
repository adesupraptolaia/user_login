package main

import (
	"flag"
	"log"
	"os"

	"github.com/adesupraptolaia/user_login/cmd/auth"
	"github.com/adesupraptolaia/user_login/cmd/user"
)

func main() {
	flag.Parse()

	app := flag.Arg(0)
	if os.Getenv("APP_NAME") != "" {
		app = os.Getenv("APP_NAME")
	}

	if app == "auth" {
		auth.Run()
	} else if app == "user" {
		user.Run()
	} else {
		log.Fatalln("insert app argument, auth or user \n ex: go run main.go user \n NOT ", app)
	}
}
