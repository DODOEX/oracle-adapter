package connector

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKraken_SetPairs(t *testing.T) {
	kraken := Kraken{}
	_ = kraken.RefreshPairs()
	pairs := kraken.GetPairs()

	assert.Contains(t, pairs, &Pair{"ETH", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "EUR"})
}

func TestKraken_GetResponse(t *testing.T) {
	kraken := Kraken{}
	price, err := kraken.GetResponse("ETH", "EUR")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Kraken isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Kraken isn't greater than 0")
}
