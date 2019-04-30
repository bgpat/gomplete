PACKAGE := github.com/bgpat/gomplete
COUNT := 1

.PHONY: test
test: go-test integration-test

.PHONY: go-test
go-test:
	GO111MODULE=on go test -cover -race -count=$(COUNT) ./...

.PHONY: integration-test
integration-test: integration-test-bash integration-test-zsh integration-test-fish

.PHONY: integration-test-bash
integration-test-bash:
	docker run -it --rm -v$(PWD):/go/src/$(PACKAGE) -w /go/src/$(PACKAGE) -e GO111MODULE=on golang:1.12 go test -tags="integration bash" -count=$(COUNT) -v ./shells/bash/integration_test.go

.PHONY: integration-test-zsh
integration-test-zsh:
	docker run -it --rm -v$(PWD):/go/src/$(PACKAGE) -w /go/src/$(PACKAGE) -e GO111MODULE=on bgpat/golang-zsh:1.12 go test -tags="integration zsh" -count=$(COUNT) -v ./shells/zsh/integration_test.go

.PHONY: integration-test-fish
integration-test-fish:
	docker run -it --rm -v$(PWD):/go/src/$(PACKAGE) -w /go/src/$(PACKAGE) -e GO111MODULE=on bgpat/golang-fish:1.12 go test -tags="integration fish" -count=$(COUNT) -v ./shells/fish/integration_test.go
