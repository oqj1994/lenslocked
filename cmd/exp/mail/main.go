package main

import (
	"lenslocked/M"
)

func main() {
	// m:=gomail.NewMessage()
	// m.SetHeader("From","fun@club.com")
	// m.SetHeader("To","fun@club.com")
	// m.SetHeader("subject","注册成功")
	// m.SetHeader("To","fun@club.com")
	// m.SetBody("text/html",`<h1>Hello</h1>`)
	// m.SetBody("text/plain",`I love cat`)
	// // m.AddAlternative("text/html",``)
	// d:=gomail.NewDialer("sandbox.smtp.mailtrap.io",25,"7eb4b9ffde3ef0","b926a640148933")
	// m.WriteTo(os.Stdout)
	// sender,err:=d.Dial()
	// defer sender.Close()
	// if err !=nil {
	// 	panic(err)
	// }
	// to:=[]string{"oqj@163.com","greenland@world.com"}
	// err=sender.Send("jia@qq.com",to,m)
	// if err !=nil{
	// 	fmt.Println("error : ",err)
	// }
	// email:=M.Email{
	// 	To:      "cat@io.com",
	// 	From:    "",
	// 	Subject: "Welcome to SongTang country",
	// 	Text:    "raw text",
	// 	HTML:    "<h1>Welcome!!</h1>",
	// }

	es := M.NewEmailService(M.SMTPConfig{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     25,
		UserName: "7b4e8c97b7bebc",
		Password: "2136315cbe2186",
	})
	// err:=es.Send(email)
	err := es.Send(M.Email{
		To:      "oqj@163.com",
		Subject: "Hello World",
		Text:    "this is a mail test",
		HTML:    "<h1>Hello </h1>",
	})
	if err != nil {
		panic(err)
	}
}
