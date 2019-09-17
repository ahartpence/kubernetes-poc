package main

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/domain"
	"k8s.io/client-go/kubernetes"
)

type Broker struct {
	KubeClient kubernetes.Clientset
}

type InstanceCreator interface {
	Create(instanceID string) error
	Destroy(instanceID string) error
	InstanceExists(instanceID string) (bool, error)
}

func (b *Broker) Services(context context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:          "redis",
			Name:        "redis-one",
			Description: "number one redis",
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{{
				ID:          "a",
				Name:        "b",
				Description: "c",
			}},
		},
	}, nil
}

func (b *Broker) Provision(context context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	spec := brokerapi.ProvisionedServiceSpec{}

	fmt.Println("Provisioning new instance: ", details.ServiceID, instanceID)

	deploymentName = details.ServiceID + "-" + instanceID
	createDeployment(&b.KubeClient, "deployment.yml")

	return spec, nil
}

func (b *Broker) Deprovision(context context.Context, instanceID string,
	details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b *Broker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, nil
}

func (b *Broker) Unbind(context context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {
	return domain.UnbindSpec{}, nil
}

func (b *Broker) Update(context context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}

func (b *Broker) LastOperation(context context.Context, instanceID string, details domain.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

func (b *Broker) GetBinding(context context.Context, instanceID, bindID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, nil
}

func (b *Broker) GetInstance(context context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, nil
}

func (b *Broker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{}, nil
}
