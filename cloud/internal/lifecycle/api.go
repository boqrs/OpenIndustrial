package lifecycle

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




type AppendRequest struct {


	ProductInstanceID string `json:"product_instance_id"`


	From string `json:"from"`


	To string `json:"to"`


	EventType string `json:"event_type"`


	Source string `json:"source"`

}





func (a *API) Append(
	w http.ResponseWriter,
	r *http.Request,
){


	var req AppendRequest



	json.NewDecoder(
		r.Body,
	).
	Decode(
		&req,
	)



	err:=a.service.AppendEvent(
		r.Context(),

		req.ProductInstanceID,

		req.From,

		req.To,

		req.EventType,

		req.Source,

		nil,
	)



	if err!=nil{


		http.Error(
			w,
			err.Error(),
			400,
		)


		return

	}



	w.WriteHeader(201)

}