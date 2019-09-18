#!/bin/bash

kubectl delete deployment --all -n ahartpence
kubectl delete service --all -n ahartpence
kubectl delete secret  --all -n ahartpence

