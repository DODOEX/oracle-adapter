package main

import (
	"github.com/Dominator008/asset-price-oracle-adapter/adapter"
	"github.com/Dominator008/asset-price-oracle-adapter/bridge"
)

func main() {
	c := adapter.NewConfig()
	adapter.StartPairsTicker(c)

	bridge.NewServer(&adapter.AssetPrice{}).Start(c.Port)
}
