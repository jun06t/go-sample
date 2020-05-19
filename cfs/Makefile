.PHONY: build push

REGISTRY_DOMAIN := jun06t

build:
	docker build -t ${REGISTRY_DOMAIN}/go-cfs:latest .

push:
	docker push ${REGISTRY_DOMAIN}/go-cfs


