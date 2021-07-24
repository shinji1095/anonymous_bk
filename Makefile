db_container = anonymous_bk_db_1
env = dev
pro_url = https://anonymous-bk.herokuapp.com/
dev_url = http://localhost:8080/

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
	@make seed-user-all
	@make seed-do-all
	@make seed-ass
	@make seed-group

seed-user:
	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"eto","lastname":"shinji", "password":"ppp", "email":"aaa@mail.com", "groupID":1}'

seed-user-all:
	@make seed-user
	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"yoshida","lastname":"takanori", "password":"ppp", "email":"bbb@mail.com", "groupID":1}'

	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"tanaka","lastname":"ryouya", "password":"ppp", "email":"ccc@mail.com", "groupID":1}'

	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"hosida","lastname":"kyousuke", "password":"ppp", "email":"ddd@mail.com", "groupID":1}'

	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"ogata","lastname":"miduho", "password":"ppp", "email":"eee@mail.com", "groupID":1}'

seed-group:
	curl -X POST http://localhost:8080/group \
	-H 'Content-Type: application/json' \
  	-d '{"name":"hibikino"}'

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

seed-do:
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":1, "status":2, "ranking":1, "updateAt":"2021-07-02T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":2, "status":$(status), "ranking":2, "updateAt":"2021-07-09T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":3, "status":2, "ranking":3, "updateAt":"2021-07-10T17:00:00+09:00"}'	
	
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":1, "status":$(status), "ranking":1, "updateAt":"2021-07-12T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":2, "status":2, "ranking":4, "updateAt":"2021-07-17T17:00:00+09:00"}'

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":3, "status":$(status), "ranking":1, "updateAt":"2021-07-20T17:00:00+09:00"}'	

seed-do-all:
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":0, "ranking":1, "updateAt":"2021-07-02T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":2, "status":0, "ranking":2, "updateAt":"2021-07-09T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":3, "status":0, "ranking":3, "updateAt":"2021-07-10T17:00:00+09:00"}'	
	
	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":0, "ranking":1, "updateAt":"2021-07-12T17:00:00+09:00"}'	

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":2, "status":0, "ranking":4, "updateAt":"2021-07-17T17:00:00+09:00"}'

	curl -X POST http://localhost:8080/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":3, "status":2, "ranking":1, "updateAt":"2021-07-20T17:00:00+09:00"}'	
	@make seed-do userID=2 status=1
	@make seed-do userID=3 status=2
	@make seed-do userID=4 status=2
	@make seed-do userID=5 status=0

	
seed-refresh:
	@make migrate-refresh
	@make seed-all


test-user-reg:
	curl -X POST http://localhost:8080/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"eto","lastname":"shinji", "password":"ppp", "email":"joe@invalid-domain", "groupID":1}'

pass=ppp
test-user-val:
	curl -X POST http://localhost:8080/validate/user \
	-H 'Content-Type: application/json' \
  	-d '{"password":"$(pass)", "email":"joe@invalid-domain"}'

test-ass-reg:
	curl -X POST http://localhost:8080/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math", "due":"monday"}'

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
	curl -X GET http://localhost:8080/do/week/$(userID)

groupID=1
test-ass:
	@make test-group-reg
	curl -X GET http://localhost:8080/assignment/$(groupID)

status=1
test-do-put:
	curl -X PUT http://localhost:8080/do?userID=1\&assignmentID=1\&status=$(status)

test-user-put:
	curl -X PUT http://localhost:8080/user/$(userID)?groupID=$(groupID)

userID=2
test-user-get:
	curl -X GET http://localhost:8080/user/$(userID)

groupID=1
test-group-get:
	curl -X GET http://localhost:8080/group/$(groupID)

test-group-get-all:
	curl -X GET http://localhost:8080/group

test-belong:
	curl -X PUT http://localhost:8080/belong/group/$(groupID)?userID=$(userID)

test-belong-all:
	@make migrate-refresh
	@make seed-user-all
	@make seed-group
	@make test-belong groupID=1 userID=1
	@make test-belong groupID=1 userID=2
	@make test-belong groupID=1 userID=3
	@make test-belong groupID=1 userID=4
	@make test-belong groupID=1 userID=5