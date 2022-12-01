IMAGE_REF=opendolphin/backend
BINARY_NAME=opendolphin-backend

build:
	mkdir -p build
	CGO_ENABLED=0 go build -o "./build/$(BINARY_NAME)" "./cmd/$(BINARY_NAME)"

docker-build:
	docker build \
		-t "$(IMAGE_REF)" \
		.

docker-run:
	docker run \
		--rm \
		-p 5000:5000 \
		"$(IMAGE_REF)"
