test: linttest
	go test ./...

linttest:
	golangci-lint -v run