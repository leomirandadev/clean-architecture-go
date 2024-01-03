package aescrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
)

type AESCryptDoer interface {
	Encrypt(toEncrypt []byte) ([]byte, error)
	Decrypt(encryptedData []byte) ([]byte, error)
	EncryptString(toEncrypt string) (string, error)
	DecryptString(encryptedData string) (string, error)
}

func New(key string) AESCryptDoer {
	if len(key) != 64 {
		slog.Debug("Hexadecimal AES256 key is not set correctly", "keylen", len(key))
		panic("hexadecimal AES256 key is not set correctly")
	}

	keyByte, err := hex.DecodeString(key)
	if err != nil {
		slog.Debug("Failed to decode hexadecimal key", "error", err)
		panic(err)
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		slog.Debug("NewCipher fails", "err", err)
		panic(err)
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		slog.Debug("NewGCM fails", "err", err)
		panic(err)
	}

	return impl{aesGCM: aesGCM, nonceSize: aesGCM.NonceSize()}
}

type impl struct {
	aesGCM    cipher.AEAD
	nonceSize int
}

func (i impl) Encrypt(toEncrypt []byte) ([]byte, error) {
	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, i.nonceSize)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		slog.Debug("ReadFull fails", "err", err)
		return nil, err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	return i.aesGCM.Seal(nonce, nonce, toEncrypt, nil), nil
}

func (i impl) Decrypt(encryptedData []byte) ([]byte, error) {
	//Extract the nonce from the encrypted data
	nonce, ciphertext := encryptedData[:i.nonceSize], encryptedData[i.nonceSize:]

	//Decrypt the data
	result, err := i.aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		slog.Debug("aesGCM.Open fails", "err", err)
		return nil, err
	}

	return result, nil
}

func (i impl) EncryptString(toEncrypt string) (string, error) {
	result, err := i.Encrypt([]byte(toEncrypt))
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (i impl) DecryptString(encryptedData string) (string, error) {
	result, err := i.Decrypt([]byte(encryptedData))
	if err != nil {
		return "", err
	}

	return string(result), nil
}
