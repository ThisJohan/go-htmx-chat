build:
	go build -o main.exe cmd/server/server.go

tail:
	npx tailwindcss -i ./tailwind/styles.css -o ./assets/styles.css --watch

templ:
	templ generate --watch