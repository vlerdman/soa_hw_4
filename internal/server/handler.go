package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"soa_hw_4/internal/queue"
	"soa_hw_4/internal/server/auth"
	"soa_hw_4/internal/tasks"
	"soa_hw_4/internal/users"
	"strings"
	"time"
)

func NewRouter(tasksDB tasks.TasksDB, usersDB users.UsersDB, publisher *queue.TaskPublisher) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users/register", func(w http.ResponseWriter, req *http.Request) {
		d := json.NewDecoder(req.Body)
		userBody := UserBody{}
		err := d.Decode(&userBody)

		if err != nil {
			http.Error(w, "Incorrect body provided: check and fill all fields", http.StatusBadRequest)
			return
		}

		err = Validate(&userBody)

		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err), http.StatusBadRequest)
			return
		}

		if userBody.Password == "" {
			http.Error(w, "Validation error: password is empty", http.StatusBadRequest)
			return
		}

		user := users.User{
			Id:               uuid.New(),
			RegistrationTime: time.Now(),
			Password:         userBody.Password,
			Username:         userBody.Username,
			Avatar:           userBody.Avatar,
			Sex:              userBody.Sex,
			Email:            userBody.Email,
		}

		err = usersDB.Register(&user)

		if err != nil {
			log.Printf("DB executing error: %s", err)

			http.Error(w, "DB executing error", http.StatusInternalServerError)
			return
		}

		resp := RegistrationResponce{user.Id, auth.GetToken(user.Id)}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	r.HandleFunc("/users/auth", func(w http.ResponseWriter, req *http.Request) {
		d := json.NewDecoder(req.Body)
		userBody := AuthUserBody{}
		err := d.Decode(&userBody)

		if err != nil {
			http.Error(w, "Incorrect body provided: check and fill all fields", http.StatusBadRequest)
			return
		}

		users, err := usersDB.GetByUsernames([]string{userBody.Username})

		if err != nil || len(users) == 0 {
			http.Error(w, "Incorrect body provided: check and fill all fields", http.StatusBadRequest)
		}

		if users[0].Password != userBody.Password {
			http.Error(w, "Incorrect password", http.StatusBadRequest)
		}

		resp := RegistrationResponce{users[0].Id, auth.GetToken(users[0].Id)}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	}).Methods("POST")

	r.HandleFunc("/users/edit", func(w http.ResponseWriter, req *http.Request) {
		d := json.NewDecoder(req.Body)
		userBody := UserBody{}
		err := d.Decode(&userBody)

		if err != nil {
			http.Error(w, "Incorrect body provided: check and fill all fields", http.StatusBadRequest)
			return
		}

		claims, err := auth.FetchToken(w, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err), http.StatusBadRequest)
			return
		}

		user, err := usersDB.GetById(claims.UserId)

		if err != nil {
			log.Printf("DB executing error: %s", err)
			http.Error(w, "DB executing error", http.StatusInternalServerError)
			return
		}

		if userBody.Password != "" {
			user.Password = userBody.Password
		}
		if userBody.Username != "" {
			user.Username = userBody.Username
		}
		if userBody.Sex == "male" || userBody.Sex == "female" {
			user.Sex = userBody.Sex
		}
		user.Avatar = userBody.Avatar
		user.Email = userBody.Email

		err = usersDB.Update(user)

		if err != nil {
			log.Printf("DB executing error: %s", err)
			http.Error(w, "DB executing error", http.StatusInternalServerError)
			return
		}

	}).Methods("POST")

	r.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		claims, err := auth.FetchToken(w, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err), http.StatusBadRequest)
			return
		}

		param := req.URL.Query().Get("usernames")
		usernames := strings.Split(param, ",")

		if param == "" || len(usernames) == 0 {
			user, err := usersDB.GetById(claims.UserId)

			if err != nil {
				http.Error(w, "DB executing error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}

		users, err := usersDB.GetByUsernames(usernames)
		if err != nil {
			log.Printf("DB executing error: %s", err)
			http.Error(w, "DB executing error", http.StatusInternalServerError)
		}

		for _, user := range users {
			user.Password = "******"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&SeveralUserBody{users})

	}).Methods("GET")

	r.HandleFunc("/users/stats", func(w http.ResponseWriter, req *http.Request) {

		claims, err := auth.FetchToken(w, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err), http.StatusBadRequest)
			return
		}

		user, err := usersDB.GetById(claims.UserId)

		if err != nil {
			log.Printf("DB executing error: %s", err)
			http.Error(w, "DB executing error", http.StatusInternalServerError)
			return
		}

		task := tasks.Task{
			Id:           uuid.New(),
			CreationTime: time.Now(),
			UserId:       user.Id,
		}
		err = tasksDB.Create(&task)

		if err != nil {
			log.Printf("DB executing error: %s", err)
			http.Error(w, "DB executing error", http.StatusInternalServerError)
			return
		}

		err = publisher.Publish(task)
		if err != nil {
			log.Printf("Publishing error: %s", err)
			http.Error(w, "Publishing error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&CreateStatsTaskResponce{task.Id})

	}).Methods("POST")

	r.HandleFunc("/users/stats/{id}", func(w http.ResponseWriter, req *http.Request) {

		_, err := auth.FetchToken(w, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err), http.StatusBadRequest)
			return
		}

		vars := mux.Vars(req)

		id, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(w, "Invalid id is provided", http.StatusBadRequest)
			return
		}

		task, err := tasksDB.GetById(id)

		if err != nil {
			log.Printf("DB executing error: %s", err)
			http.Error(w, "DB executing error", http.StatusInternalServerError)
			return
		}

		if len(task.Result) == 0 {
			http.Error(w, "Stats task isn't finished", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(task.Result)

	}).Methods("GET")

	return r
}
