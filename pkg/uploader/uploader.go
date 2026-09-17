package uploader

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var allowedMIME = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
}

const MaxFileSize = 2 * 1024 * 1024 // 2MB

type Uploader struct {
	dir     string
	baseURL string
}

func New(dir, baseURL string) *Uploader {
	return &Uploader{dir: dir, baseURL: baseURL}
}

// SaveLogo validates MIME type & size, then stores the file under a UUID name.
// Returns the public URL to the stored file.
func (u *Uploader) SaveLogo(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > MaxFileSize {
		return "", fmt.Errorf("ukuran file melebihi 2MB")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := detectContentType(buf[:n])
	if !allowedMIME[contentType] {
		return "", fmt.Errorf("tipe file tidak didukung, gunakan PNG atau JPEG")
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	if err := os.MkdirAll(u.dir, 0o755); err != nil {
		return "", err
	}

	ext := filepath.Ext(fileHeader.Filename)
	fileName := uuid.New().String() + ext
	dstPath := filepath.Join(u.dir, fileName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/uploads/%s", u.baseURL, fileName), nil
}

func detectContentType(buf []byte) string {
	switch {
	case len(buf) > 3 && buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47:
		return "image/png"
	case len(buf) > 2 && buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF:
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}
