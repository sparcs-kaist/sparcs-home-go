package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sparcs-home-go/internal/app/configure"
	"github.com/sparcs-home-go/internal/utils"
)

// UploadSeminarRequest : upload seminar request body
type UploadSeminarRequest struct {
	Title   string `json:"title"`
	Speaker string `json:"speaker"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

// GetSeminarSchema : get seminar
type GetSeminarSchema struct {
	ID      int    `db:"id"`
	Title   string `db:"title"`
	Speaker string `db:"speaker"`
	Date    string `db:"date"`
	Path    string `db:"path"`
}

// GetSeminarResponse : get seminar
type GetSeminarResponse struct {
	Title   string   `json:"title"`
	Speaker string   `json:"speaker"`
	Date    string   `json:"date"`
	Sources []string `json:"sources"`
}

// GetSeminar : GET /seminar
func GetSeminar() ([]GetSeminarResponse, error) {
	var err error
	seminarList := []GetSeminarSchema{}
	querySeminar := `SELECT Seminar.id, title, speaker, date, path FROM Seminar JOIN SeminarResource S on Seminar.id = S.seminar_id`
	if err = configure.SQL.Select(&seminarList, querySeminar); err != nil {
		log.Println("Failed on GetSeminar query")
		return nil, err
	}
	seminarMap := map[int]GetSeminarResponse{}
	for _, seminar := range seminarList {
		if s, ok := seminarMap[seminar.ID]; ok {
			s.Sources = append(s.Sources, seminar.Path)
			seminarMap[seminar.ID] = s
		} else {
			newSeminar := GetSeminarResponse{
				seminar.Title,
				seminar.Speaker,
				seminar.Date,
				[]string{seminar.Path},
			}
			seminarMap[seminar.ID] = newSeminar
		}
	}
	seminarReponse := []GetSeminarResponse{}
	for _, seminar := range seminarMap {
		seminarReponse = append(seminarReponse, seminar)
	}
	log.Println("seminarReponse: ", seminarReponse)
	return seminarReponse, err
}

// UploadSeminar : POST /seminar/upload
func UploadSeminar(req UploadSeminarRequest) error {
	var err error
	insertSeminar := `INSERT INTO Seminar (title, speaker, date) VALUE (:title, :speaker, :date)`
	row, err := configure.SQL.NamedExec(insertSeminar, req)
	lastInserted, err := row.LastInsertId()
	if err != nil {
		log.Println("Failed to insert seminar", err)
		return err
	}
	seminarDir := fmt.Sprintf("%s/seminar", configure.AppProperties.StaticFilePath)
	seminarName := req.Speaker + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	seminarPath := seminarDir + string(os.PathSeparator) + seminarName
	if seminarPath, err = utils.DecodeAndSaveBase64(seminarPath, req.Content, utils.FileBase64); err != nil {
		log.Println("Failed to DecodeAndSaveBase64")
		return err
	}
	insertSeminarResource := `INSERT INTO SeminarResource (seminar_id, path) VALUE (?, ?)`
	if _, err := configure.SQL.Exec(insertSeminarResource, lastInserted, seminarPath); err != nil {
		log.Println("Failed to insert seminar resource")
		return err
	}
	return nil
}
