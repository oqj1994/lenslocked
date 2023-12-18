package M

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
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
	sqlStr := `insert into sessions(user_id,token_hash) values ($1,$2) returning id;`
	token, err := rand.SessionToken(ss.BytesPerToken)
	if err != nil {
		return nil, err
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	//change the Token field to TokenHashField
	row := ss.DB.QueryRow(sqlStr, session.UserID, session.TokenHash)
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
	var id int
	sqlStr := `select user_id from sessions where token_hash=$1 ;`
	tokenHash := ss.hash(token)
	row := ss.DB.QueryRow(sqlStr, tokenHash)
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

func (ss SessionService) hash(token string) string {
	tokenHsh := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHsh[:])
}
