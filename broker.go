package main

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/domain"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Broker struct {
	KubeClient  kubernetes.Clientset
	Deployments map[string][]string
	Secret      *apiv1.Secret
}

func (b *Broker) Services(context context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:          "redis",
			Name:        "redis",
			Description: "a proof of concept redis",
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{{
				ID:          "redis-poc",
				Name:        "redis",
				Description: "poc",
			}},
		},
		brokerapi.Service{
			ID:          "postgres",
			Name:        "postgres",
			Description: "a proof of concept postgresql instance",
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{{
				ID:          "postgres-poc",
				Name:        "postgres",
				Description: "poc",
			}},
		},
	}, nil
}

func (b *Broker) Provision(context context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	fmt.Print(fmt.Sprintf("Provisioning new instance: %s %s/%s \n", instanceID, details.ServiceID, details.PlanID))
	spec := brokerapi.ProvisionedServiceSpec{}

	deploymentName := details.ServiceID + "-" + instanceID
	b.Secret, _ = createSecret(&b.KubeClient, "secret.yml", deploymentName)

	_, err := createDeployment(&b.KubeClient, "deployment.yml", deploymentName, details.ServiceID)
	if err != nil {
		bailWith("Failed to create deployment: %s", err)
	}

	_, err = createService(&b.KubeClient, "service.yml", deploymentName, details.ServiceID)
	if err != nil {
		bailWith("failed to create service: %s", err)
	}

	b.Deployments[details.ServiceID] = append(b.Deployments[details.ServiceID], instanceID)

	return spec, nil
}

func (b *Broker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fmt.Print(fmt.Sprintf("Deleting  instance: %s %s/%s \n", instanceID, details.ServiceID, details.PlanID))
	deploymentName := details.ServiceID + "-" + instanceID
	err := deleteDeployment(&b.KubeClient, deploymentName)
	if err != nil {
		bailWith("Failed to delete deployment: %s", err)
	}

	err = deleteService(&b.KubeClient, deploymentName)
	if err != nil {
		bailWith("Failed to delete service: %s", err)
	}

	err = deleteSecret(&b.KubeClient, deploymentName)

	b.Fuhgettaboutit(b.Deployments[details.ServiceID], instanceID)

	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b *Broker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	secretMap := make(map[string]string)

	for key, secret := range b.Secret.Data {
		secret := string([]byte(secret))
		secretMap[key] = secret
	}

	binding.Credentials = secretMap

	//fmt.Print(fmt.Sprintf("Connect to your serivce using the following creds: \nUsername: %s \nPassword: %s", secretMap["username"], secretMap["password"]))

	return binding, nil

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

func (b *Broker) Fuhgettaboutit(slice []string, strToRemove string) []string {
	var i int
	for position, word := range slice {
		if word == strToRemove {
			i = position
		}
	}
	slice[i] = slice[len(slice)-1]
	slice[len(slice)-1] = ""
	return slice[:len(slice)-1]
	//this is my favorite function ive ever written - ahartpence
}
