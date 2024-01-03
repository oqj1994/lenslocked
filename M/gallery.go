package M

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

type Gallery struct {
	ID         int
	Title      string
	UserID     int
	Desciption string
}

type Image struct {
	FileName  string
	GalleryID int
}

type GalleryService struct {
	ImageDir string
	DB       *sql.DB
}

func (gs GalleryService) Create(title, desciption string, userID int) (*Gallery, error) {
	//TODO: volidation in here
	gallery := Gallery{
		Title:      title,
		Desciption: desciption,
		UserID:     userID,
	}
	row := gs.DB.QueryRow(`insert into gallerys(title,desciption,user_id) values($1,$2,$3) returning id ;`, title, desciption, userID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery : %w", err)
	}
	return &gallery, nil
}
func (gs GalleryService) ByID(galleryID int) (*Gallery, error) {
	g := Gallery{ID: galleryID}
	row := gs.DB.QueryRow(`select title,desciption,user_id from gallerys where id =$1 ;`, galleryID)
	err := row.Scan(&g.Title, &g.Desciption, &g.UserID)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
func (gs GalleryService) GalleryDir(galleryID int) string {
	imageDir := gs.ImageDir
	if imageDir == "" {
		imageDir = "images"
	}
	return filepath.Join(imageDir, fmt.Sprintf("gallery-%d", galleryID))
}

func (gs GalleryService) List(userID int) ([]Gallery, error) {
	gallerys := []Gallery{}
	rows, err := gs.DB.Query(`select * from gallerys where user_id =$1 ;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var gallery Gallery
		err = rows.Scan(&gallery.ID, &gallery.Title, &gallery.UserID, &gallery.Desciption)
		if err != nil {
			return nil, err
		}
		gallerys = append(gallerys, gallery)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return gallerys, nil
}

func (gs GalleryService) Images(galleryID int) (imgs []Image, err error) {
	dir := filepath.Join(gs.GalleryDir(galleryID), "*")
	files, err := filepath.Glob(dir)
	fmt.Println(files)
	if err != nil {
		return nil, fmt.Errorf("read file :%w", err)
	}
	for _, file := range files {
		if hasExtension(file, gs.extensions()) {

			imgs = append(imgs, Image{FileName: filepath.Base(file), GalleryID: galleryID})
		}
	}
	fmt.Println(imgs)
	return imgs, nil
}

func hasExtension(fileName string, extension []string) bool {
	fileNameLower := strings.ToLower(fileName)
	for _, ext := range extension {
		ext = strings.ToLower(ext)
		if filepath.Ext(fileNameLower) == ext {
			return true
		}
	}
	return false
}

func (gs GalleryService) Update(title, desciption string, id int) error {
	_, err := gs.DB.Exec(`update gallerys set title=$1,desciption =$2 where id=$3`, title, desciption, id)
	if err != nil {
		return err
	}
	return nil
}
func (gs GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(`delete from gallerys where id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (gs GalleryService) UserHave(userID, galleryID int) error {
	row := gs.DB.QueryRow(`select id from gallerys where user_id=$1 and id=$2 ;`, userID, galleryID)
	err := row.Scan(&galleryID)
	if err != nil {
		return err
	}
	return nil
}

func (gs GalleryService) extensions() []string {
	return []string{
		".jpg", ".png", ".gif", ".jpeg",
	}
}
