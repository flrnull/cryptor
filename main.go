package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	filePath string
	content  []byte
	ctx      context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) OpenFile() (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a file to open",
	})
	if err != nil {
		return "", err
	}

	if filePath == "" {
		return "", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	a.filePath = filePath
	return string(data), nil
}

func (a *App) Save(content string) (string, error) {
	var data = []byte(content)
	var err error

	if a.filePath != "" {
		err = os.WriteFile(a.filePath, data, 0644)
		if err != nil {
			return "", err
		}
		return "File saved successfully.", nil
	}

	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: "untitled.txt",
	})
	if err != nil {
		return "", err
	}

	if savePath == "" {
		return "", errors.New("save operation cancelled")
	}

	err = os.WriteFile(savePath, data, 0644)
	if err != nil {
		return "", err
	}

	a.filePath = savePath
	return "File saved successfully.", nil
}

func (a *App) Encrypt(content, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Генерируем хеш текста
	hash := sha256.Sum256([]byte(content))
	contentWithHash := append([]byte(content), hash[:]...)

	ciphertext := make([]byte, aes.BlockSize+len(contentWithHash))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], contentWithHash)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *App) Decrypt(content, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	if len(data) < sha256.Size {
		return "", errors.New("invalid data length")
	}

	// Извлекаем хеш и проверяем его
	decryptedContent := data[:len(data)-sha256.Size]
	expectedHash := data[len(data)-sha256.Size:]
	actualHash := sha256.Sum256(decryptedContent)

	if !equalHashes(expectedHash, actualHash[:]) {
		return "", errors.New("decryption failed: invalid password or corrupted data")
	}

	return string(decryptedContent), nil
}

// Проверка на равенство хешей
func equalHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "Cryptor",
		Width:  800,
		Height: 600,
		Bind: []interface{}{
			app,
		},
		Assets:    assets,
		OnStartup: app.startup,
	})
	if err != nil {
		println("Error:", err.Error())
	}
}