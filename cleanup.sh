#!/bin/bash

kubectl delete deployment postgres-deployment
kubectl delete service postgres-service
kubectl delete secret postgres-credentials
