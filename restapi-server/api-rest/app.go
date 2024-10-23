package main

import (
	"fabricrest-go/api-rest/web"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	//Fluidos
	router.HandleFunc("/xadatu/auth/register", web.RegisterAuthReq).Methods("POST")
	router.HandleFunc("/xadatu/auth/queryByDate", web.QueryAuthReqByDate).Methods("GET")
	router.HandleFunc("/xadatu/auth/{id}", web.QueryAuthReq).Methods("GET")

	//XACML
	router.HandleFunc("/ngsi-ld/v1/entities/", web.QueryAllXACMLEntities).Methods("GET")
	router.HandleFunc("/ngsi-ld/v1/entities/", web.RegisterXACMLEntities).Methods("POST")
	router.HandleFunc("/ngsi-ld/v1/entities/{entity}/attrs", web.UpdateXACMLEntity).Methods("PATCH")
	router.HandleFunc("/ngsi-ld/v1/entities/{entity}", web.QueryXACMLEntity).Methods("GET")

	fmt.Println("Listening (http://localhost:3002/)...")
	if err := http.ListenAndServe(":3002", router); err != nil {
		fmt.Println(err)
	}
}
