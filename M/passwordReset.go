package M

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"lenslocked/rand"
	"time"
)

const (
	DefaultDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiredAt time.Time
}

type PasswordResetService struct {
	BytesPerToken int
	DB            *sql.DB
	Duration      time.Duration
}

func (prs PasswordResetService) Create(email string) (*PasswordReset, error) {
	var pr PasswordReset
	token, err := rand.ResetSetPasswordToken(prs.BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("generate token :%w", err)
	}
	pr.Token = token
	//Get user by email
	row := prs.DB.QueryRow(`select id from users where email=$1`, email)
	err = row.Scan(&pr.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("can't find user by emial: %w", err)
		}
		return nil, err
	}

	duration := prs.Duration
	if duration == 0 {
		duration = DefaultDuration
	}
	pr.TokenHash = prs.hash(token)
	pr.ExpiredAt = time.Now().Add(duration)

	row = prs.DB.QueryRow(`insert into password_reset(user_id,token_hash,expired_at) values ($1,$2,$3)
    on conflict (user_id) do update set token_hash=$2,expired_at=$3    returning id`,
		pr.UserID, pr.TokenHash, pr.ExpiredAt)
	err = row.Scan(&pr.ID)
	if err != nil {
		return nil, fmt.Errorf("insert password_reset error:%w", err)
	}

	return &pr, nil
}
func (prs PasswordResetService) Consume(token string) (*User, error) {
	var u User
	var passwordReset PasswordReset
	tokenHash := prs.hash(token)
	row := prs.DB.QueryRow(`select u.id,p.id,p.expired_at from users u join password_reset p on u.id=p.user_id 
where p.token_hash=$1`,
		tokenHash)
	err := row.Scan(&u.ID, &passwordReset.ID, &passwordReset.ExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("consume :%w", err)
	}
	if time.Now().After(passwordReset.ExpiredAt) {
		return nil, fmt.Errorf("token expired :%v", token)
	}
	err = prs.Delete(passwordReset.ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (prs PasswordResetService) Delete(id int) error {
	_, err := prs.DB.Exec(`delete from password_reset where id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete passwordReset :%w", err)
	}
	return nil
}

func (prs PasswordResetService) hash(token string) string {
	tokenHsh := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHsh[:])
}
