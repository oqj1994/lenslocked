package M

import (
	"database/sql"
	"fmt"
	"lenslocked/rand"
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	sqlStr := `insert into sessions(user_id,token_hash) values ($1,$2) returning id;`
	token, err := rand.SessionToken()
	if err != nil {
		return nil, err
	}
	session := Session{
		UserID: userID,
		Token:  token,
	}
	//change the Token field to TokenHashField
	row := ss.DB.QueryRow(sqlStr, session.UserID, session.Token)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, err
	}
	//TODO:generate tokenHash by token
	return &session, nil
}

func (ss SessionService) User(token string) (*User, error) {
	var id int
	sqlStr := `select user_id from sessions where token_hash=$1 ;`
	fmt.Println(token)
	row := ss.DB.QueryRow(sqlStr, token)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	user := User{}
	sqlStr = `select id,name,email from users where id=$1 ;`
	row = ss.DB.QueryRow(sqlStr, id)
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
