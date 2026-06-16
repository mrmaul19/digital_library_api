package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Uploader struct {
	uploadDir   string
	maxFileSize int64
}

func NewUploader(uploadDir string, maxFileSize int64) *Uploader {
	os.MkdirAll(uploadDir+"/covers", os.ModePerm)
	os.MkdirAll(uploadDir+"/books", os.ModePerm)
	return &Uploader{uploadDir: uploadDir, maxFileSize: maxFileSize}
}

var allowedImageTypes = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

var allowedBookTypes = map[string]bool{
	".pdf": true, ".epub": true,
}

func (u *Uploader) UploadCover(c *fiber.Ctx, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", nil
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageTypes[ext] {
		return "", fmt.Errorf("format gambar tidak didukung, gunakan jpg/jpeg/png/webp")
	}

	if file.Size > u.maxFileSize {
		return "", fmt.Errorf("ukuran file melebihi batas %d bytes", u.maxFileSize)
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	savePath := filepath.Join(u.uploadDir, "covers", filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return "", fmt.Errorf("gagal menyimpan cover: %w", err)
	}

	return "/uploads/covers/" + filename, nil
}

func (u *Uploader) UploadBook(c *fiber.Ctx, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", nil
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedBookTypes[ext] {
		return "", fmt.Errorf("format file tidak didukung, gunakan pdf/epub")
	}

	if file.Size > u.maxFileSize {
		return "", fmt.Errorf("ukuran file melebihi batas %d bytes", u.maxFileSize)
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	savePath := filepath.Join(u.uploadDir, "books", filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return "", fmt.Errorf("gagal menyimpan file buku: %w", err)
	}

	return "/uploads/books/" + filename, nil
}