package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"os"
	"os/signal"
	"soa_hw_4/internal/queue"
	"soa_hw_4/internal/tasks"
	"soa_hw_4/internal/users"
	"syscall"
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

	consumer, err := queue.StartTaskConsumer(conn, tasks.NewTasksDB(db), users.NewUsersDB(db))
	if err != nil {
		log.Fatalf("Starting consuming error: %s", err)
	}
	log.Printf("Consumer started")
	defer consumer.Close()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Printf("Exit by signal...")
}
