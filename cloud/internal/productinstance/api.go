package productinstance

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



func (a *API) Create(
	w http.ResponseWriter,
	r *http.Request,
){



	var req CreateRequest



	err:=json.NewDecoder(
		r.Body,
	).Decode(
		&req,
	)



	if err!=nil{

		http.Error(
			w,
			err.Error(),
			400,
		)

		return

	}



	ctx:=r.Context()



	instance,err:=
	a.service.Create(
		ctx,
		req,
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
	Encode(instance)

}