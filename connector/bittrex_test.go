package connector

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBittrex_SetPairs(t *testing.T) {
	bittrex := Bittrex{}
	_ = bittrex.RefreshPairs()
	pairs := bittrex.GetPairs()

	assert.Contains(t, pairs, &Pair{"BTC", "ETH"})
	assert.Contains(t, pairs, &Pair{"BTC", "LTC"})
}

func TestBittrex_GetResponse(t *testing.T) {
	bittrex := Bittrex{}
	price, err := bittrex.GetResponse("BTC", "ETH")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Bittrex isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Bittrex isn't greater than 0")
}
