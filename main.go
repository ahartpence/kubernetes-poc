package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/jhunt/go-ansi"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
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

	//deploy the secret from a yaml file
	_, err = createSecret(k8sClient, "secret.yml")
	if err != nil {
		bailWith("Failed when creating secret: %s", err)
	}

	fmt.Print("Secret created")

	//create the deployment from the file deployment.yml
	_, err = createDeployment(k8sClient, "deployment.yml")
	if err != nil {
		bailWith("Failed to create deployment: %s", err)
	}

	fmt.Print("deployment created")

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

	fmt.Println("testing directory for yml files")
	regex, _ := regexp.Compile("yml")

	files, err := listDir("test", regex)
	if err != nil {
		bailWith("Failed to do something with the directory")
	}

	fmt.Print(files)

}

func bailWith(format string, args ...interface{}) {
	ansi.Fprintf(os.Stderr, "@R{"+format+"}\n", args...)
	os.Exit(1)
}

func createSecret(client *kubernetes.Clientset, fileName string) (*apiv1.Secret, error) {
	secretsInterface := client.CoreV1().Secrets("ahartpence")
	secretFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		bailWith("Failed to read file: %s", err)
	}

	var secret apiv1.Secret
	err = yaml.Unmarshal(secretFile, &secret)

	return secretsInterface.Create(&secret)
}

func createDeployment(client *kubernetes.Clientset, fileName string) (*appsv1.Deployment, error) {
	deploymentsInterface := client.AppsV1().Deployments("ahartpence")
	deploymentFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		bailWith("Failed to read deployment file: %s", err)
	}

	var deployment appsv1.Deployment
	err = yaml.Unmarshal(deploymentFile, &deployment)

	return deploymentsInterface.Create(&deployment)

}

func listDir(directory string, filter *regexp.Regexp) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		bailWith("Unable to read directory")
	}

	for _, file := range fileInfo {
		fmt.Println(file.Name())
		if filter.MatchString(file.Name()) {
			files = append(files, file.Name())
		}
	}

	return files, nil

}

func int32Ptr(i int32) *int32 { return &i }
