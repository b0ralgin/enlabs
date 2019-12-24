NAME := enlabs

.PHONY:all
all: clean modtidy modverify moddownload lint golangcilint test_local build

.PHONY: build
build:
	go build -o bin/$(NAME) enlabs/cmd/.

.PHONY: clean
clean:
	rm -f bin/$(NAME)

.PHONY:modtidy
modtidy:
	go mod tidy

.PHONY: modverify
modverify:
	go mod verify

.PHONY:moddownload
moddownload:
	go mod download

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: golangcilint
golangcilint:
	golangci-lint run --enable-all  -D misspell -D funlen -D wsl --deadline=15m

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: test_local
test_local:
	go test -v -cover -tags local ./...
