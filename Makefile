build::
	docker build -t vxlabs/vault-config-extractor .
	docker create --name artifacts vxlabs/vault-config-extractor
	docker cp artifacts:/bin/vault-config-extractor vault-config-extractor
	docker rm artifacts

