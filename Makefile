build:
	cd frontend && npm install && npm run build
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

dev:
	# Requires air and concurrently to be installed
	npx concurrently "cd frontend && npm run dev" "air -c .air.toml"

clean:
	rm -rf bin tmp ./workspaces
	rm -rf frontend/dist
	rm -f ./data/agentforge.db

migrate:
	# Migrations happen on application startup natively

test:
	go test ./...
