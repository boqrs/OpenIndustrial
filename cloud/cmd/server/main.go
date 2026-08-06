package main

import (
	"log"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/bootstrap"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/infrastructure/database"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/lifecycle"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/productinstance"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/trace"
	transport "github.com/OpenGongChang/OpenIndustrial/cloud/internal/transport/http"
)

func main() {
	db, err := database.NewPostgres(database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "cloud",
		Password: "cloud",
		DBName:   "industrial",
	})
	if err != nil {
		panic(err)
	}

	container := bootstrap.NewContainer(db)

	productAPI := productinstance.NewAPI(
		container.ProductInstance,
	)



	lifecycleAPI:=
	lifecycle.NewAPI(
		container.Lifecycle,
	)



	traceAPI:=
	trace.NewAPI(
		container.Trace,
	)



	router:=
	transport.NewRouter(
		productAPI,
		lifecycleAPI,
		traceAPI,
	)



	server:=
	transport.NewServer(
		router,
	)



	log.Println(
		"cloud server start :8080",
	)



	err = server.Start()


	if err!=nil{

		panic(err)

	}

}