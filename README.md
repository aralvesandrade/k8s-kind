Criar cluster usando cli `KIND`(https://kind.sigs.k8s.io/)

```
kind create cluster --config=k8s/kind-config.yaml
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

Deletar cluster

```
kind delete cluster
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
kubectl apply -f k8s/metrics-server.yaml
kubectl apply -f k8s/result-analyzer-program/deployment.yaml
```

Comandos

```
kubectl get deployments
kubectl get pods
kubectl logs svc/producer -f
kubectl logs deployment/consumer -f
kubectl logs pod/{name-pod} -f
kubectl port-forward svc/rabbitmq 15672:15672
kubectl port-forward svc/producer 5001:5001
kubectl get services
```

K8s Dashboard UI

```
sudo snap install helm --classic
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm repo update
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
kubectl get pods -n kubernetes-dashboard
kubectl apply -f k8s/kubernetes-dashboard.yaml
kubectl -n kubernetes-dashboard create token admin-user
kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443
```

Keda

```
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda --namespace keda --create-namespace
kubectl get pods -n keda
echo -n 'amqp://guest:guest@rabbitmq:5672' | base64
echo -n 'amqp://guest:guest@rabbitmq.default:5672/' | base64
```

ArgoCD

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl get pods -n argocd
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
kubectl port-forward svc/argocd-server 5002:443 -n argocd
kubectl apply -f argocd/rabbitmq/values.yaml -f argocd/producer/values.yaml -f argocd/consumer/values.yaml
kubectl apply -f argocd/consumer2/values.yaml
kubectl apply -f argocd/consumer3/values.yaml
```

Prometheus e Grafana

```
kubectl create namespace monitoring
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus-operator prometheus-community/kube-prometheus-stack -n monitoring
kubectl get pods -n monitoring
kubectl get services -n monitoring
kubectl port-forward svc/prometheus-operated -n monitoring 9090:9090
kubectl port-forward svc/prometheus-operator-grafana -n monitoring 3000:80
kubectl get secret --namespace monitoring prometheus-operator-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

Aplicando um stress test

```
kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 10 -t 10s -c 4 "http://producer:5001/hello"
```
