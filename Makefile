BINARY_NAME := tftpuhpd
BUILD_DIR := build

TARGETS := \
    linux_amd64 \
    darwin_amd64 \
    windows_amd64

.PHONY: all clean deadcode $(TARGETS)

all: $(TARGETS)

$(TARGETS):
	GOOS=$(word 1,$(subst _, ,$@)) \
	GOARCH=$(word 2,$(subst _, ,$@)) \
	go build -o $(BUILD_DIR)/$@/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

check:
	go test ./cmd/... ./pkg/...
	go vet ./cmd/... ./pkg/...

deadcode:
	~/go/bin/deadcode ./cmd/...

clean:
	go clean
	rm -rf $(BUILD_DIR)

