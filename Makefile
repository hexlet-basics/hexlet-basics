include make-services-web.mk
include make-compose.mk

ansible-terraform-vars-generate:
	docker run --rm -e RUNNER_PLAYBOOK=ansible/terraform.yml \
		-v $(CURDIR)/ansible/production:/runner/inventory \
		-v $(CURDIR):/runner/project \
		ansible/ansible-runner ansible-playbook --vault-password-file project/tmp/ansible-vault-password project/ansible/terraform.data.yml -i project/ansible/production -vv

ansible-vaults-edit:
	docker run -it --rm \
		-v $(CURDIR):/runner/project \
		ansible/ansible-runner ansible-vault edit --vault-password-file project/tmp/ansible-vault-password project/ansible/production/group_vars/all/vault.yml

tag:
	git tag $(TAG) && git push --tags
