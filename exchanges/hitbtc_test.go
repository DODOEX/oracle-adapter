package connector

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHitBtc_SetPairs(t *testing.T) {
	hitBtc := HitBtc{}
	_ = hitBtc.RefreshPairs()
	pairs := hitBtc.GetPairs()

	assert.Contains(t, pairs, &Pair{"ETH", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "BTC"})
}

func TestHitBtc_GetResponse(t *testing.T) {
	hitBtc := HitBtc{}
	price, err := hitBtc.GetResponse("ETH", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from HitBTC isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from HitBTC isn't greater than 0")
}
