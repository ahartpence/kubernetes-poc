## A Kubernetes Service Broker, Deploying a 3-pod Redis and Postgres

![a 4-split terminal doing the serice broker thing](image.png)


### Running

To use the service broker, clone the repository then use:
`go build` 


## Important Notes

The service broker will target whatever Kubernetes instance your `kubeconfig` is currently pointed at. Make sure to update your `kubeconfig` accordingly.
