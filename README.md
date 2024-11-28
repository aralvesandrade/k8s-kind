No exemplo abaixo cito duas opções de ferramentas de linha de comando (CLI) que ajudam a criar e gerenciar clusters Kubernetes localmente para desenvolvimento e testes, sendo elas: `kind`(https://kind.sigs.k8s.io/) ou `minikube`(https://minikube.sigs.k8s.io/docs/)

# Kind

Criar um cluster usando `kind`

```
kind create cluster
#ou
kind create cluster --config=k8s/kind-config.yaml
#ou
kind create cluster --config=k8s/kind-config.yaml --name {nome_cluster}
#ou usar mesma rede do docker já existeste
docker network ls
export KIND_EXPERIMENTAL_DOCKER_NETWORK=my-network
kind create cluster --config=k8s/kind-config.yaml
```

Listar informações do cluster

```
kind get clusters
```

Exibir informações do cluster Kubernetes que está ativo no contexto especificado, usando o comando `kubectl cluster-info --context kind-{nome_cluster}`

```
kubectl cluster-info --context kind-kind
```

Listar todos os clusters

```
kubectl config get-clusters
```

Listar todos os nodes

```
kubectl get nodes -o wide
```

Deletar cluster

```
kind delete cluster
```

# Minikube

```
minikube start
kubectl get po -A
minikube kubectl -- get po -A
minikube dashboard
kubectl create deployment hello-minikube --image=kicbase/echo-server:1.0
kubectl expose deployment hello-minikube --type=NodePort --port=8080
kubectl get services hello-minikube
minikube service hello-minikube
kubectl port-forward service/hello-minikube 7080:8080
minikube pause
minikube unpause
minikube stop
minikube delete --all
```

# Docker

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

Alguns exemplos de comandos

```
kubectl get deployments
kubectl get pods
kubectl logs svc/producer -f
kubectl logs deployment/consumer -f
kubectl logs pod/{name-pod} -f
kubectl port-forward svc/rabbitmq 15672:15672
kubectl port-forward svc/producer 5001:5001
kubectl get services
kubectl top pods
```

K8s Dashboard UI (no exemplo no kind não se aplica)

```
sudo snap install helm --classic
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm repo update
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
kubectl get pods -n kubernetes-dashboard
kubectl apply -f k8s/kubernetes-dashboard.yaml
kubectl get svc -n kubernetes-dashboard
kubectl -n kubernetes-dashboard create token admin-user
kubectl patch svc kubernetes-dashboard-kong-proxy -p '{"spec": {"type": "NodePort"}}' -n kubernetes-dashboard
kubectl get svc -n kubernetes-dashboard
#ou
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

OpenFaaS

```
helm repo add openfaas https://openfaas.github.io/faas-netes/
helm repo update
kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml
helm upgrade openfaas --install openfaas/openfaas --namespace openfaas
kubectl get pods -n openfaas
kubectl -n openfaas get deployments -l "release=openfaas, app=openfaas"
kubectl get svc -n openfaas
echo $(kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode)
```

ArgoCD

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl get pods -n argocd
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
kubectl port-forward svc/argocd-server 5002:443 -n argocd
#ou
kubectl get svc -n argocd
kubectl patch svc argocd-server -p '{"spec": {"type": "NodePort"}}' -n argocd
kubectl get svc -n argocd
kubectl apply -f argocd/rabbitmq -f argocd/producer -f argocd/consumer
kubectl apply -f argocd/result-analyzer-program
kubectl apply -f argocd/consumer2
kubectl apply -f argocd/consumer3
#ou
kubectl patch svc producer -p '{"spec": {"type": "NodePort"}}'
```

Usando o `minikube`

```
kubectl delete svc argocd-server -n argocd
kubectl expose deployment argocd-server -n argocd --type=NodePort --port=8080
minikube service argocd-server -n argocd
```

Prometheus

```
kubectl create namespace monitoring
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus --namespace monitoring --set=kube-state-metrics.enabled=true
kubectl get pods -n monitoring
kubectl get services -n monitoring
kubectl -n monitoring port-forward svc/prometheus-server 9090:80
```

Grafana

```
kubectl create namespace monitoring
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana --namespace monitoring
kubectl get pods -n monitoring
kubectl get services -n monitoring
kubectl -n monitoring port-forward svc/grafana 3000:80
kubectl get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

Criar conexão Data sources > Prometheus, configurar URL: `http://prometheus-server.monitoring:80`

Para monitorar os containers do K8s precisa importar o dashboard `12740`(https://grafana.com/grafana/dashboards/12740-kubernetes-monitoring/)

Ou algo mais completo!

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

Instalar argocd CLI

```
VERSION=$(curl -L -s https://raw.githubusercontent.com/argoproj/argo-cd/stable/VERSION)
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/download/v$VERSION/argocd-linux-amd64
sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd
rm argocd-linux-amd64
```

Fazer login usando argocd CLI

```
argocd login localhost:5002 --username admin --password $(kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d) --insecure
```

Aplicando um stress test

```
kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 10 -t 10s -c 4 "http://producer:5001/hello"
```

Instalar e aplicar stress teste com K6

```
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

Executar stress test

```
k6 run k6/teste.js
#ou
docker run -i grafana/k6 run - < teste2.js
```
