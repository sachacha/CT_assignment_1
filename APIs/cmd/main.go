package main

import (
	"APIs"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/conservation/v1/country/", APIs.HandlerCountry) // ensure to type complete URL when requesting
	http.HandleFunc("/conservation/v1/species/", APIs.HandlerSpecies)
	http.HandleFunc("/conservation/v1/diag/", APIs.HandlerDiag)
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
