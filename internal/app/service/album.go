package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sparcs-home-go/internal/utils"

	"github.com/sparcs-home-go/internal/app/configure"
)

// YearSchema : year information
type YearSchema struct {
	Year        int `db:"year" json:"year"`
	EventNumber int `db:"event_number" json:"eventNumber"`
	PhotoNumber int `db:"photo_number" json:"photoNumber"`
}

// PhotoSchema : photo info
type PhotoSchema struct {
	AlbumID    int    `db:"id"`
	AlbumYear  int    `db:"year"`
	AlbumTitle string `db:"title"`
	AlbumDate  string `db:"date"`
	Path       string `db:"path"`
}

// AlbumSchema : album information
type AlbumSchema struct {
	ID          int      `json:"id"`
	Year        int      `json:"year"`
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	PhotoNumber int      `json:"photoNumber"`
	Photos      []string `json:"photos"`
}

// UploadPhotoRequest :
type UploadPhotoRequest struct {
	AlbumID   int      `json:"albumId"`
	PhotoList []string `json:"photoList"`
}

// CreateAlbumRequest :
type CreateAlbumRequest struct {
	Year       int    `json:"year"`
	AlbumTitle string `json:"albumTitle"`
	AlbumDate  string `json:"albumDate"` // YYYY-MM-DD
}

// ListAlbum : List albums
func ListAlbum() ([]YearSchema, []AlbumSchema, error) {
	var err error
	yearList := []YearSchema{}
	if err = configure.SQL.Select(&yearList, `
		SELECT a.year as year, COUNT(DISTINCT a.id) as event_number, COUNT(p.id) as photo_number
		FROM Album AS a LEFT JOIN Photo AS p
		ON a.id = p.album_id 
		GROUP BY a.year
	`); err != nil {
		log.Println("Failed on YearSchema")
	}

	log.Println("yearList: ", yearList)

	photoList := []PhotoSchema{}
	if err = configure.SQL.Select(&photoList, `
		SELECT a.id as id, a.year as year, a.title as title, a.date as date, p.path as path
		FROM Album AS a LEFT JOIN Photo AS p 
		ON a.id = p.album_id
	`); err != nil {
		log.Println("Failed on PhotoSchema")
	}

	log.Println("photoList: ", photoList)

	albumMap := map[int]AlbumSchema{}
	for _, photo := range photoList {
		if album, ok := albumMap[photo.AlbumID]; ok {
			album.PhotoNumber++
			album.Photos = append(album.Photos, photo.Path)
			albumMap[photo.AlbumID] = album
		} else {
			newAlbum := AlbumSchema{
				photo.AlbumID,
				photo.AlbumYear,
				photo.AlbumTitle,
				photo.AlbumDate,
				1,
				[]string{photo.Path},
			}
			albumMap[newAlbum.ID] = newAlbum
		}
	}
	albumList := []AlbumSchema{}
	for _, album := range albumMap {
		albumList = append(albumList, album)
	}
	log.Println("albumList: ", albumList)
	return yearList, albumList, err
}

// UploadPhoto : upload photo
func UploadPhoto(req UploadPhotoRequest) error {
	var err error
	row := configure.SQL.QueryRow(`
		SELECT year, title FROM Album WHERE id = ?
	`, req.AlbumID)
	var (
		year  int
		title string
	)
	if err = row.Scan(&year, &title); err != nil {
		log.Println("Error while scanning row")
		return err
	}
	albumPath := fmt.Sprintf("%s/%d/%s", configure.AppProperties.StaticFilePath, year, title)
	insertPhoto := `INSERT INTO Photo (album_id, path) VALUES (?, ?)`
	for _, photoBase64 := range req.PhotoList {
		photoName := strconv.FormatInt(time.Now().Unix(), 10)
		photoPath := albumPath + string(os.PathSeparator) + photoName
		// TODO: Transaction management
		if err = utils.DecodeAndSaveBase64Image(photoPath, photoBase64); err != nil {
			log.Println("Failed on DecodeAndSaveBase64Image")
			return err
		}
		if _, err = configure.SQL.Query(insertPhoto, req.AlbumID, photoPath); err != nil {
			log.Println("Failed on inserting photo")
			return err
		}
	}
	return nil
}

// CreateAlbum : create album
func CreateAlbum(req CreateAlbumRequest) error {
	insertAlbum := `INSERT INTO Album (year, title, date) VALUES (?, ?, ?)`
	if _, err := configure.SQL.Query(insertAlbum, req.Year, req.AlbumTitle, req.AlbumDate); err != nil {
		log.Println("Failed on inserting album")
		return err
	}
	return nil
}

// DeleteAlbum : delete album
func DeleteAlbum(albumID int) error {
	deleteAlbum := `DELETE FROM Album WHERE id = ?`
	if _, err := configure.SQL.Query(deleteAlbum, albumID); err != nil {
		log.Println("Failed to delete album")
		return err
	}
	return nil
}

// DeletePhoto : delete photo
func DeletePhoto(photoID int) error {
	deletePhoto := `DELETE FROM Photo WHERE id = ?`
	if _, err := configure.SQL.Query(deletePhoto, photoID); err != nil {
		log.Println("Failed to delete photo")
		return err
	}
	return nil
}
