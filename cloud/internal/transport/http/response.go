package http

import (
	"encoding/json"
	"net/http"
)



type Response struct {


	Code int `json:"code"`


	Message string `json:"message"`


	Data any `json:"data"`

}



func JSON(
	w http.ResponseWriter,
	code int,
	data any,
){


	w.Header().
	Set(
		"Content-Type",
		"application/json",
	)



	json.NewEncoder(w).
	Encode(
		Response{


			Code:code,


			Message:"success",


			Data:data,

		},
	)

}