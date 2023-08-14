run:
	docker build -t word-of-wisdom .
	docker run --network host --name word-of-wisdom-container word-of-wisdom

start:
	go run main.go

test:
	go test ./...
	go test -cover ./...

clean:
	@if docker ps -a --filter "name=word-of-wisdom-container" -q | grep -q . ; then \
            docker stop word-of-wisdom-container || true; \
            docker rm word-of-wisdom-container || true; \
            echo "Container 'word-of-wisdom-container' stopped and removed."; \
        else \
            echo "Container 'word-of-wisdom-container' does not exist."; \
        fi

.PHONY: run start test clean
