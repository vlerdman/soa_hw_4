package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"math/rand"
	"soa_hw_4/internal/tasks"
	"soa_hw_4/internal/users"
)

func StartTaskConsumer(conn *rabbitmq.Conn, tasksDb tasks.TasksDB, usersDb users.UsersDB) (*rabbitmq.Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			task := tasks.Task{}
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Printf("Consuming error: %s", err)
				return rabbitmq.NackDiscard
			}
			err = processTask(tasksDb, usersDb, task)
			if err != nil {
				log.Printf("Processing error: %s", err)
				return rabbitmq.NackRequeue
			}
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		"tasks",
		rabbitmq.WithConsumerOptionsRoutingKey("ROUTE"),
		rabbitmq.WithConsumerOptionsExchangeName("tasks"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	return consumer, err
}

func processTask(tasksDb tasks.TasksDB, usersDb users.UsersDB, task tasks.Task) error {
	user, err := usersDb.GetById(task.UserId)
	if err != nil {
		return err
	}

	avatarStr := fmt.Sprintf("<left>Avatar: <a href=\"%s\">%s</a></left>", user.Avatar)
	if user.Avatar == "" {
		avatarStr = fmt.Sprintf("<left>Avatar: </left>", user.Avatar)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	_, lineHt := pdf.GetFontSize()

	wins := rand.Intn(1000) + 1
	loses := rand.Intn(1000) + 1
	time := wins + loses + rand.Intn(3)
	htmlStr := "<center><b>Mafia user activity summary</b></center>" +
		fmt.Sprintf("<left>Username: %s</left>", user.Username) +
		avatarStr +
		fmt.Sprintf("<left>Sex: %s</left>", user.Sex) +
		fmt.Sprintf("<left>E-mail: %s</left>", user.Email) +
		fmt.Sprintf("<left>Wins: %d</left>", wins) +
		fmt.Sprintf("<left>Loses: %d</left>", loses) +
		fmt.Sprintf("<left>Total games: %d</left>", wins+loses) +
		fmt.Sprintf("<left>Win rate: %.2f</left>", float64(wins)/float64(wins+loses)) +
		fmt.Sprintf("<left>Total time in game: %d h</left>", time)

	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)

	var b bytes.Buffer
	pdf.Output(&b)
	buf := make([]byte, b.Len())
	b.Read(buf)

	task.Result = buf
	tasksDb.Update(&task)

	return nil
}
