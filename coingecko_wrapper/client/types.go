package client

import (
	"sync"
)

type Client struct {
	Storage Storage
}

type Storage struct {
	Mu   sync.Mutex
	Data []Coin
}

type ApiErrResponse struct {
	Status struct {
		ErrorCode    uint16 `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}
}

type Coin struct {
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

type ClientResponse struct {
	Code uint   `json:"code"`
	Data []Coin `json:"data"`
}
