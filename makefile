dev:
	go test -cover ./...
	go run httpd/main.go