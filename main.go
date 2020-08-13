package main

import (
	"github.com/jsd/asset-price-cl-ea/app"
	"github.com/jsd/bridges"
)

func main() {
	c := app.NewConfig()
	app.StartPairsTicker(c)

	bridges.NewServer(&app.AssetPrice{}).Start(c.Port)
}
