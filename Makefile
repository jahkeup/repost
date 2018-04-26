# This Makefile assumes the environment as described in
# `shell.nix'. Really though, there's likely not a huge difference if
# you were to run with the tool configured without nix.
OUT := out

all: test build

$(OUT):
	mkdir -p $(OUT)

build: $(OUT)
	go build -o $(OUT)/repostd ./cmd/repostd

test:
	go test ./...

generate:
	go generate -v ./...

clean:
	rm -fv $(OUT)
