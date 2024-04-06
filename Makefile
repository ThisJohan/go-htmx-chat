run:
	go run cmd/server/server.go

build:
	go build -o main.exe cmd/server/server.go

tail:
	npx tailwindcss -i ./tailwind/styles.css -o ./assets/styles.css --watch

templ:
	templ generate --watch

migrate_up:
	migrate -path migration/ -database "postgres://Johan:junglebook@localhost:5432/chat?sslmode=disable" -verbose up

migrate_down:
	migrate -path migration/ -database "postgres://Johan:junglebook@localhost:5432/chat?sslmode=disable" -verbose down