package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func main() {

	// Construct a new API object and get env api key
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch user details on the account
	u, err := api.UserDetails(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Print user details
	fmt.Printf("ID: %v, EMAIL: %v\n", u.ID, u.Email)
	// GET all id domain
	z, err := api.ListZonesContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total domain: %v\n", z.Total)
	fmt.Printf("%#v\n", z.Result)
}
