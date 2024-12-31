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
    "encoding/base64"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	filePath    string
	content     []byte
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) OpenFile() (string, error) {
	// Показываем диалог выбора файла
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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

func (a *App) Save(content, password string) (string, error) {
	// Если указан пароль, шифруем содержимое
	var data []byte
	var err error

	if password != "" {
		encrypted, err := a.Encrypt(content, password)
		if err != nil {
			return "", err
		}
		data = []byte(encrypted)
	} else {
		data = []byte(content)
	}

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
	// Кодируем в base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt расшифровывает текст
func (a *App) Decrypt(content, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	// Декодируем из base64
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
	return string(data), nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
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
        Assets: assets,
        OnStartup: app.startup,
    })
	if err != nil {
		println("Error:", err.Error())
	}
}
