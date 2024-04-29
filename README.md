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

## Step 5: Test changing settings
0. [Optional] For pods reload control you can open separated terminal and run command  
```bash
kubectl get pods -n app --watch
```

1. Get logs of any server-1 pod
```bash
kubectl logs server-1-deployment-c686894f8-sl22m -n app
```

2. find there time priod value by output:
```bash
Wait N seconds
```

3. Change value timeRange for the server-1 in the sample/cr-sample.yaml

```yaml
    - name: server-1
      replicas: 5
      port: 8080
      timeRange: 10
```

4. Apply changes

```bash
kubectl apply -n app -f sample/cr-sample.yaml
```
Server-1 pods will not reload

5. Wait 5 minutes and reproduce steps 1-2

## Step 6: Test how discovery works
0. [Optional] For pods reload control you can open separated terminal and run command  
```bash
kubectl get pods -n app --watch
```

1. Get logs of any server-1 pod
```bash
kubectl logs server-1-deployment-c686894f8-sl22m -n app
```
2. find there time priod value by output:
```bash
2 ips connected to the service server-2-headless found.
contact ip: 10.1.1.206:8081
Response from 10.1.1.206:8081: pong
```
3. Change value timeRange for the server-1 in the sample/cr-sample.yaml

```yaml
    - name: server-2
      replicas: 5
```

4. Apply changes

```bash
kubectl apply -n app -f sample/cr-sample.yaml
```
5. Wait while all pods will up and running and reproduce steps 1-2. You will see new ip addresses and its responses
