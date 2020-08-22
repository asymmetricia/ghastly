all: ${HOME}/local-bin/ghastly ${HOME}/.terraform.d/plugins/terraform-provider-homeassistant

docs/api.md: $(shell find api -name \*.go)
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown ./api > docs/api.md
${HOME}/local-bin/ghastly: $(shell find api cmd -name \*.go)
	go build -o ${HOME}/local-bin/ghastly .
${HOME}/.terraform.d/plugins/terraform-provider-homeassistant: $(shell find api terraform -name \*.go)
	go build -o ${HOME}/.terraform.d/plugins/terraform-provider-homeassistant ./terraform
