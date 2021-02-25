CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

CONTEXT := --kube-context gke_hexlet-basics_europe-west3-a_hexlet-basics-6

# kubectl exec app-deployment-5b566d5fb5-q9kdr -c phoenix -- docker pull hexletbasics/exercises-php

production-console:
	    kubectl exec -it $(shell kubectl get pod -lapp.kubernetes.io/name=app-console-pod -o=custom-columns=NAME:metadata.name | sed -n 2p) -- bin/rails c -s

init:
	brew install kubernetes-helm

local-cluster-setup:
	sudo snap install doctl
	doctl auth init
	doctl kubernetes cluster kubeconfig save hexlet-basics-3

helm-install: helm-install-nginx-ingress helm-install-metrics
	helm install $(CONTEXT) app ./app-chart
	helm install $(CONTEXT) -f ./values-external-dns.yml external-dns stable/external-dns

helm-install-ingress-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm install $(CONTEXT) nginx-ingress ingress-nginx/ingress-nginx

helm-install-metrics:
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	helm repo update
	helm install $(CONTEXT) metrics-server stable/metrics-server --set args={--kubelet-preferred-address-types=InternalIP}

# helm install --name app ./app-chart
helm-upgrade-app:
	helm upgrade $(CONTEXT) app $(CURRENT_DIR)/app-chart
	git reset HEAD
	git add $(CURRENT_DIR)/app-chart/values.yaml
	git commit -m "deploy"
	git push

# helm install -f ./values-external-dns.yml --name external-dns stable/external-dns
helm-upgrade-external-dns:
	helm upgrade $(CONTEXT) -f $(CURRENT_DIR)/values-external-dns.yml external-dns stable/external-dns

get-pods:
	kubectl get pod -w -o wide

describe-app:
	kubectl describe pod app

start-bash:
	kubectl delete pods alpine || true
	kubectl run -i --tty alpine --image=alpine --restart=Never -- sh
