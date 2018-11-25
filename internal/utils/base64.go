package utils

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// Base64Type : enum for regexp
type Base64Type string

const (
	// ImageBase64 : enum for regexp
	ImageBase64 Base64Type = "image_base64"
	// FileBase64 : enum for regexp
	FileBase64 Base64Type = "file_base64"
)

var (
	imageRegex = regexp.MustCompile(`^data:image\/[a-z]+;base64,`)
	fileRegex  = regexp.MustCompile(`^data:application\/[a-z]+;base64,`)
)

//DecodeAndSaveBase64 : decode base64 string and save into image; returns filepath, error
func DecodeAndSaveBase64(path string, s string, base64Type Base64Type) (string, error) {
	var re *regexp.Regexp
	switch base64Type {
	case ImageBase64:
		re = imageRegex
	case FileBase64:
		re = fileRegex
	default:
		re = imageRegex
	}
	encodedPhoto := re.ReplaceAllString(s, "")
	semicolon := strings.Index(s, ";")
	slash := strings.Index(s, "/")
	extension := s[slash+1 : semicolon]
	fullPath := path + "." + extension
	unbased, err := base64.StdEncoding.DecodeString(encodedPhoto)
	if err != nil {
		log.Println("Cannot decode base64")
		return "", err
	}

	dirPath := path[:strings.LastIndex(path, "/")]
	if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Println("Cannot mkdir for directory on path: ", dirPath)
	}

	if err = ioutil.WriteFile(fullPath, unbased, 0644); err != nil {
		log.Println("Cannot write file in path: ", fullPath)
		return "", err
	}
	return fullPath, nil
}
