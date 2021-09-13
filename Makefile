.PHONY:build
build:
	@go build -o TechTrainApi cmd/main.go

.PHONY:run
run:
	@go run cmd/main.go
