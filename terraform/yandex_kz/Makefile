init:
	@sops -d --extract '["terraform"]["yc_kz"]["service_key_data"]' ../../k8s/secrets.yaml > yc_config.json
	@echo "access_key = \"$(shell sops -d --extract '["terraform"]["yc_kz"]["state_storage"]["access_key"]' ../../k8s/secrets.yaml)\"" > backend.secrets.auto.tfvars
	@echo "secret_key = \"$(shell sops -d --extract '["terraform"]["yc_kz"]["state_storage"]["secret_key"]' ../../k8s/secrets.yaml)\"" >> backend.secrets.auto.tfvars
	@terraform init \
		-backend-config="access_key=$(shell sops -d --extract '["terraform"]["yc_kz"]["state_storage"]["access_key"]' ../../k8s/secrets.yaml)" \
		-backend-config="secret_key=$(shell sops -d --extract '["terraform"]["yc_kz"]["state_storage"]["secret_key"]' ../../k8s/secrets.yaml)"
