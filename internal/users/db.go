package users

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersDB interface {
	Register(*User) error
	GetById(id uuid.UUID) (*User, error)
	GetByUsernames(usernames []string) ([]*User, error)
	Update(*User) error
}

type SQLUsersDB struct {
	db *sqlx.DB
}

func NewUsersDB(db *sqlx.DB) UsersDB {
	return &SQLUsersDB{db}
}

func (u *SQLUsersDB) Register(user *User) error {
	sqlQuery := `INSERT INTO "users" (id, registration_time, password, username, avatar, sex, email) 
				VALUES (:id, :registration_time, :password, :username, :avatar, :sex, :email)`
	_, err := u.db.NamedExec(sqlQuery, user)
	return err
}

func (u *SQLUsersDB) GetById(id uuid.UUID) (*User, error) {
	user := &User{}
	sqlQuery := `SELECT * FROM "users" WHERE id = $1`
	err := u.db.Get(user, sqlQuery, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *SQLUsersDB) GetByUsernames(usernames []string) ([]*User, error) {
	query, args, err := sqlx.In(`SELECT * FROM users WHERE username in (?)`, usernames)
	if err != nil {
		return nil, err
	}

	query = u.db.Rebind(query)
	rows, err := u.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	var users []*User

	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (u *SQLUsersDB) Update(user *User) error {
	_, err := u.db.NamedExec(`UPDATE users 
									SET password = :password, username = :username, avatar = :avatar, sex = :sex, email = :email 
									WHERE id = :id`, user)

	return err
}
