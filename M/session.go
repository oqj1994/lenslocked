package M

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
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
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// sqlStr := `update sessions set token_hash=$2 where user_id=$1 returning id;`
	token, err := rand.SessionToken(ss.BytesPerToken)
	if err != nil {
		return nil, err
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	// //change the Token field to TokenHashField
	// row := ss.DB.QueryRow(sqlStr, session.UserID, session.TokenHash)
	// err = row.Scan(&session.ID)
	// if err == sql.ErrNoRows {
	// 	sqlStr = `insert into sessions(user_id,token_hash) values ($1,$2) returning id;`
	// 	row = ss.DB.QueryRow(sqlStr, session.UserID, session.TokenHash)
	// 	err = row.Scan(&session.ID)
	// }
	sqlStr := `insert into sessions(user_id,token_hash) values($1, $2) 
            on CONFLICT (user_id) DO update set token_hash =$2 returning id`
	row := ss.DB.QueryRow(sqlStr, userID, ss.hash(token))
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, err
	}
	//TODO:generate tokenHash by token
	return &session, nil
}

func (ss SessionService) Delete(userID int) error {
	sqlStr := `delete from sessions where user_id = $1;`
	_, err := ss.DB.Exec(sqlStr, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ss SessionService) User(token string) (*User, error) {
	user := User{}
	
	sqlStr := `select u.id,name,email from users u join sessions s on u.id=s.user_id where s.token_hash=$1;`
	row := ss.DB.QueryRow(sqlStr, ss.hash(token))
	err := row.Scan(&user.ID,&user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("create :%w", err)
	}
	return &user, nil
}

func (ss SessionService) hash(token string) string {
	tokenHsh := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHsh[:])
}
