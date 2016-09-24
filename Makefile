LINTER_ARGS = -j 8 --vendor --enable=misspell --enable=gofmt --enable=goimports --disable=dupl --disable=gocyclo --disable=errcheck --disable=golint --disable=interfacer --deadline=10m --tests

LOCAL_DOCKER_TAG = engops/ddvote
GOOGLE_CONTAINER_REGISTRY_TAG = gcr.io/keen-autumn-144321/engops/ddvote

test: ## Run tests on non-vendored packages
	go test -v $$(glide nv)

binary: ## Construct a binary
	go build -o dd-vote

docker: ## Invoke 'docker build' with the correct arguments
	docker build -t engops/ddvote .
	docker tag $(LOCAL_DOCKER_TAG) $(GOOGLE_CONTAINER_REGISTRY_TAG)

push: test binary docker ## Invoke 'gcloud docker push' with the correct arguments
	gcloud docker push $(GOOGLE_CONTAINER_REGISTRY_TAG)

get-deps: ## install build dependencies
	mkdir -p $$GOPATH/bin
	curl https://glide.sh/get | sh
	glide -v
	glide install --cache --cache-gopath --use-gopath --strip-vcs --strip-vendor

get-gimme: ## install gimme (for installing Go versions)
	scripts/install-gimme.sh

lint: metalint

metalint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	glide novendor | xargs gometalinter $(LINTER_ARGS)

help: ## print list of tasks and descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := binary

.PHONY: all
