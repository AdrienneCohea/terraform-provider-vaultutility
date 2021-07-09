all:
	go build
	mkdir -p ~/.terraform.d/plugins/localhost/providers/vaultutility/0.1/linux_amd64
	cp terraform-provider-vaultutility ~/.terraform.d/plugins/localhost/providers/vaultutility/0.1/linux_amd64

clean:
	rm -f terraform-provider-vaultutility ~/.terraform.d/plugins/localhost/providers/vaultutility/0.1/linux_amd64/terraform-provider-vaultutility

release:
	goreleaser release --rm-dist
