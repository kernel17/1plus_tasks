package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type ApiErrResponse struct {
	Status struct {
		ErrorCode    uint16 `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}
}

type ApiSuccessfulResponseEntry struct {
	Id                           string  `json:"id"`
	Symbol                       string  `json:"symbol"`
	Name                         string  `json:"name"`
	ImageUri                     string  `json:"image"`
	CurrentPrice                 float64 `json:"current_price"`
	MarketCap                    float64 `json:"market_cap"`
	MarketCapRank                float64 `json:"market_cap_rank"`
	FullyDilutedValuation        float64 `json:"fully_diluted_valuation"`
	TotalValue                   float64 `json:"total_volume"`
	High24H                      float64 `json:"high_24h"`
	Low24H                       float64 `json:"low_24h"`
	PriceChange24H               float64 `json:"price_change_24h"`
	PriceChangePercentage24H     float64 `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64 `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64 `json:"circulating_supply"`
	TotalSupply                  float64 `json:"total_supply"`
	MaxSupply                    float64 `json:"max_supply"`
	Ath                          float64 `json:"ath"`
	AthChangePercentage          float64 `json:"ath_change_percentage"`
	AthDate                      string  `json:"ath_date"`
	Atl                          float64 `json:"atl"`
	AtlChangePercentage          float64 `json:"atl_change_percentage"`
	AtlDate                      string  `json:"atl_date"`
	Roi                          struct {
		Times      float64 `json:"times"`
		Currency   string  `json:"currency"`
		Percentage float64 `json:"percentage"`
	} `json:"roi"`
	LastUpdated string `json:"last_updated"`
}

type Store struct {
	Mu   sync.Mutex
	Data []ApiSuccessfulResponseEntry
}

type ClientResponse struct {
	Code uint                         `json:"code"`
	Data []ApiSuccessfulResponseEntry `json:"data"`
}

func getCoins() (data []ApiSuccessfulResponseEntry, err error) {
	r, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1")
	if err != nil {
		log.Panic(err)
	}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var (
		d     []ApiSuccessfulResponseEntry
		error ApiErrResponse
	)
	err = json.Unmarshal(rBody, &d)
	if err != nil {
		err = json.Unmarshal(rBody, &error)
		if err != nil {
			return nil, err
		}
		if error.Status.ErrorCode == 429 {
			log.Println("Error 429")
			return nil, nil
		} else {
			log.Println("WARNING: ", error.Status.ErrorMessage)
			return nil, nil
		}
	}
	return d, nil

}

func main() {
	r := chi.NewRouter()
	var store Store
	go func() {
		for {
			log.Println("tick")
			data, err := getCoins()
			if err != nil {
				log.Fatalln(err)
			}
			if data != nil {
				store.Mu.Lock()
				store.Data = data
				store.Mu.Unlock()
				time.Sleep(10 * time.Minute)
			}

		}
	}()

	// handle coins sending
	r.Get("/get_{coin}", func(w http.ResponseWriter, r *http.Request) {
		req := chi.URLParam(r, "coin")
		var resp ClientResponse
		switch req {
		case "", "all":
			log.Println("123")
			store.Mu.Lock()
			if store.Data != nil {
				resp.Code = 0
				resp.Data = store.Data
				d, err := json.Marshal(resp)
				if err != nil {
					log.Fatalln(err)
				}
				w.Write(d)
			} else {
				resp.Code = 1
				resp.Data = nil
				d, err := json.Marshal(resp)
				if err != nil {
					log.Fatalln(err)
				}
				w.Write(d)
			}
			store.Mu.Unlock()
		default:
			store.Mu.Lock()
			if store.Data != nil {
				for _, elem := range store.Data {
					if elem.Symbol == req {
						resp.Code = 0
						resp.Data = append(resp.Data, elem)
						d, err := json.Marshal(resp)
						if err != nil {
							log.Fatalln(err)
						}
						w.Write(d)
						break
					}
				}
				resp.Code = 1
				resp.Data = nil
				d, err := json.Marshal(resp)
				if err != nil {
					log.Fatalln(err)
				}
				w.Write(d)

			} else {
				resp.Code = 1
				resp.Data = nil
				d, err := json.Marshal(resp)
				if err != nil {
					log.Fatalln(err)
				}
				w.Write(d)
			}
			store.Mu.Unlock()
		}

	})

	http.ListenAndServe(":7878", r)

}
