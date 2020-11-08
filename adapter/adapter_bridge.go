package adapter

import (
	"net/http"

	"github.com/DODOEX/oracle-adapter/bridge"
)

type AssetPrice struct{}

func (ap *AssetPrice) Run(h *bridge.Helper) (interface{}, error) {
	return GetPrice(h.GetParam("base"), h.GetParam("quote"))
}

func (ap *AssetPrice) Opts() *bridge.Opts {
	return &bridge.Opts{
		Name:   "Asset Price",
		Lambda: true,
		Path:   "/price",
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	StartPairsTicker(nil)
	bridge.NewServer(&AssetPrice{}).Handler(w, r)
}
