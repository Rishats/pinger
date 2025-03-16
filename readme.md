This document explains how to build, tag, push the Docker container to DockerHub, and deploy it using Kubernetes.

---

# Pinger Service

This repository contains a simple Go-based pinger service that monitors multiple hosts and exposes metrics for Prometheus.

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [DockerHub account](https://hub.docker.com/)
- [Kubernetes (kubectl)](https://kubernetes.io/docs/tasks/tools/)
- [Minikube (optional for local testing)](https://minikube.sigs.k8s.io/docs/start/)

---

## Running Locally with Docker
To test the container locally:
```sh
docker run -d -p 2112:2112 --name pinger -e PINGER_DOMAINS="vk.com,google.com,megaline.kg,aknet.kg,yandex.ru" skeletondocker/pinger:latest
```
Then, visit:
```
http://localhost:2112/metrics
```

---

## Deploying to Kubernetes

1. Ensure you have a running Kubernetes cluster.
2. Apply the manifest:

```sh
kubectl apply -f k8s-manifest.yaml
```

3. Check the deployment and service:

```sh
kubectl get pods
kubectl get svc pinger-service
```

---

## Accessing Metrics in Kubernetes
If running locally with Minikube, you can access the service using:
```sh
minikube service pinger-service
```
For external access, configure an ingress or use `kubectl port-forward`:
```sh
kubectl port-forward svc/pinger-service 2112:2112
```
Now, visit:
```
http://localhost:2112/metrics
```
```

---

### `k8s-manifest.yaml`:
This Kubernetes manifest deploys the `pinger` service using the `skeletondocker/pinger:latest` image and exposes it on port `2112`.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pinger-deployment
  labels:
    app: pinger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pinger
  template:
    metadata:
      labels:
        app: pinger
    spec:
      containers:
        - name: pinger
          image: skeletondocker/pinger:latest
          ports:
            - containerPort: 2112
          env:
            - name: PINGER_DOMAINS
              value: "vk.com,google.com,megaline.kg,aknet.kg,yandex.ru"

---
apiVersion: v1
kind: Service
metadata:
  name: pinger-service
spec:
  selector:
    app: pinger
  ports:
    - protocol: TCP
      port: 2112
      targetPort: 2112
  type: ClusterIP
```

---

### Instructions for Deploying the Kubernetes Manifest:

1. Save the `k8s-manifest.yaml` file.
2. Apply the manifest:
   ```sh
   kubectl apply -f k8s-manifest.yaml
   ```
3. Check if the pod and service are running:
   ```sh
   kubectl get pods
   kubectl get svc pinger-service
   ```
4. To access the metrics, use `kubectl port-forward`:
   ```sh
   kubectl port-forward svc/pinger-service 2112:2112
   ```
5. Open a browser or use `curl`:
   ```sh
   curl http://localhost:2112/metrics
   ```

---

---
## Building and Pushing the Docker Image

### 1. Clone the Repository
```sh
git clone https://github.com/Rishats/pinger.git
cd pinger
```

### 2. Build the Docker Image
```sh
docker build -t yourdockerhubusername/pinger:latest .
docker tag yourdockerhubusername/pinger:latest yourdockerhubusername/pinger:v1.0.0
```

### 3. Log in to DockerHub
```sh
docker login
```

### 4. Push the Images to DockerHub
```sh
docker push yourdockerhubusername/pinger:latest
docker push yourdockerhubusername/pinger:v1.0.0
```

---



This setup ensures your `pinger` service is running inside Kubernetes and exposing metrics properly. 🚀
