# Deploy of the Go stack (app-chart) to the codebasics namespace. No CI deploys
# it: after a release is merged, `git pull` the bumped Chart.yaml and run
# `make helm-upgrade-app` here.
#
# Secrets stay in the sops-encrypted legacy/k8s/secrets.yaml (age recipient in
# legacy/k8s/.sops.yaml): legacy/terraform reads the same file by relative
# path, so it moves when legacy/ is retired, not before.
CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NAMESPACE := codebasics
SECRETS := $(CURRENT_DIR)/../legacy/k8s/secrets.yaml

k8s-utils-install:
	brew install helm kubectl sops age
	helm plugin install https://github.com/jkroepke/helm-secrets --version v4.6.5 || true

secrets-edit:
	helm secrets edit $(SECRETS)

helm-template-app:
	helm secrets template codebasics $(CURRENT_DIR)/app-chart -n $(NAMESPACE) -f $(SECRETS)

helm-upgrade-app:
	helm secrets upgrade codebasics $(CURRENT_DIR)/app-chart -n $(NAMESPACE) -f $(SECRETS)

# The cutover upgrade only: the database already has the schema Rails left
# behind, so atlas records the baseline migration as applied instead of running
# it (ADR-0015). Every later release uses helm-upgrade-app.
helm-upgrade-app-baseline:
	helm secrets upgrade codebasics $(CURRENT_DIR)/app-chart -n $(NAMESPACE) -f $(SECRETS) \
		--set migrate.baseline=20260727053619
