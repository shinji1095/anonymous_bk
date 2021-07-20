db_container = anonymous_db_1

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

migrate-refresh:
	@make migrate-rollback
	@make migrate

seed-all:
	@make seed-do-single
	@make seed-ass
	@make seed-group

seed-group:
	@make test-group-reg

seed-ass:
	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math", "due":"monday","groupID":1}'

	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"english", "due":"tuesday","groupID":1}'

	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"statistics", "due":"friday","groupID":1}'

	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"japanese", "due":"monday","groupID":2}'

	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math2", "due":"tuesday","groupID":2}'

	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"statistics", "due":"friday","groupID":2}'

seed-do-single:
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":2, "ranking":1, "updateAt":"2021-07-02T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":2, "status":2, "ranking":2, "updateAt":"2021-07-09T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":3, "status":2, "ranking":3, "updateAt":"2021-07-10T17:00:00+09:00"}'	
	
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":2, "ranking":1, "updateAt":"2021-07-12T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":2, "status":2, "ranking":4, "updateAt":"2021-07-17T17:00:00+09:00"}'

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":3, "status":2, "ranking":1, "updateAt":"2021-07-20T17:00:00+09:00"}'		

	
seed-refresh:
	@make migrate-refresh
	@make seed-all


test-user-reg:
	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"eto","lastname":"shinji", "password":"ppp", "email":"joe@invalid-domain"}'

pass=ppp
test-user-val:
	curl -X POST http://localhost:8080/validate/user \
	-H 'Content-Type: application/json' \
  	-d '{"password":"$(pass)", "email":"joe@invalid-domain"}'

test-ass-reg:
	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math", "due":"2021-07-23T17:00:00+09:00"}'

test-group-reg:
	curl -X POST http://localhost:8080/group \
	-H 'Content-Type: application/json' \
  	-d '{"name":"hibikino"}'

test-share-reg:
	curl -X POST http://localhost:8080/share \
	-H 'Content-Type: application/json' \
  	-d '{"groupID":1,"assignmentID":1}'

test-do-reg:
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":2, "ranking":1, "updateAt":"2021-07-20T17:00:00+09:00"}'		

test-do-get:
	curl -X GET http://localhost:8080/do?userID=1\&year=2021\&month=7 

test-do-get-spec:
	curl -X GET http://localhost:8080/do/userID/1

test-do-week:
	curl -X GET http://localhost:8080/do/week

groupID=1
test-ass:
	@make test-group-reg
	curl -X GET http://localhost:8080/assignment/$(groupID)

test-do-put:
	curl -X PUT http://localhost:8080/do?userID=1\&assignmentID=2\&status=1 
