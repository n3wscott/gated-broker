package controller

import "github.com/n3wscott/gated-broker/pkg/apis/broker/v1"

func (b *BrokerController) CreateServiceInstance(ID string, req *v1.CreateServiceInstanceRequest) (*v1.CreateServiceInstanceResponse, int, error) {
	return nil, 0, nil
}

func (b *BrokerController) UpdateServiceInstance(ID string, req *v1.CreateServiceInstanceRequest) (*v1.ServiceInstance, int, error) {
	return nil, 0, nil
}

func (b *BrokerController) DeleteServiceInstance(ID string, req *v1.DeleteServiceInstanceRequest) (*v1.DeleteServiceInstanceResponse, int, error) {
	return nil, 0, nil
}

func (b *BrokerController) PollServiceInstance(ID string, req *v1.LastOperationRequest) (*v1.LastOperationResponse, int, error) {
	return nil, 0, nil
}
