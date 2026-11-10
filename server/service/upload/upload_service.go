package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Uploader interface {
	Upload(file multipart.File, header *multipart.FileHeader, subDir string) (string, error)
}

type LocalUploader struct {
	baseDir      string
	maxSize      int64
	allowedTypes map[string]bool
}

func NewLocalUploader(baseDir string, maxSizeMB int64, allowedTypes []string) *LocalUploader {
	at := make(map[string]bool)
	for _, t := range allowedTypes {
		at[t] = true
	}
	return &LocalUploader{
		baseDir:      baseDir,
		maxSize:      maxSizeMB * 1024 * 1024,
		allowedTypes: at,
	}
}

func (u *LocalUploader) Validate(file multipart.File, header *multipart.FileHeader) error {
	if header.Size > u.maxSize {
		return fmt.Errorf("文件大小超出限制")
	}

	magicBuf := make([]byte, 512)
	n, err := file.Read(magicBuf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("读取文件失败")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("读取文件失败")
	}
	if !validateMagicNumber(magicBuf[:n]) {
		return fmt.Errorf("文件内容校验失败，仅支持 JPEG/PNG/WebP 格式")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	mimeType := header.Header.Get("Content-Type")
	allowed := false
	for t := range u.allowedTypes {
		if mimeType == t {
			allowed = true
			break
		}
	}
	if !allowed {
		switch ext {
		case ".jpg", ".jpeg":
			allowed = u.allowedTypes["image/jpeg"]
		case ".png":
			allowed = u.allowedTypes["image/png"]
		case ".webp":
			allowed = u.allowedTypes["image/webp"]
		}
	}
	if !allowed {
		return fmt.Errorf("不支持的文件类型")
	}

	return nil
}

func validateMagicNumber(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return true
	}

	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return true
	}

	if len(data) >= 12 &&
		data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return true
	}

	return false
}

func (u *LocalUploader) Upload(file multipart.File, header *multipart.FileHeader, subDir string) (string, error) {
	if err := u.Validate(file, header); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	newName := uuid.New().String() + ext

	targetDir := filepath.Join(u.baseDir, subDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	targetPath := filepath.Join(targetDir, newName)
	dst, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		absPath = targetPath
	}
	logrus.Infof("文件上传成功: %s", absPath)

	url := "/uploads/" + subDir + "/" + newName
	return url, nil
}
