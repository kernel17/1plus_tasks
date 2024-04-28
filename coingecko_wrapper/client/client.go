package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func updateData(c *Client) {
	c.Storage.Mu.Lock()
	data, err := c.getCoins()
	if err != nil {
		log.Fatalln(err)
	}
	if data != nil {
		c.Storage.Data = data
	}
	c.Storage.Mu.Unlock()

}

func NewClient() *Client {
	var p Client

	updateData(&p)

	go func(c *Client) {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			log.Println("tick")
			updateData(&p)
		}
	}(&p)
	return &p
}

func (c *Client) getCoins() (data []Coin, err error) {
	r, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1")
	if err != nil {
		log.Panic(err)
	}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var (
		d     []Coin
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

func (c *Client) StartServer() {
	r := chi.NewRouter()
	r.Get("/get_{coin}", func(w http.ResponseWriter, r *http.Request) {
		req := chi.URLParam(r, "coin")
		var resp ClientResponse
		data := c.GetFromStorage(req)
		if data != nil {
			resp.Data = data
			d, err := json.Marshal(resp)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(d)
		} else {
			resp.Code, resp.Data = 1, nil
			d, err := json.Marshal(resp)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(d)

		}

	})

	http.ListenAndServe(":7878", r)

}

func (c *Client) GetFromStorage(req string) []Coin {
	c.Storage.Mu.Lock()
	defer c.Storage.Mu.Unlock()
	switch req {
	case "all":
		return c.Storage.Data
	default:
		if c.Storage.Data != nil {
			for _, elem := range c.Storage.Data {
				if elem.Symbol == req {
					return []Coin{elem}
				}
			}
		}
		return nil
	}

}
