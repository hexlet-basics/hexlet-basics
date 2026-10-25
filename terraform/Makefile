init:
	@echo "access_key = \"$(shell sops -d --extract '["terraform"]["state_storage"]["access_key"]' ../k8s/secrets.yaml)\"" > backend.secrets.auto.tfvars
	@echo "secret_key = \"$(shell sops -d --extract '["terraform"]["state_storage"]["secret_key"]' ../k8s/secrets.yaml)\"" >> backend.secrets.auto.tfvars
	terraform init -backend-config=backend.secrets.auto.tfvars
