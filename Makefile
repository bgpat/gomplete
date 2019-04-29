PACKAGE := github.com/bgpat/gomplete

.PHONY: test
test: go-test e2e-test

.PHONY: go-test
go-test:
	GO111MODULE=on go test -cover -race ./...

.PHONY: e2e-test
e2e-test: e2e-bash e2e-zsh e2e-fish

.PHONY: e2e-bash
e2e-bash:
	docker run -it --rm -v$(PWD):/go/src/$(PACKAGE) -w /go/src/$(PACKAGE) -e GO111MODULE=on golang:1.12 go test -tags=bash ./_test/e2e

.PHONY: e2e-zsh
e2e-zsh:
	docker run -it --rm -v$(PWD):/go/src/$(PACKAGE) -w /go/src/$(PACKAGE) -e GO111MODULE=on bgpat/golang-zsh:1.12 go test -tags=zsh ./_test/e2e

.PHONY: e2e-fish
e2e-fish:
	docker run -it --rm -v$(PWD):/go/src/$(PACKAGE) -w /go/src/$(PACKAGE) -e GO111MODULE=on bgpat/golang-fish:1.12 go test -tags=fish ./_test/e2e
