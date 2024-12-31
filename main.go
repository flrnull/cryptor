package main

import (
        "embed"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
        "context"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	filePath    string
	isEncrypted bool
	content     []byte
}

func NewApp() *App {
	return &App{}
}

func (a *App) OpenFile(ctx context.Context) (string, error) {
	// Показываем диалог выбора файла
	filePath, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select a file to open",
	})
	if err != nil {
		println("Error while opening file dialog:", err.Error())
		return "", err
	}

	if filePath == "" {
		println("No file selected.")
		return "", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		println("Error while reading file:", err.Error())
		return "", err
	}

	a.filePath = filePath
	return string(data), nil
}

// Encrypt шифрует текст
func (a *App) Encrypt(content, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(content))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(content))
	return string(ciphertext), nil
}

// Decrypt расшифровывает текст
func (a *App) Decrypt(content, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	data := []byte(content)
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return string(data), nil
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
            Assets: assets,
    OnStartup: func(ctx context.Context) {
    },
        })
	if err != nil {
		println("Error:", err.Error())
	}
}
