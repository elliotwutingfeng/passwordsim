tests:
	go test -v -race -covermode atomic -coverprofile coverage.out && go tool cover -html coverage.out -o coverage.html

tests_without_race:
	go test -v -covermode atomic -coverprofile coverage.out && go tool cover -html coverage.out -o coverage.html

format:
	go fmt . ./cmd/...

build_cli:
	go build -o ./dist/passwordsim ./cmd/main.go
