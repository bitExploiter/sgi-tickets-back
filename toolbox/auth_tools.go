package toolbox

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword genera un hash bcrypt de una contraseña
// Cost factor: 12 (balance seguridad/performance)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash compara una contraseña con su hash
// Retorna true si coinciden, false si no
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateResetToken genera un token aleatorio de 32 bytes (64 chars hex)
// Retorna: token plano (para enviar por email), hash SHA-256 (para almacenar en DB)
func GenerateResetToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	token := hex.EncodeToString(bytes) // 64 caracteres

	// Hash del token para almacenar en DB
	hash := sha256.Sum256([]byte(token))
	hashString := hex.EncodeToString(hash[:])

	return token, hashString, nil
}

// ValidateResetToken valida que el token sea correcto y no haya expirado
// Verifica que el hash coincida y que no haya pasado más de 1 hora desde expiry
func ValidateResetToken(tokenPlain, tokenHash string, expiry time.Time) bool {
	// Verificar expiración (expiry debe ser mayor a NOW)
	if time.Now().After(expiry) {
		return false
	}

	// Verificar que el hash coincida
	hash := sha256.Sum256([]byte(tokenPlain))
	hashString := hex.EncodeToString(hash[:])

	return hashString == tokenHash
}

// GenerateSessionToken genera un token UUID v4 para cookies
func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Formato UUID v4
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return hex.EncodeToString(bytes), nil
}
