CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# kubectl exec app-deployment-5b566d5fb5-q9kdr -c phoenix -- docker pull hexletbasics/exercises-php

production-console:
	    kubectl exec -it $(shell kubectl get pod -lapp.kubernetes.io/name=app-console-pod -o=custom-columns=NAME:metadata.name | sed -n 2p) -- bin/rails c -s

init:
	brew install kubernetes-helm

local-cluster-setup:
	sudo snap install doctl
	doctl auth init
	doctl kubernetes cluster kubeconfig save hexlet-basics-3

apply-exercises: apply-exercises-go apply-exercises-ruby apply-exercises-elixir apply-exercises-html apply-exercises-css apply-exercises-js apply-exercises-java apply-exercises-racket apply-exercises-php apply-exercises-python

pull-exercise-image:
	kubectl get pod -l app.kubernetes.io/name=app-web-pod -o custom-columns=name:metadata.name --no-headers --field-selector=status.phase=Running | xargs -I{} kubectl exec {} -c phoenix docker pull hexletbasics/exercises-$(LANG)

apply-exercises-ruby:
	make pull-exercise-image LANG=ruby
	kubectl delete jobs exercises-ruby-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-ruby.job.yml

apply-exercises-elixir:
	make pull-exercise-image LANG=elixir
	kubectl delete jobs exercises-elixir-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-elixir.job.yml

apply-exercises-html:
	make pull-exercise-image LANG=html
	kubectl delete jobs exercises-html-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-html.job.yml

apply-exercises-css:
	make pull-exercise-image LANG=css
	kubectl delete jobs exercises-css-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-css.job.yml

apply-exercises-js:
	make pull-exercise-image LANG=javascript
	kubectl delete jobs exercises-javascript-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-javascript.job.yml

apply-exercises-java:
	make pull-exercise-image LANG=java
	kubectl delete jobs exercises-java-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-java.job.yml

apply-exercises-racket:
	make pull-exercise-image LANG=racket
	kubectl delete jobs exercises-racket-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-racket.job.yml

apply-exercises-php:
	make pull-exercise-image LANG=php
	kubectl delete jobs exercises-php-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-php.job.yml

apply-exercises-python:
	make pull-exercise-image LANG=python
	kubectl delete jobs exercises-python-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-python.job.yml

apply-exercises-go:
	make pull-exercise-image LANG=go
	kubectl delete jobs exercises-go-job || true
	kubectl apply -f $(CURRENT_DIR)/exercises-go.job.yml

helm-install: helm-install-nginx-ingress helm-install-metrics
	helm install app ./app-chart
	helm install -f ./values-external-dns.yml external-dns stable/external-dns

helm-install-ingress-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm install nginx-ingress ingress-nginx/ingress-nginx

helm-install-metrics:
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	helm repo update
	helm install metrics-server stable/metrics-server --set args={--kubelet-preferred-address-types=InternalIP}


# helm install --name app ./app-chart
helm-upgrade-app:
	helm upgrade  app $(CURRENT_DIR)/app-chart
	git reset HEAD
	git add $(CURRENT_DIR)/app-chart/values.yaml
	git commit -m "deploy"
	git push

# helm install -f ./values-external-dns.yml --name external-dns stable/external-dns
helm-upgrade-external-dns:
	helm upgrade  -f $(CURRENT_DIR)/values-external-dns.yml external-dns stable/external-dns

get-pods:
	kubectl get pod -w -o wide

describe-app:
	kubectl describe pod app

log-app:
	kubectl logs -lapp=app -c phoenix

start-bash:
	kubectl delete pods alpine || true
	kubectl run -i --tty alpine --image=alpine --restart=Never -- sh
