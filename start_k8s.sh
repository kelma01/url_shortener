#!/bin/bash

minikube delete
minikube start
eval $(minikube docker-env)
docker build -t testing .
minikube image load testing
kubectl apply -f k8s/

