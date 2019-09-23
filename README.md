## A kubernetes service broker, deploying a 3-pod redis and postgres

![a 4-split terminal doing the serice broker thing](image.png)


### Running

To use the service broker, clone the repository then use
`go build` 


## Important notes

The service broker will target whatever kubernetes instance your `kubeconfig` is current pointed at. Make sure to update your `kubeconfig` accordingly.
