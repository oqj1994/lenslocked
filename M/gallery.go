package M

import (
	"database/sql"
	"fmt"
)

type Gallery struct {
	ID          int
	Title       string
	UserID      int
	Description string
}

type GalleryService struct {
	DB *sql.DB
}


func (gs GalleryService) Create(title , desciption string,userID int)(*Gallery,error)  {
	//TODO: volidation in here
	 gallery :=Gallery{
		Title: title,
		Description: desciption,
		UserID: userID,
	 }
	row:=gs.DB.QueryRow(`insert into gallerys(title,desciption,user_id) values($1,$2,$3) returning id ;`,title,desciption,userID)
	err:=row.Scan(&gallery.ID)
	if err !=nil{
		return nil,fmt.Errorf("create gallery : %w",err)
	}
	return &gallery,nil
}

func (gs GalleryService) List(userID int)([]Gallery,error)  {
	gallerys:=[]Gallery{}
	rows,err:=gs.DB.Query(`select * from gallerys where user_id =$1 ;`,userID)
	if err !=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		var gallery Gallery
		err=rows.Scan(&gallery.ID,&gallery.Title,&gallery.UserID,&gallery.Description)
		if err !=nil{
			return nil,err
		}
		gallerys=append(gallerys, gallery)
	}
	if rows.Err() !=nil{
		return nil,rows.Err()
	}
	return gallerys,nil
}

func (gs GalleryService) Update(title, desciption string,id int)error  {
	_,err:=gs.DB.Exec(`update gallerys set title=$1,desciption =$2 where id=$3`,title,desciption,id)
	if err !=nil{
		return err
	}
	return nil
}
func (gs GalleryService) Delete(id int)error  {
	_,err:=gs.DB.Exec(`delete from gallerys where id=$1`,id)
	if err !=nil{
		return err
	}
	return nil
}

func (gs GalleryService) UserHave(userID,galleryID int)error  {
	row:=gs.DB.QueryRow(`select id from gallerys where user_id=$1 and id=$2 ;`,userID,galleryID)
	err:=row.Scan(&galleryID)
	if err !=nil{
		return err
	}
	return nil
}

