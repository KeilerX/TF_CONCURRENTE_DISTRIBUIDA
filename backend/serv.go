package main

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Persona struct {
	Id       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Edad     int    `json:"edad"`
}

func RootEndpointGET(response http.ResponseWriter, request *http.Request) {
	persJson := `[{"id":1,"nombre": "Emmanuel","apellido": "German","edad": 22},{"id":2,"nombre":"Juan Diego","apellido":"Peña","edad":50}]`
	var pers []Persona
	json.Unmarshal([]byte(persJson), &pers)
	data, _ := json.Marshal(pers)
	for _, per := range pers {
		fmt.Println("Nombre: ", per.Nombre)
		fmt.Println("Apellido: ", per.Apellido)
	}
	//fmt.Println(pers)
	fmt.Println("JSON: ", string(data))
	response.Write(data)
}
func RootEndpointPOST(response http.ResponseWriter, request *http.Request) {
	persJson := `[{"id":1,"nombre": "Emmanuel","apellido": "German","edad": 22},{"id":2,"nombre":"Juan Diego","apellido":"Peña","edad":50}]`
	var pers []Persona
	json.Unmarshal([]byte(persJson), &pers)
	data, _ := json.Marshal(pers)
	for _, per := range pers {
		fmt.Println("Nombre: ", per.Nombre)
		fmt.Println("Apellido: ", per.Apellido)
	}
	//fmt.Println(pers)
	fmt.Println("JSON: ", string(data))
	response.Write(data)
}

func main() {
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/", RootEndpointGET).Methods("GET")
	router.HandleFunc("/knn", RootEndpointPOST).Methods("POST")
	fmt.Println("Now server is running on port 8000")
	http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(router))
}
