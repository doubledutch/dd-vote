LINTER_ARGS = -j 8 --vendor --enable=misspell --enable=gofmt --enable=goimports --disable=dupl --disable=gocyclo --disable=errcheck --disable=golint --disable=interfacer --deadline=10m --tests

test: ## Run tests on non-vendored packages
	go test -v $$(glide nv)

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

.DEFAULT_GOAL := test

.PHONY: all
