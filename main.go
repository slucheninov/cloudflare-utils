package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// bytes human readable (stackoverflow)
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func main() {
	// Init vars
	var idzone = make(map[string]string)
	// date UTC now
	t := time.Now().UTC()
	year, month, _ := t.Date()
	firstDayOfThisMonth := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	// Init all sum
	var summa int64
	// Construct a new API object and get env api key
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CF_API_TOKEN"))
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
	z, err := api.ListZones(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%#v\n", z)
	for _, zi := range z {
		if zi.Status == "active" {
			idzone[zi.ID] = zi.Name
		}
	}
	// Count bandwidth all domain
	for id, domain := range idzone {
		fmt.Printf("Domain: %v %v ", domain, id)
		s, err := trafficDomainId(os.Getenv("CF_API_TOKEN"), id, firstDayOfThisMonth.Format("2006-01-02"), t.Format("2006-01-02"))
		if err != nil {
			log.Fatal(err)
		}
		summa = summa + s
		fmt.Printf("%v \n", ByteCountIEC(s))
	}
	fmt.Printf("Summa: %v \n", ByteCountIEC(summa))
}
