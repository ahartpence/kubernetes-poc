#!/bin/bash

kubectl delete deployment -l created_by=blacksmith
kubectl delete service -l created_by=blacksmith
kubectl delete secret  -l created_by=blacksmith
