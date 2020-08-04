all:
	go build
	mkdir -p ~/.terraform.d/plugins
	cp terraform-provider-vaultutility ~/.terraform.d/plugins

clean:
	rm -f terraform-provider-vaultutility ~/.terraform.d/plugins/terraform-provider-vaultutility

release:
	goreleaser release --rm-dist
