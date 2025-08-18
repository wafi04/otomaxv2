run:
	cd cmd/server && go run main.go

up:
	cd docker && docker compose -f docker-compose-db.yml up -d