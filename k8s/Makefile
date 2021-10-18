CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

production-console:
	    kubectl exec -it $(shell kubectl get pod -lapp.kubernetes.io/name=app-console-pod -o=custom-columns=NAME:metadata.name | sed -n 2p) -- bin/rails c -s

init:
	brew install kubernetes-helm

local-cluster-setup:
	snap install helm --classic
	snap install kubectl --classic
	snap install doctl
	doctl auth init
	doctl kubernetes cluster kubeconfig save hexlet-basics-4

helm-install: helm-install-ingress-nginx
	helm install app ./app-chart

helm-install-ingress-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm install ingress-nginx ingress-nginx/ingress-nginx

helm-install-metrics:
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	helm repo update
	helm install metrics-server stable/metrics-server --set args={--kubelet-preferred-address-types=InternalIP}

# helm install --name app ./app-chart
helm-upgrade-app:
	helm upgrade app $(CURRENT_DIR)/app-chart
	git reset HEAD
	git add $(CURRENT_DIR)/app-chart/values.yaml
	git commit -m "deploy"
	git push --no-verify

# helm install -f ./values-external-dns.yml --name external-dns stable/external-dns
helm-upgrade-external-dns:
	helm upgrade -f $(CURRENT_DIR)/values-external-dns.yml external-dns stable/external-dns

get-pods:
	kubectl get pod -w -o wide

describe-app:
	kubectl describe pod app

start-bash:
	kubectl delete pods alpine || true
	kubectl run -i --tty alpine --image=alpine --restart=Never -- sh
