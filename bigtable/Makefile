.PHONY: prepare build proto

REGISTRY_DOMAIN := jun06t

build:
	docker build -f Dockerfile -t ${REGISTRY_DOMAIN}/bigtable-emulator:latest .

push:
	docker push ${REGISTRY_DOMAIN}/bigtable-emulator
 
