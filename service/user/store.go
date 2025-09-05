package user

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/zinx110/golang-backend-rest/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {

	newStore := &Store{db: db}
	if err := createUserTableIfNotExists(newStore); err != nil {
		log.Fatal("failed to create user table:", err)
	}
	return newStore
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	return u, nil
}

func createUserTableIfNotExists(s *Store) error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY, 
		firstName VARCHAR(100) NOT NULL,
		lastName VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`)
	return err
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (s *Store) CreateUser(user types.User) error {
	return nil
}
