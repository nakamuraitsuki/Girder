package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
	"github.com/nakamuraitsuki/Girder/internal/api"
)

func main() {
	ctx := context.Background()

	libvirtClient, err := libvirt.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer libvirtClient.Close()

	ovnEndpoint := os.Getenv("GIRDER_OVN_NB_ENDPOINT")
	if ovnEndpoint == "" {
		ovnEndpoint = "tcp:127.0.0.1:6641" // Development default
	}

	ovnClient, err := ovn.Connect(ctx, ovnEndpoint)
	if err != nil {
		log.Fatalf("failed to connect to OVN NB DB: %v", err)
	}
	defer ovnClient.Close()

	router := api.NewRouter(libvirtClient, ovnClient)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
