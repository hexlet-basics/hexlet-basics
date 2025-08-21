CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NAMESPACE := codebasics

k8s-utils-install:
	brew install helm kubectl sops age
	helm plugin install https://github.com/jkroepke/helm-secrets --version v4.6.5 || true

k8s-init:
	kubectl create namespace $(NAMESPACE)
	kubectl config set-context --current --namespace=$(NAMESPACE)

secrets-edit:
	helm secrets edit secrets.yaml

production-console:
	kubectl exec -it $(shell kubectl get pod -l app.kubernetes.io/name=$(NAMESPACE)-console-pod -o custom-columns=:metadata.name --no-headers) -- bin/rails c -s

helm-install-app:
	helm secrets install codebasics $(CURRENT_DIR)/app-chart -f secrets.yaml

helm-upgrade-app:
	helm secrets upgrade codebasics $(CURRENT_DIR)/app-chart -f secrets.yaml
	git reset HEAD
	git add $(CURRENT_DIR)/app-chart/values.yaml
	git commit --no-verify -m ":rocket: deploy"
	git push --no-verify
