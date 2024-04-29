## Step 1: Build Operator Docker Image

Build the Docker image with the following command:
```bash
docker build . -t radmgrp/operator-test:latest
```
## Step 2: Push Docker Image

Push the Docker image to the registry:
```bash
docker push radmgrp/operator-test:latest
```

## Step 3: Deploy Helm Chart

Upgrade or install the Helm chart with the following command:
```bash
helm upgrade --install --create-namespace -n system operator ./operator-helm
```
## Step 4: Create Namespace

Create a Kubernetes namespace for the application:
```bash
kubectl create ns app
```
## Step 4: Apply Sample Custom Resource

Apply the custom resource YAML file in the newly created namespace:
```bash
kubectl apply -n app -f sample/cr-sample.yaml
```