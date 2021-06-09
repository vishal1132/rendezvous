test: linttest
	go test ./...

linttest:
	golangci-lint run