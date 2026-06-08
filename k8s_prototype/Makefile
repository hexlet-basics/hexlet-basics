CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NAMESPACE := codebasics
CHART_DIR := $(CURRENT_DIR)/app-chart
RELEASE_NAME := codebasics

# https://github.com/jkroepke/helm-secrets/wiki/Installation
install:
	brew install helm kubectl sops age
	helm plugin install --verify=false https://github.com/jkroepke/helm-secrets/releases/download/v4.7.4/secrets-4.7.4.tgz
	helm plugin install --verify=false https://github.com/jkroepke/helm-secrets/releases/download/v4.7.4/secrets-getter-4.7.4.tgz
	helm plugin install --verify=false https://github.com/jkroepke/helm-secrets/releases/download/v4.7.4/secrets-post-renderer-4.7.4.tgz

k8s-init:
	kubectl create namespace $(NAMESPACE)
	kubectl config set-context --current --namespace=$(NAMESPACE)

secrets-edit:
	helm secrets edit $(CURRENT_DIR)/secrets.yaml

production-console:
	kubectl exec -it $(shell kubectl get pod -l app.kubernetes.io/name=$(NAMESPACE)-console-pod -o custom-columns=:metadata.name --no-headers) -- bin/rails c -s

helm-install-app:
	helm secrets install $(RELEASE_NAME) $(CHART_DIR) -f $(CURRENT_DIR)/secrets.yaml

helm-upgrade-app:
	helm secrets upgrade $(RELEASE_NAME) $(CHART_DIR) -f $(CURRENT_DIR)/secrets.yaml
	git reset HEAD
	git add $(CHART_DIR)/values.yaml
	git commit --no-verify -m ":rocket: deploy"
	git push --no-verify

helm-lint:
	helm lint $(CHART_DIR)

helm-template:
	helm template $(RELEASE_NAME) $(CHART_DIR) -n $(NAMESPACE) > /tmp/hexlet-basics-helm-render.yaml

helm-check: helm-lint helm-template

.PHONY: install k8s-init secrets-edit production-console helm-install-app helm-upgrade-app helm-lint helm-template helm-check
