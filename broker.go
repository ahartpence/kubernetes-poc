package main

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/domain"
	"k8s.io/client-go/kubernetes"
)

type Broker struct {
	KubeClient  kubernetes.Clientset
	Deployments map[string][]string
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
				Name:        "postgres",
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
	spec := brokerapi.ProvisionedServiceSpec{}

	fmt.Println("Provisioning new instance: ", details.ServiceID, instanceID)

	deploymentName := details.ServiceID + "-" + instanceID
	secret, err := createSecret(&b.KubeClient, "secret.yml", deploymentName)
	if err != nil {
		bailWith("Failed to create secret: %s", err)
	}

	fmt.Println("\t Created Secret")
	_, err = createDeployment(&b.KubeClient, "deployment.yml", deploymentName, details.ServiceID)
	if err != nil {
		bailWith("Failed to create deployment: %s", err)
	}
	fmt.Println("\t Created deployment")
	_, err = createService(&b.KubeClient, "service.yml", deploymentName, details.ServiceID)
	if err != nil {
		bailWith("failed to create service: %s", err)
	}
	fmt.Println("\t Created service")

	b.Deployments[details.ServiceID] = append(b.Deployments[details.ServiceID], instanceID)
	b.GetServices(b.Deployments)
	secretMap := make(map[string]string)
	for key, secret := range secret.Data {
		secret := string([]byte(secret))
		secretMap[key] = secret
	}

	fmt.Print(fmt.Sprintf("Service sucessfully created, your credentials are: \nUsername: %s \nPassword: %s", secretMap["username"], secretMap["password"]))

	return spec, nil
}

func (b *Broker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fmt.Println("Deprovisioning Service :", details.ServiceID)
	deploymentName := details.ServiceID + "-" + instanceID

	fmt.Println(deploymentName)
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

	b.GetServices(b.Deployments)

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

func (b *Broker) Fuhgettaboutit(s []string, strToRemove string) []string {
	var i int
	for p, w := range s {
		if w == strToRemove {
			i = p
		}
	}
	s[i] = s[len(s)-1]
	s[len(s)-1] = ""
	return s[:len(s)-1]
}

func (b *Broker) GetServices(services map[string][]string) error {
	for key, value := range services {
		fmt.Println("Service:", key, "Instances", value)
	}

	return nil
}
