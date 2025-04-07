CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

production-console:
	kubectl exec -it $(shell kubectl get pod -l app.kubernetes.io/name=app-console-pod -o custom-columns=:metadata.name --no-headers) -- bin/rails c -s

k8s-ubuntu-setup:
	snap install helm --classic
	snap install kubectl --classic
	snap install doctl

k8s-macos-setup:
	brew install helm doctl kubectl

helm-install-app:
	helm install app $(CURRENT_DIR)/app-chart

helm-upgrade-app:
	helm upgrade app $(CURRENT_DIR)/app-chart
	git reset HEAD
	git add $(CURRENT_DIR)/app-chart/values.yaml
	git commit --no-verify -m ":rocket: deploy"
	git push --no-verify

helm-install-grafana-alloy:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	helm install grafana-alloy -f grafana-alloy-values.yaml grafana/alloy

helm-upgrade-grafana-alloy:
	helm upgrade grafana-alloy -f grafana-alloy-values.yaml grafana/alloy

kube-docker-install:
	kubectl apply -f docker-installer.daemonset.yaml

get-pods:
	kubectl get pod -w -o wide

describe-app:
	kubectl describe pod app
