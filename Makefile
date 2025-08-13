.PHONY: build up down migrate logs

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

migrate:
	docker compose run --rm backend npx prisma migrate deploy --schema=prisma/schema.prisma

logs:
	docker compose logs -f
