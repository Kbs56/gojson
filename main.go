package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// https://rapidapi.com/apidojo/api/realtor
// json result stuct from api call
type ResStruct struct {
	Data struct {
		HomeSearch struct {
			Results []struct {
				Location struct {
					Address struct {
						City       string `json:"city"`
						Line       string `json:"line"`
						PostalCode string `json:"postal_code"`
						StateCode  string `json:"state_code"`
					} `json:"address"`
				} `json:"location"`
				Description struct {
					Sqft  int `json:"sqft"`
					Beds  int `json:"beds"`
					Baths int `json:"baths"`
				} `json:"description"`
				Href               string `json:"href"`
				ListPrice          int    `json:"list_price"`
				PriceReducedAmount int    `json:"price_reduced_amount"`
				LastSoldPrice      int    `json:"last_sold_price"`
				ListDate           string `json:"list_date"`
				Status             string `json:"status"`
			} `json:"results"`
		} `json:"home_search"`
	} `json:"data"`
}

type ReqBody struct {
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	PostalCode string   `json:"postal_code"`
	Status     []string `json:"status"`
	SortFields struct {
		Direction string `json:"direction"`
		Field     string `json:"field"`
	} `json:"sort_fields"`
	ListPrice struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"list_price"`
}

func callApi(limit int, searchArea string, min int, max int) {
	url := "https://realtor.p.rapidapi.com/properties/v3/list"

	statusSlice := []string{"for_sale", "ready_to_build"}
	reqBody := &ReqBody{
		Limit:      limit,
		Offset:     0,
		PostalCode: searchArea,
		Status:     statusSlice,
		SortFields: struct {
			Direction string "json:\"direction\""
			Field     string "json:\"field\""
		}{Direction: "desc", Field: "list_date"},
		ListPrice: struct {
			Min int "json:\"min\""
			Max int "json:\"max\""
		}{Min: min, Max: max},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("realtorApiKey"))
	req.Header.Add("X-RapidAPI-Host", "realtor.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	resp, _ := io.ReadAll(res.Body)

	fmt.Println(string(resp))
	// createOutputFile(resp)
}

func createOutputFile(body []byte) {
	f, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}

	l, err := f.WriteString(string(body))
	if err != nil {
		f.Close()
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%d bytes written succesfully successful", l))
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func readJSONFile(fileName string) error {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	res := &ResStruct{}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	json.Unmarshal(byteValue, res)
	for i := 0; i < len(res.Data.HomeSearch.Results); i++ {
		fullStreet := res.Data.HomeSearch.Results[i].Location.Address.Line
		city := res.Data.HomeSearch.Results[i].Location.Address.City
		state := res.Data.HomeSearch.Results[i].Location.Address.StateCode
		zip := res.Data.HomeSearch.Results[i].Location.Address.PostalCode
		fmt.Println(fmt.Sprintf("%s %s, %s %s", fullStreet, city, state, zip))
		fmt.Println("Link:", res.Data.HomeSearch.Results[i].Href)
		fmt.Println("List Price:", res.Data.HomeSearch.Results[i].ListPrice)
		fmt.Println("List Date:", res.Data.HomeSearch.Results[i].ListDate)
		fmt.Println("Status:", res.Data.HomeSearch.Results[i].Status)
		fmt.Println("Price Reduced Amount:", res.Data.HomeSearch.Results[i].PriceReducedAmount)
		fmt.Println("Last Sold Price:", res.Data.HomeSearch.Results[i].LastSoldPrice)
		fmt.Println("Sqft:", res.Data.HomeSearch.Results[i].Description.Sqft)
		fmt.Println("Beds:", res.Data.HomeSearch.Results[i].Description.Beds)
		fmt.Println("Baths:", res.Data.HomeSearch.Results[i].Description.Baths)
		fmt.Println()
	}
	return nil
}

func main() {
	err := readJSONFile("output.json")
	if err != nil {
		panic(err)
	}
	// callApi(10, "75204", 0, 1_000_000)
}
