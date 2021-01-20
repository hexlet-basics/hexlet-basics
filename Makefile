include make-app.mk
include make-compose-app.mk
include make-compose.mk
include k8s/Makefile

project-setup: ansible-generate-env compose-setup

ansible-generate-env:
	docker run --rm -e RUNNER_PLAYBOOK=ansible/development.yml \
		-v $(CURDIR)/ansible/development:/runner/inventory \
		-v $(CURDIR):/runner/project \
		ansible/ansible-runner

ansible-terraform-vars-generate:
	docker run --rm -e RUNNER_PLAYBOOK=ansible/terraform.yml \
		-v $(CURDIR)/ansible/production:/runner/inventory \
		-v $(CURDIR):/runner/project \
		-e ANSIBLE_VAULT_PASSWORD_FILE=tmp/ansible-vault-password \
		ansible/ansible-runner

ansible-vaults-edit:
	docker run -it --rm \
		-v $(CURDIR):/runner/project \
		ansible/ansible-runner ansible-vault edit --vault-password-file project/tmp/ansible-vault-password project/ansible/production/group_vars/all/vault.yml

tag:
	git tag $(TAG) && git push --tags

editor-setup:
	bundle exec solargraph bundle
