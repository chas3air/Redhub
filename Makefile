build-up:
	@docker compose up --build

down:
	@docker compose down

refresh:
	@docker compose down
	@docker compose up --build -d