package utils

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^data:image\/[a-z]+;base64,`)

//DecodeAndSaveBase64Image : decode base64 string and save into image
func DecodeAndSaveBase64Image(path string, s string) error {
	encodedPhoto := re.ReplaceAllString(s, "")
	semicolon := strings.Index(s, ";")
	extension := s[11:semicolon]
	fullPath := path + "." + extension
	unbased, err := base64.StdEncoding.DecodeString(encodedPhoto)
	if err != nil {
		log.Println("Cannot decode base64")
		return err
	}

	dirPath := path[:strings.LastIndex(path, "/")]
	if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Println("Cannot mkdir for directory on path: ", dirPath)
	}

	if err = ioutil.WriteFile(fullPath, unbased, 0644); err != nil {
		log.Println("Cannot write file in path: ", fullPath)
		return err
	}
	return nil
}
