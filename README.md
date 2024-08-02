Criar cluster usando cli `KIND`(https://kind.sigs.k8s.io/)

```
kind create cluster
```

Listar informações do cluster

```
kind get clusters
#kubectl cluster-info --context kind-{name-cluster}
kubectl cluster-info --context kind-kind
```

Listar todos os clusters

```
kubectl config get-clusters
```

Criar imagem e publicar imagem docker

```
docker buildx build -t aralvesandrade/producer ./src/producer
docker push aralvesandrade/producer

docker buildx build -t aralvesandrade/consumer ./src/consumer
docker push aralvesandrade/consumer
```

Aplicar manifestos

```
kubectl apply -f k8s/rabbitmq-deployment.yaml
kubectl apply -f k8s/producer/deployment.yaml
kubectl apply -f k8s/consumer/deployment.yaml
kubectl apply -f k8s/result-analyzer-program/deployment.yaml
```

Comandos

```
kubectl get deployments
kubectl get pods
kubectl logs svc/producer -f
kubectl logs deployment/consumer -f
kubectl logs pod/{name-pod} -f
kubectl port-forward svc/producer 5001:5001
kubectl get services
```

K8s Dashboard UI

```
sudo snap install helm --classic
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
kubectl apply -f k8s/kubernetes-dashboard.yaml
kubectl -n kubernetes-dashboard create token admin-user
kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443
```

ArgoCD

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
kubectl port-forward svc/argocd-server 5002:443 -n argocd
kubectl apply -f argocd/consumer/values.yaml
kubectl apply -f argocd/producer/values.yaml
```

Keda

```
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda --namespace keda --create-namespace
echo -n 'amqp://guest:guest@rabbitmq.default:5672/' | base64
```

Aplicando um stress test

```
kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://producer:5001/hello"
```
