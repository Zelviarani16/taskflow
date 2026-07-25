package helpers

import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah password teks biasa menjadi hash bcrypt.
// Cost 10 adalah kompromi standar antara keamanan dan kecepatan.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword membandingkan password teks biasa dengan hash yang tersimpan.
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
