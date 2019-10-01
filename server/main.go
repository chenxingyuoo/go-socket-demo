package main

import (
	"net/http"
	"go-socket-demo/router"
)

func main() {
	r := router.SetupRouter()
	
	http.Handle("/", r)

	r.Run(":8000")
}