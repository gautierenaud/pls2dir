.PHONY: build
build:
	mkdir -p build
	env GOOS=linux GOARCH=amd64 go build -o build/pls2dir-linux-amd64 .
	env GOOS=linux GOARCH=386 go build -o build/pls2dir-linux-386 .
	env GOOS=darwin GOARCH=arm64 go build -o build/pls2dir-darwin-arm64 .
	env GOOS=darwin GOARCH=amd64 go build -o build/pls2dir-darwin-amd64 .
	env GOOS=windows GOARCH=386 go build -o build/pls2dir-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o build/pls2dir-windows-amd64.exe .