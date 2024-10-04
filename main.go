package main

import (
	"flag"
	"jwt_auth.com/api"
	"jwt_auth.com/internal/mail"
	"jwt_auth.com/pkg/srorage/db/postgres"
	redisConn "jwt_auth.com/pkg/srorage/redis"
	"jwt_auth.com/pkg/web"
)

func main() {
	addr := flag.String("web_addr", ":8000", "Сетевой адрес HTTP")

	dbHost := flag.String("db_host", "db", "Сетевой хост БД")
	dbPort := flag.Int("db_port", 5432, "Сетевой порт БД")
	dbUser := flag.String("db_user", "vladislav", "Пользователь БД")
	dbPassword := flag.String("db_password", "examplePass142", "Пароль пользователя БД")
	dbName := flag.String("db_name", "auth_db", "Название БД")

	redisAddr := flag.String("redis_addr", "redis:6379", "Сетевой адрес Redis")
	redisPass := flag.String("redis_password", "", "Пароль Redis")
	redisDb := flag.Int("redis_db", 0, "Номер БД Redis")

	flag.Parse()

	postgres.InitConn(*dbPort, *dbHost, *dbUser, *dbPassword, *dbName)
	defer postgres.Connection.Close()

	redisConn.NewRedisConnection(*redisAddr, *redisPass, *redisDb)

	mail.NewMailConnection()

	mux := api.RegisterMux()

	srv := web.NewServerConnection(*addr, mux)

	err := srv.ListenAndServe()
	panic(err)
}
