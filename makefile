up-build:
	docker  compose -f docker-compose.yml up --build
build:
	docker  compose -f docker-compose.yml build
up:
	docker  compose -f docker-compose.yml up 
up-db:
	docker  compose -f docker-compose.yml up -d db
down:
	docker  compose -f docker-compose.yml down -v