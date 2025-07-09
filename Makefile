live:
	make live/templ

# live/tailwind:
	# npx @tailwindcss/cli -i ./views/root.css -o ./static/styles.css --minify --watch

live/templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false --cmd="go run cmd/main.go"

up:
	docker compose up -d

down:
	docker compose down

kill:
	docker compose down -v

