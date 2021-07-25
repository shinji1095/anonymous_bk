db_container = anonymous_bk_db_1
env = dev
pro_url = https://anonymous-bk.herokuapp.com/
dev_url = https://anonymous-bk.herokuapp.com/

up:
	docker-compose up -d

down: 
	docker-compose down

rs :
	docker-compose restart	

exec-db:
	docker exec -it $(db_container) bash

migrate:
	curl -X GET https://anonymous-bk.herokuapp.com/migrate/up

migrate-rollback:
	curl -X GET https://anonymous-bk.herokuapp.com/migrate/down

migrate-refresh:
	@make migrate-rollback
	@make migrate

seed-all:
	@make seed-group
	@make seed-user-all
	@make seed-ass

seed-user:
	curl -X POST https://anonymous-bk.herokuapp.com/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"eto","lastname":"shinji", "password":"ppp", "email":"aaa@mail.com", "groupID":1}'

seed-user-all:
	@make seed-user
	curl -X POST https://anonymous-bk.herokuapp.com/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"yoshida","lastname":"takanori", "password":"ppp", "email":"bbb@mail.com", "groupID":1}'

	curl -X POST https://anonymous-bk.herokuapp.com/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"tanaka","lastname":"ryouya", "password":"ppp", "email":"ccc@mail.com", "groupID":1}'

	curl -X POST https://anonymous-bk.herokuapp.com/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"hosida","lastname":"kyousuke", "password":"ppp", "email":"ddd@mail.com", "groupID":1}'

	curl -X POST https://anonymous-bk.herokuapp.com/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"ogata","lastname":"miduho", "password":"ppp", "email":"eee@mail.com", "groupID":1}'

seed-group:
	curl -X POST https://anonymous-bk.herokuapp.com/group \
	-H 'Content-Type: application/json' \
  	-d '{"name":"hibikino"}'

seed-ass:
	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math", "due":"monday","groupID":1}'

	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"english", "due":"tuesday","groupID":1}'

	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"statistics", "due":"friday","groupID":1}'

	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"japanese", "due":"monday","groupID":2}'

	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math2", "due":"tuesday","groupID":2}'

	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"statistics", "due":"friday","groupID":2}'

seed-do:
	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":1, "status":2, "ranking":1, "updateAt":"2021-07-02T17:00:00+09:00"}'	

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":2, "status":$(status), "ranking":2, "updateAt":"2021-07-09T17:00:00+09:00"}'	

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":3, "status":2, "ranking":3, "updateAt":"2021-07-10T17:00:00+09:00"}'	
	
	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":1, "status":$(status), "ranking":1, "updateAt":"2021-07-12T17:00:00+09:00"}'	

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":2, "status":2, "ranking":4, "updateAt":"2021-07-17T17:00:00+09:00"}'

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":$(userID),"assignmentID":3, "status":$(status), "ranking":1, "updateAt":"2021-07-20T17:00:00+09:00"}'	

seed-do-all:
	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":0, "ranking":1, "updateAt":"2021-07-02T17:00:00+09:00"}'	

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":2, "status":0, "ranking":2, "updateAt":"2021-07-09T17:00:00+09:00"}'	

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":3, "status":0, "ranking":3, "updateAt":"2021-07-10T17:00:00+09:00"}'	
	
	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":0, "ranking":1, "updateAt":"2021-07-12T17:00:00+09:00"}'	

	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":2, "status":0, "ranking":4, "updateAt":"2021-07-17T17:00:00+09:00"}'

	curl -X POST https://anonymous-bk.herokuapp.com/do \
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
	curl -X POST https://anonymous-bk.herokuapp.com/user \
	-H 'Content-Type: application/json' \
  	-d '{"firstname":"eto","lastname":"shinji", "password":"ppp", "email":"joe@invalid-domain", "groupID":1}'

pass=ppp
test-user-val:
	curl -X POST https://anonymous-bk.herokuapp.com/validate/user \
	-H 'Content-Type: application/json' \
  	-d '{"password":"$(pass)", "email":"joe@invalid-domain"}'

test-ass-reg:
	curl -X POST https://anonymous-bk.herokuapp.com/assignment \
	-H 'Content-Type: application/json' \
  	-d '{"name":"math", "due":"monday"}'

test-group-reg:
	curl -X POST https://anonymous-bk.herokuapp.com/group \
	-H 'Content-Type: application/json' \
  	-d '{"name":"hibikino"}'

test-share-reg:
	curl -X POST https://anonymous-bk.herokuapp.com/share \
	-H 'Content-Type: application/json' \
  	-d '{"groupID":1,"assignmentID":1}'

test-do-reg:
	curl -X POST https://anonymous-bk.herokuapp.com/do \
	-H 'Content-Type: application/json' \
  	-d '{"userID":1,"assignmentID":1, "status":2, "ranking":1, "updateAt":"2021-07-20T17:00:00+09:00"}'		

test-do-get:
	curl -X GET https://anonymous-bk.herokuapp.com/do?userID=1\&year=2021\&month=7 

test-do-get-spec:
	curl -X GET https://anonymous-bk.herokuapp.com/do/userID/1

test-do-week:
	curl -X GET https://anonymous-bk.herokuapp.com/do/week/$(userID)

groupID=1
test-ass:
	@make test-group-reg
	curl -X GET https://anonymous-bk.herokuapp.com/assignment/$(groupID)

status=1
test-do-put:
	curl -X PUT https://anonymous-bk.herokuapp.com/do?userID=1\&assignmentID=1\&status=$(status)

test-user-put:
	curl -X PUT https://anonymous-bk.herokuapp.com/user/$(userID)?groupID=$(groupID)

userID=2
test-user-get:
	curl -X GET https://anonymous-bk.herokuapp.com/user/$(userID)

groupID=1
test-group-get:
	curl -X GET https://anonymous-bk.herokuapp.com/group/$(groupID)

test-group-get-all:
	curl -X GET https://anonymous-bk.herokuapp.com/group

test-belong:
	curl -X PUT https://anonymous-bk.herokuapp.com/belong/group/$(groupID)?userID=$(userID)

test-belong-all:
	@make migrate-refresh
	@make seed-user-all
	@make seed-group
	@make test-belong groupID=1 userID=1
	@make test-belong groupID=1 userID=2
	@make test-belong groupID=1 userID=3
	@make test-belong groupID=1 userID=4
	@make test-belong groupID=1 userID=5