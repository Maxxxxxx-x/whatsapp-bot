#!make

main_bin_name=bot
main_cmd_path=./cmd/${main_bin_name}


tidy:
	go mod tidy
	go fmt ./...

clean:
	@if [ -d ./tmp ]; then rm -rf ./tmp; fi

build: clean tidy
	go build -o=./tmp/bin/${main_bin_name} ${main_cmd_path}

.PHONY: test-run
test-run: build
	./tmp/bin/${main_bin_name}

.PHONY: test
test: build
	go test ./...
