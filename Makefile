#default value
env=dev

terraform-plan:
	git clean -xfd
	cd terraform; terraform init -backend-config bucket="jxpress-gou-$(env)"
	cd terraform; terraform plan -var-file="$(env).tfvars"

terraform-apply:
	git clean -xfd
	cd terraform; terraform init -backend-config bucket="jxpress-gou-$(env)"
	cd terraform; terraform apply -auto-approve -var-file="$(env).tfvars"
