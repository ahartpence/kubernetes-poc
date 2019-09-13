package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/jhunt/go-ansi"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//login to the bufflab k8s cluster
	cfg, err := clientcmd.BuildConfigFromFlags("https://10.128.4.17:6443", fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")))
	if err != nil {
		bailWith("Failed creating config: %s", err)
	}
	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		bailWith("Failed creating client set: %s", err)
	}

	//get the deployment file, unmarshal it into `deployment`, and then create it
	deploymentFile, err := ioutil.ReadFile("deployment.yml")
	if err != nil {
		bailWith("Failed to read deployment from file: %s", err)
	}

	var deployment appsv1.Deployment
	err = yaml.Unmarshal(deploymentFile, &deployment)
	if err != nil {
		bailWith("Failed to parse yaml: %s", err)
	}

	fmt.Print(&deployment)

	//get the current deployment name, set it to nil then generate a new name based on what was already present in the file
	deploymentName := deployment.ObjectMeta.Name
	deployment.ObjectMeta.Name = ""
	deployment.ObjectMeta.GenerateName = deploymentName + "-"

	deploymentInterface := k8sClient.AppsV1().Deployments("ahartpence")
	_, err = deploymentInterface.Create(&deployment)
	if err != nil {
		bailWith("Failed to deploy postgres: %s", err)
	}

	fmt.Println("Created deployment")

	username := "tom"
	password := "andrew"
	_, err = createSecret(k8sClient, username, password)
	if err != nil {
		bailWith("Failed when creating secret: %s", err)
	}

	fmt.Println("Created Secret")

	//create the service based on a file, "service.yml"
	serviceFile, err := ioutil.ReadFile("service.yml")
	if err != nil {
		bailWith("Failed to read service file: %s", err)
	}
	var service apiv1.Service
	err = yaml.Unmarshal(serviceFile, &service)
	if err != nil {
		bailWith("Failed to parse service yaml %s", err)
	}
	serviceInterface := k8sClient.CoreV1().Services("ahartpence")

	_, err = serviceInterface.Create(&service)
	if err != nil {
		bailWith("Failed to create service: %s", err)
	}

	fmt.Println("Created Service")

}

func bailWith(format string, args ...interface{}) {
	ansi.Fprintf(os.Stderr, "@R{"+format+"}\n", args...)
	os.Exit(1)
}

func createSecret(client *kubernetes.Clientset, username, password string) (*apiv1.Secret, error) {
	secretsInterface := client.CoreV1().Secrets("ahartpence")
	secretsSpec := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "postgres-credentials",
		},
		StringData: map[string]string{
			"username": username,
			"password": password,
		},
		Type: apiv1.SecretTypeOpaque,
	}

	return secretsInterface.Create(secretsSpec)
}

func int32Ptr(i int32) *int32 { return &i }
