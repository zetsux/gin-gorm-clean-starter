tidy:
	go mod tidy

run: 
	go run main.go

build: 
	go build -o main main.go

run-build: build
	./main