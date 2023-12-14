package M

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
}

type NewUser struct {
	Name     string
	Email    string
	Password string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(user NewUser) (*User, error) {
	user.Email = strings.ToLower(user.Email)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generate passwordHash  error:%w", err)
	}
	sqlStr := `insert into users(name,email,password_hash)
values ($1,$2,$3) returning id;`

	u := User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: string(passwordHash),
	}
	row := us.DB.QueryRow(sqlStr, user.Name, user.Email, passwordHash)

	err = row.Scan(&u.ID)
	if err != nil {
		return nil, fmt.Errorf("insert into DB error:%w", err)
	}
	return &u, nil

}
