# gin-demo/Makefile
APP_NAME=gin-demo
MAIN_FILE=main.go
PORT=8085


.PHONY: run build clean tidy test air stress-get stress-post load-test-get

# -z 10s: Run the test for 10 seconds
load-test-get:
	hey -z 10s -c 100 \
	-H "x-api-key: my-super-secret-key" \
	http://localhost:8085/api/v2/news


# 🔥 Stress test GET
# -n 1000 → total 1000 requests
# -c 50 → 50 concurrent users
# -H → add header
stress-get:
	hey -n 1000 -c 50 \
	-H "x-api-key: my-super-secret-key" \
	http://localhost:$(PORT)/api/v2/news

# 🔥 Stress test POST
stress-post:
	hey -n 200 -c 20 \
	-m POST \
	-H "Content-Type: application/json" \
	-H "x-api-key: my-super-secret-key" \
	-D internal/test/body.json \
	http://localhost:$(PORT)/api/v2/products


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
