package controller

import "github.com/n3wscott/gated-broker/pkg/apis/broker/v1"

func (b *BrokerController) CreateServiceBinding(instanceID, bindingID string, req *v1.BindingRequest) (*v1.CreateServiceBindingResponse, error) {
	return nil, nil
}

func (b *BrokerController) DeleteServiceBinding(instanceID, bindingID, serviceID, planID string) error {
	return nil
}
