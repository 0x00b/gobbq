.PHONY: lintAll vet fmt lint note imports

PROTOC_GIT_TAG=v1.2.0
GOLANGCI_LINT_TAG=v1.33.0
WORK_PATH=$(shell cd .. && pwd)
 
#提交之前建议 make lintAll
lintAll:
	$(MAKE) fmt
	$(MAKE) vet
#	$(MAKE) note
	$(MAKE) lint

downloadlint = $(shell command -v golangci-lint>/dev/null && echo 0 || echo 1)
downloadgonote = $(shell command -v gonote>/dev/null && echo 0 || echo 1)
downloadgoimports = $(shell command -v goimports>/dev/null && echo 0 || echo 1)
 
install_protoc:
	-go get -d -u github.com/golang/protobuf/protoc-gen-go
	git -C "$(shell go env GOPATH)"/src/github.com/golang/protobuf checkout ${PROTOC_GIT_TAG}
	go install github.com/golang/protobuf/protoc-gen-go

install_tools:
ifeq ($(downloadgoimports),1)
	go get golang.org/x/tools/cmd/goimports
endif
ifeq ($(downloadlint),1)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_TAG}
endif
#	下载.golangci.yml
ifeq ($(wildcard .golangci.yml),)
	wget https://lint-1257095926.cos.ap-guangzhou.myqcloud.com/.golangci.yml
endif
ifeq ($(downloadgonote),1)
	go get git.code.oa.com/vlib/tools/gonote
	gonote -u
# 	新项目建议开启
#	gonote -p
endif
#	自动注册pre-commit，提交前自动检查代码
ifneq ($(shell grep -c "make lintAll" .git/hooks/pre-commit), 1)
	echo 'make lintAll' >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
endif


format:
	$(MAKE) fmt
	$(MAKE) vet

lint: install_tools
	golangci-lint run

note: install_tools
	gonote ./...

fmt:
	go fmt ./...

imports: install_tools
	goimports -w ./

vet:
	go vet ./...