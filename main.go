package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"code.cloudfoundry.org/lager"
	"github.com/ghodss/yaml"
	"github.com/jhunt/go-ansi"
	"github.com/pivotal-cf/brokerapi"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/sethvargo/go-password/password"
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

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "andrew",
		Password: "tom",
	}

	deploymentMap := make(map[string][]string)

	logger := lager.NewLogger("poc")

	broker := &Broker{
		KubeClient:  *k8sClient,
		Deployments: deploymentMap,
	}
	handler := brokerapi.New(broker, logger, brokerCredentials)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		bailWith("Server quit: %s", err)
	}

}
func bailWith(format string, args ...interface{}) {
	ansi.Fprintf(os.Stderr, "@R{"+format+"}\n", args...)
	os.Exit(1)
}

func createSecret(client *kubernetes.Clientset, fileName string, uid string) (*apiv1.Secret, error) {
	secretsInterface := client.CoreV1().Secrets("ahartpence")
	secretFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		bailWith("Failed to read file: %s", err)
	}

	var secret apiv1.Secret
	err = yaml.Unmarshal(secretFile, &secret)

	secret.ObjectMeta.Name = uid
	secret.ObjectMeta.Labels = map[string]string{
		"created_by": "blacksmith",
		"service":    uid,
	}

	b64Password, err := password.Generate(10, 2, 0, true, true)
	if err != nil {
		bailWith("Failed to generate password: %s", err)
	}
	secret.Data["password"] = []byte(b64Password)

	return secretsInterface.Create(&secret)
}

func createDeployment(client *kubernetes.Clientset, fileName string, uid string, deployKind string) (*appsv1.Deployment, error) {
	deploymentsInterface := client.AppsV1().Deployments("ahartpence")
	deploymentFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		bailWith("Failed to read deployment file: %s", err)
	}
	var deployment appsv1.Deployment
	err = yaml.Unmarshal(deploymentFile, &deployment)

	deployment.ObjectMeta.Name = uid
	deployment.ObjectMeta.Labels = map[string]string{
		"service": uid,
	}

	deployment.Spec.Selector.MatchLabels = map[string]string{
		"service": uid,
	}

	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{
		"service": uid,
	}

	deployment.Spec.Template.Spec.Containers = []apiv1.Container{{
		Name:  uid,
		Image: deployKind,
		Env: []apiv1.EnvVar{{
			Name: "POSTGRES_USER",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: uid,
					},
					Key: "username",
				},
			},
		},
			{
				Name: "POSTGRES_PASSWORD",
				ValueFrom: &apiv1.EnvVarSource{
					SecretKeyRef: &apiv1.SecretKeySelector{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: uid,
						},
						Key: "password",
					},
				},
			},
		},
	}}

	return deploymentsInterface.Create(&deployment)

}

func createService(client *kubernetes.Clientset, fileName string, uid string, deployKind string) (*apiv1.Service, error) {
	serviceInterface := client.CoreV1().Services("ahartpence")
	serviceFile, err := ioutil.ReadFile("service.yml")
	if err != nil {
		bailWith("Failed to read service file: %s", err)
	}
	var service apiv1.Service
	err = yaml.Unmarshal(serviceFile, &service)
	if err != nil {
		bailWith("Failed to parse service yaml %s", err)
	}

	service.Spec.Selector = map[string]string{
		"service": uid,
	}

	service.ObjectMeta.Name = uid
	service.ObjectMeta.Labels = map[string]string{
		"created_by": "blacksmith",
		"service":    uid,
	}

	return serviceInterface.Create(&service)
}

func deleteDeployment(client *kubernetes.Clientset, uid string) error {
	deploymentsInterface := client.AppsV1().Deployments("ahartpence")

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	return deploymentsInterface.Delete(uid, deleteOptions)
}

func deleteService(client *kubernetes.Clientset, uid string) error {
	serviceInterface := client.CoreV1().Services("ahartpence")

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	return serviceInterface.Delete(uid, deleteOptions)
}

func deleteSecret(client *kubernetes.Clientset, uid string) error {
	secretsInterface := client.CoreV1().Secrets("ahartpence")

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	return secretsInterface.Delete(uid, deleteOptions)
}

func listDir(directory string, filter *regexp.Regexp) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		bailWith("Unable to read directory")
	}

	for _, file := range fileInfo {
		if filter.MatchString(file.Name()) {
			files = append(files, file.Name())
		}
	}

	return files, nil

}

func int32Ptr(i int32) *int32 { return &i }
