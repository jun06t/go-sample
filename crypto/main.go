package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"

	"golang.org/x/crypto/argon2"
)

const (
	keySize  = 32
	saltSize = 16
)

func main() {
	pass := "hogehoge"
	data := []byte("secret text")

	cek, err := generateCEK()
	if err != nil {
		log.Fatal(err)
	}
	ciphertext, err := encrypt(cek, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ciphertext: %x\n", ciphertext)

	salt, ecek, err := encryptCEK(pass, cek)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("salt: %x\n", salt)
	fmt.Printf("encrypted CEK: %x\n", ecek)

	dcek, err := decryptCEK(pass, salt, ecek)
	if err != nil {
		log.Fatal(err)
	}
	ddata, err := decrypt(dcek, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("plaintext:", string(ddata))
}

func encryptCEK(password string, cek []byte) ([]byte, []byte, error) {
	kek, salt, err := generateKEK(password, nil)
	if err != nil {
		return nil, nil, err
	}

	encrypted, err := encrypt(kek, cek)
	if err != nil {
		return nil, nil, err
	}
	return salt, encrypted, nil
}

func decryptCEK(password string, salt, ecek []byte) ([]byte, error) {
	kek, salt, err := generateKEK(password, salt)
	if err != nil {
		return nil, err
	}

	cek, err := decrypt(kek, ecek)
	if err != nil {
		return nil, err
	}
	return cek, nil
}

func generateCEK() ([]byte, error) {
	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func generateKEK(password string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, saltSize)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}
	key := argon2.Key([]byte(password), salt, 3, 32*1024, 4, keySize)

	return key, salt, nil
}

func encrypt(key, plaintext []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, err
	}
	nonce, encrypted := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
