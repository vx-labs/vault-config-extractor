build::
	docker build -t quay.io/vxlabs/vault-config-extractor .
	docker create --name artifacts quay.io/vxlabs/vault-config-extractor
	docker cp artifacts:/bin/vault-config-extractor quay.io/vault-config-extractor
	docker rm artifacts

