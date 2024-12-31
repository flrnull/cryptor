package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type App struct {
	filePath    string
	isEncrypted bool
	content     []byte
}

func NewApp() *App {
	return &App{}
}

func (a *App) OpenFile(filePath string) (string, error) {
	a.filePath = filePath
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	a.content = data
	if isEncrypted(filePath) {
		a.isEncrypted = true
		return "Encrypted", nil
	}
	a.isEncrypted = false
	return string(data), nil
}

func (a *App) SaveFile(content string, password string) error {
	data := []byte(content)
	if a.isEncrypted {
		encryptedData, err := encrypt(data, password)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(a.filePath, encryptedData, 0644)
	}
	return ioutil.WriteFile(a.filePath, data, 0644)
}

func (a *App) Decrypt(password string) (string, error) {
	if !a.isEncrypted {
		return string(a.content), nil
	}
	decryptedData, err := decrypt(a.content, password)
	if err != nil {
		return "", err
	}
	a.isEncrypted = false
	return string(decryptedData), nil
}

func isEncrypted(filePath string) bool {
	return len(filePath) > 4 && filePath[len(filePath)-4:] == ".enc"
}

func encrypt(data []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

func decrypt(data []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
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
        })
	if err != nil {
		println("Error:", err.Error())
	}
}
