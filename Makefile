test_redis:
	docker run -d --name test_redis -p 6380:6379 redis
	#go mod tidy - прописать, если тесты не работают и ругаются на зависимости
	go test ./pkg/srorage/redis -v
	docker stop test_redis
	docker rm test_redis

test_db:
	docker compose up -d test_db
	go test ./pkg/srorage/db/postgres -v
	docker compose down test_db

test_auth:
	docker compose up -d test_db
	docker run -d --name test_redis -p 6380:6379 redis
	go test ./internal/auth -v
	docker compose down test_db
	docker stop test_redis
	docker rm test_redis

test_api:
	docker compose up -d test_db
	docker run -d --name test_redis -p 6380:6379 redis
	go test ./api -v
	docker compose down test_db
	docker stop test_redis
	docker rm test_redis

test_mail:
	go test ./internal/mail -v
# Для того чтобы быстро выключать рабчие тестовые образы
kill_docker_test:
	docker compose down test_db
	docker stop test_redis
	docker rm test_redis

start_auth:
	docker compose up db redis app --build
