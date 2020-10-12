# Keep test at the top so that it is default when `make` is called.
# This is used by Travis CI.
coverage.txt:
	go test -race -covermode=atomic -coverprofile=coverage.txt
view-cover: clean coverage.txt
	go tool cover -html=coverage.txt
test: build
	go test ./...
build:
	go build ./...
install: build
	go install ./...
inspect: build
	golint ./...
pre-commit: clean coverage.txt inspect
	go mod tidy
clean:
	rm -f coverage.txt