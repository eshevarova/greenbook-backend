package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

type application struct {
	servicePort      int
	router           *fasthttprouter.Router
}

func NewApplication() *application {
	return &application{
		servicePort:      7888,
		router:           fasthttprouter.New(),
	}
}

func (app *application) Start() {
	requestHandler := createRequestHandler()
	httpServer := fasthttp.Server{
		Name:            "greenbook_api",
		Handler: requestHandler,
	}

	log.Println("Http api started :", app.servicePort)
	log.Fatal(httpServer.ListenAndServe(":" + strconv.Itoa(app.servicePort)))
}
