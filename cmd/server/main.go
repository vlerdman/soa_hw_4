package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"net/http"
	"soa_hw_4/internal/queue"
	"soa_hw_4/internal/server"
	"soa_hw_4/internal/tasks"
	"soa_hw_4/internal/users"
)

func main() {
	db, err := sqlx.Open("pgx", "host=postgres port=5432 user=user password=userpass dbname=userdb")
	if err != nil {
		log.Fatalf("DB open error: %s", err)
	}

	conn, err := rabbitmq.NewConn("amqp://user:userpass@rabbitmq:5672/")

	if err != nil {
		log.Fatalf("RQ open error: %s", err)
	}

	tp, err := queue.NewTaskPublisher(conn)

	if err != nil {
		log.Fatalf("RQ open error: %s", err)
	}

	log.Printf("Server started")

	r := server.NewRouter(tasks.NewTasksDB(db), users.NewUsersDB(db), tp)
	http.ListenAndServe("0.0.0.0:8080", r)
}
