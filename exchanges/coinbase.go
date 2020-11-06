package connector

import (
	"fmt"
	"strconv"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

type Coinbase struct {
	Exchange
}

func (exc *Coinbase) GetResponse(base, quote string) (*Response, error) {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*coinbasepro.Client)

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, &Error{exc.GetConfig().Name, "500 ERROR", err.Error()}
	}

	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return nil, &Error{exc.GetConfig().Name, "Parse price error", err.Error()}
	}
	volume, err := strconv.ParseFloat(string(ticker.Volume), 64)
	if err != nil {
		return nil, &Error{exc.GetConfig().Name, "Parse volume error", err.Error()}
	}
	return &Response{exc.GetConfig().Name, price, volume * price}, nil
}

func (exc *Coinbase) RefreshPairs() error {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*coinbasepro.Client)

	products, err := client.GetProducts()
	if err != nil {
		return &Error{Exchange: exc.GetConfig().Name, Message: err.Error()}
	}

	var pairs []*Pair
	for _, product := range products {
		pairs = append(pairs, &Pair{product.BaseCurrency, product.QuoteCurrency})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Coinbase) GetConfig() *Config {
	return &Config{Name: "Coinbase", Client: coinbasepro.NewClient()}
}
