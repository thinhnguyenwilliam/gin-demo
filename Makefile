# =========================
# App Info
# =========================
APP_NAME=gin-demo
MAIN_FILE=main.go
PORT=8080

# =========================
# Commands
# =========================

.PHONY: run build clean tidy test air

air:
	air

run:
	go run $(MAIN_FILE)

build:
	go build -o $(APP_NAME) $(MAIN_FILE)

start: build
	./$(APP_NAME)

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -f $(APP_NAME)
