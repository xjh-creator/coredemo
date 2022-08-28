package main

import (
	"coredemo/framework"
	"net/http"
)

func main()  {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr: ":8080",
	}
	server.ListenAndServe()
}
