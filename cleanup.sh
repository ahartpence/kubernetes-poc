#!/bin/bash

kubectl delete deployment --all
kubectl delete service postgres-service
kubectl delete secret postgres-credentials
