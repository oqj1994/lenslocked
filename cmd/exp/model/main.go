package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"lenslocked/M"
)

func main() {
	cfg := M.DefaultConfig()

	//sql.Register("pgx", stdlib.GetDefaultDriver())
	db, err := M.Open(cfg)
	if err != nil {
		panic(err)
	}

	us := M.UserService{DB: db}
	user, err := us.Create(M.CreateUserParms{
		Name:     "jia",
		Email:    "oqj@foo.com",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	//	err = db.Ping()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("database connected!")
	//	_, err = db.Exec(`CREATE TABLE  IF NOT EXISTS users(
	//    id serial primary key ,
	//    email text unique not null
	//);
	//create table if not exists tweets(
	//    id serial primary key ,
	//    owner_id int,
	//    content text,
	//    created_at time,
	//    updated_at time,
	//    deleted_at time
	//);
	//create table if not exists postLikes(
	//    id serial primary key ,
	//    post_id int,
	//    user_id int
	//)
	//`)
	//	if err != nil {
	//		panic(err)
	//	}
	//	result, err := db.Exec(`delete  from users where email='oqj@163.com'`)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(result.RowsAffected())
	//
	//	//insert into users
	//	result, err = db.Exec(`insert into users(email) values ('oqj@163.com')`)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(result.RowsAffected())
	//	result, err = db.Exec(`update users set email=concat("email",'*****.com')`)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(result.RowsAffected())
	//
	//	row := db.QueryRow(`select email from users `)
	//	var email string
	//	err = row.Scan(&email)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("email address is %s", email)
}
