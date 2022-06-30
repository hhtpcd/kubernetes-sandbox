BUILD_ID=dev
IMAGE_REPOSITORY=localhost
IMAGE_ID=echo-server
IMAGE_NAME=${IMAGE_REPOSITORY}/${IMAGE_ID}
GITSHA=$(shell git rev-parse --short HEAD)

.PHONY: test
test:
	go test -v .

.PHONY: build
build:
	go build -o server

.PHONY: build-image
build-image:
	skaffold build

.PHONY: push-image
push-image:
	skaffold build --push

.PHONY: start-kube
start-kube:
	# 1.22 is EOL in 2022-10, still OK for now.
	kind create cluster --name kind-cluster --image kindest/node:v1.22.9

.PHONY: clean
clean:
	rm -rf server
	kind delete cluster --name kind-cluster

.PHONY: kube-lint
kube-lint:
	kustomize build kubernetes/base/ | kubeconform -summary -strict -verbose -kubernetes-version 1.22.9 -ignore-missing-schemas
