package main

import (
	"github.com/DODOEX/oracle-adapter/adapter"
	"github.com/DODOEX/oracle-adapter/bridge"
)

func main() {
	c := adapter.NewConfig()
	adapter.StartPairsTicker(c)

	bridge.NewServer(&adapter.AssetPrice{}).Start(c.Port)
}
