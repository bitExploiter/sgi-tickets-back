package toolbox

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret genera un nuevo secreto TOTP para un usuario
// Retorna el secreto codificado en base32 para almacenar en DB
func GenerateTOTPSecret(email string, issuer string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,       // "SGI Tickets"
		AccountName: email,
		Period:      30,            // Código válido por 30 segundos
		SecretSize:  32,            // 256 bits de entropía
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})

	if err != nil {
		return "", err
	}

	return key.Secret(), nil
}

// GenerateQRCodeBase64 genera la imagen QR en formato base64 para escanear
// Input: URL del TOTP (otpauth://...)
func GenerateQRCodeBase64(totpURL string) (string, error) {
	// Parsear el URL TOTP
	key, err := otp.NewKeyFromURL(totpURL)
	if err != nil {
		return "", err
	}

	// Generar imagen QR (200x200 pixels)
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}

	// Convertir a PNG en memoria
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	// Codificar en base64 para devolver al frontend
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return "data:image/png;base64," + imgBase64, nil
}

// GetTOTPURL genera el URL otpauth:// para generar QR
func GetTOTPURL(secret, email, issuer string) string {
	// Reconstruir el TOTP key desde el secreto guardado
	return "otpauth://totp/" + issuer + ":" + email + "?secret=" + secret + "&issuer=" + issuer
}

// ValidateTOTPCode valida un código de 6 dígitos contra el secreto del usuario
// Permite ±1 periodo de tiempo (90 segundos de ventana total)
func ValidateTOTPCode(code, secret string) bool {
	return totp.Validate(code, secret)
}

// Require2FA determina si un rol requiere 2FA obligatorio
// Roles CON 2FA: admin, supervisor, agente, contratista
// Roles SIN 2FA: entidad
func Require2FA(rol string) bool {
	rolesWithout2FA := []string{"entidad"}

	for _, r := range rolesWithout2FA {
		if r == rol {
			return false
		}
	}

	return true
}
