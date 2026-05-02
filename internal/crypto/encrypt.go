package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
)

var secretKey []byte

func Init(key string) error {
    if len(key) == 0 {
        return errors.New("encryption key is empty")
    }
    hash := sha256.Sum256([]byte(key))
    secretKey = hash[:]
    return nil
}

func Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }
    aead, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, aead.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    ciphertext := aead.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string) (string, error) {
    ciphertext, err := base64.URLEncoding.DecodeString(encoded)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }
    aead, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    if len(ciphertext) < aead.NonceSize() {
        return "", errors.New("ciphertext too short")
    }
    nonce, ciphertext := ciphertext[:aead.NonceSize()], ciphertext[aead.NonceSize():]
    plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}