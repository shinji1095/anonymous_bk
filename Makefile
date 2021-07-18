up:
	docker-compose up -d

down:
	docker-compose down

rs :
	docker-compose restart	

test-user-reg:
	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"name":"Joe","email":"joe@invalid-domain"}'

