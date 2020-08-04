all: ${HOME}/local-bin/ghastly ${HOME}/.terraform.d/plugins/terraform-provider-homeassistant

${HOME}/local-bin/ghastly: $(shell find api cmd -name \*.go)
	go build -o ${HOME}/local-bin/ghastly .
${HOME}/.terraform.d/plugins/terraform-provider-homeassistant: $(shell find api terraform -name \*.go)
	go build -o ${HOME}/.terraform.d/plugins/terraform-provider-homeassistant ./terraform
