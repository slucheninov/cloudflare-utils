package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/go-resty/resty"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	startDate, stopDate string
)

func init() {

	// date UTC now
	t := time.Now().UTC()
	year, month, _ := t.Date()
	firstDayOfThisMonth := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())

	flag.StringVar(&startDate, "startDate", firstDayOfThisMonth.Format("2006-01-02"), "start date for count")
	flag.StringVar(&stopDate, "stopDate", t.Format("2006-01-02"), "stop date for count")
	flag.Parse()

	if !(len(os.Getenv("CF_API_TOKEN")) > 0) {
		log.Fatal("Please provide CF_API_TOKEN")
	}
}

func trafficDomainId(token string, zoneTag string, dategeq string, datelt string) (summa int64, err error) {

	type Response struct {
		Viewer struct {
			Zones []struct {
				HTTPRequests1DGroups []struct {
					Sum struct {
						Bytes int64 `json:"bytes"`
					} `json:"sum"`
				} `json:"httpRequests1dGroups"`
			} `json:"zones"`
		} `json:"viewer"`
	}

	type graphErr struct {
		Message string
	}

	type graphResponse struct {
		Data   interface{}
		Errors []graphErr
	}
	var reeree Response
	//var summa int64

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(token).
		SetResult(&graphResponse{}).
		SetBody(fmt.Sprintf(`{ "query":"query {viewer {zones(filter: {zoneTag: \"%s\"}) {httpRequests1dGroups(limit: 10, filter: {date_geq: \"%s\", date_lt: \"%s\"}) {sum{bytes }}}}}"}`, zoneTag, dategeq, datelt)).
		Post("https://api.cloudflare.com/client/v4/graphql")

	if err != nil {
		return 0, err
	}

	if resp.StatusCode() == http.StatusOK {
		result := resp.Result().(*graphResponse)
		// Convert to json graphql data: ....
		j, err := json.Marshal(result.Data)
		if err != nil {
			return 0, err
		}
		// parses the JSON-encoded data Viewer
		if err := json.Unmarshal(j, &reeree); err != nil {
			return 0, err
		}
		// to access HTTPRequests1DGroups
		for _, z := range reeree.Viewer.Zones {
			if len(z.HTTPRequests1DGroups) == 0 {
				return 0, nil
			}
			zt := z.HTTPRequests1DGroups[0]
			summa = zt.Sum.Bytes
		}
		return summa, nil
	} else {
		return 0, fmt.Errorf("no http status 200")
	}
}

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
	//fmt.Printf("%#v\n", z)
	for _, zi := range z {
		if zi.Status == "active" && zi.Plan.IsSubscribed {
			idzone[zi.ID] = zi.Name
		}
	}
	// Count bandwidth all domain
	for id, domain := range idzone {
		fmt.Printf("Domain: %v %v ", domain, id)
		s, err := trafficDomainId(os.Getenv("CF_API_TOKEN"), id, startDate, stopDate)
		if err != nil {
			log.Fatal(err)
		}
		summa = summa + s
		fmt.Printf("%v \n", ByteCountIEC(s))
	}
	fmt.Printf("Summa: %v \n", ByteCountIEC(summa))
}
