package main

import (
	"fabricrest-go/api-rest/web"
	"fmt"
	"net/http"
)

func main() {
	//Fluidos
	http.HandleFunc("/xadatu/auth/register", web.RegisterAuthReq)
	http.HandleFunc("/xadatu/auth", web.QueryAuthReq)
	http.HandleFunc("/xadatu/auth/queryByDate", web.QueryAuthReqByDate)

	//XACML
	http.HandleFunc("/ngsi-ld/v1/entities", web.ManageEntities)

	fmt.Println("Listening (http://localhost:3002/)...")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		fmt.Println(err)
	}
}
