package main

import (
	"context"
	"log"
	"encoding/json"
	"errors"
	"strings"
	"net/http"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)


// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

type uppercaseRequest struct{
	S string
}
type uppercaseReponse struct {
	V string
	Err string
}

type StringService interface {
	Uppercase(string) (string, error)
}

type stringService struct{}
func (stringService) Uppercase(s string) (string, error) {
	if s == ""{
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func makeUppercaseEndpoint(svc StringService)  endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req :=request.(uppercaseRequest)
		v, err :=svc.Uppercase(req.S)

		if err != nil {
			return uppercaseReponse{v, err.Error()}, nil	}
			return uppercaseReponse{v, ""}, nil}
	}

func decodeUppercaseRequest(_ context.Context, r *http.Request)  (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err !=nil{
		return nil, err
	}
	return request, nil
}

func encodeReponse(_ context.Context, w http.ResponseWriter, reponse interface{}) error{
	return json.NewEncoder(w).Encode(reponse)
}

func  main() {
	svc :=stringService{}
	endPoint := makeUppercaseEndpoint(svc)
	uppercaseHandler  :=  httptransport.NewServer(endPoint, decodeUppercaseRequest, encodeReponse,)
	http.Handle("/uppercase", uppercaseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
