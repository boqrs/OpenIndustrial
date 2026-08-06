package trace

import (
	"encoding/json"

	"net/http"
)



type API struct {


	service *Service

}



func NewAPI(
	service *Service,
)*API{


	return &API{
		service:service,
	}

}



func (a *API) Query(
	w http.ResponseWriter,
	r *http.Request,
){


	id:=r.URL.Query().
	Get(
		"id",
	)



	result,err:=
	a.service.Query(
		r.Context(),
		id,
	)



	if err!=nil{

		http.Error(
			w,
			err.Error(),
			500,
		)

		return

	}



	json.NewEncoder(w).
	Encode(result)

}