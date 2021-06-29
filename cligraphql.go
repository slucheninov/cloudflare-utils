package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty"
)

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

		j, err := json.Marshal(result.Data)
		if err != nil {
			return 0, err
		}

		if err := json.Unmarshal(j, &reeree); err != nil {
			return 0, err
		}

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
