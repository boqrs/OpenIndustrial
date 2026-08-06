package http

import (
	nethttp "net/http"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/productinstance"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/lifecycle"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/trace"
)



func NewRouter(
	productAPI *productinstance.API,
	lifecycleAPI *lifecycle.API,
	traceAPI *trace.API,
) *nethttp.ServeMux {



	mux:=nethttp.NewServeMux()



	mux.HandleFunc(
		"/api/product-instances",
		productAPI.Create,
	)



	mux.HandleFunc(
		"/api/lifecycle/events",
		lifecycleAPI.Append,
	)



	mux.HandleFunc(
		"/api/trace",
		traceAPI.Query,
	)



	return mux

}