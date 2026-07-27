.PHONY: all server web build clean seed

# Start both the API server and web UI dev server
all: server web

# Build the orisd binary
build:
	go build -o bin/orisd ./cmd/orisd

# Run the API server (seeds demo data on first run)
server: build
	./bin/orisd -workspace ./data -addr :8080 -seed 10000 -static ./web/dist &

# Run the web dev server (connects to :8080 API)
web:
	cd web && bun run dev &

# Seed demo data into the demo collection
seed: build
	./bin/orisd -workspace ./data -addr :8081 -seed 100000 &
	@sleep 1
	@kill %1 2>/dev/null || true

# Clean data and build artifacts
clean:
	rm -rf bin/ data/ web/dist/
