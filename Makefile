db_container = anonymous_bk_db_1

up:
	docker-compose up -d

down: 
	docker-compose down

rs :
	docker-compose restart	

exec-db:
	docker exec -it $(db_container) bash

migrate:
	curl -X GET http://localhost:8080/migrate/up

migrate-rollback:
	curl -X GET http://localhost:8080/migrate/down

test-user-reg:
	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"eto","lastname":"shinji", "password":"ppp", "email":"joe@invalid-domain"}'

pass=ppp
test-user-val:
	curl -X POST http://localhost:8080/validate/user \
	-H 'Content-Type: application/json' \
  	-d '{"password":"$(pass)", "email":"joe@invalid-domain"}'


