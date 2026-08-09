package main

import (
	"log"
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"github.com/nakamuraitsuki/Girder/internal/api"
)

func main() {
	libvirtClient, err := libvirt.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer libvirtClient.Close()

	router := api.NewRouter(libvirtClient)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
