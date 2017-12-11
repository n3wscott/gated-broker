package controller

import (
	"errors"

	"github.com/n3wscott/gated-broker/pkg/client"
	"github.com/n3wscott/osb-framework-go/pkg/apis/broker/v2"
)

func (b *BrokerController) GetCatalog() (*v2.Catalog, error) {
	b.hub.Broadcast <- []byte("GetCatalog")

	r := client.NewRequest(b.hub)

	resp := <-r.Send

	if resp.Approved {

		catalog := v2.Catalog{}

		return &catalog, nil
	}
	return nil, errors.New("not approved")
}
