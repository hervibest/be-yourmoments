package adapter

import (
	"be-yourmoments/upload-svc/internal/helper/utils"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
)

type CompressAdapter interface {
	CompressImage(uploadFile multipart.File, dirname string) (string, string, error)
}

type compressAdapter struct {
	compressQuality int
}

func NewCompressAdapter() CompressAdapter {
	compressQuality, _ := strconv.Atoi(utils.GetEnv("COMPRESS_QUALITY")) //75

	return &compressAdapter{
		compressQuality: compressQuality,
	}
}

func (a *compressAdapter) CompressImage(uploadFile multipart.File, dirname string) (string, string, error) {
	log.Println("open in compress image")

	buffer, err := io.ReadAll(uploadFile)
	if err != nil {
		log.Println("error in compress image")
		return "", "", err
	}

	filename := strings.Replace(uuid.New().String(), "-", "", -1) + ".jpg"

	options := bimg.Options{
		Quality: a.compressQuality,
	}

	processed, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return filename, "", err
	}

	filePath := fmt.Sprintf("./%s/%s", dirname, filename)

	// Pastikan direktori ada
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		log.Println("error creating directory:", err)
		return filename, "", err
	}

	if err := bimg.Write(filePath, processed); err != nil {
		log.Println("error when writing")
		return filename, "", err
	}

	return filename, filePath, nil
}
