package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Script para generar hash bcrypt de una contraseña
// Uso: go run scripts/hash_password.go "tu_contraseña"
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run scripts/hash_password.go \"tu_contraseña\"")
		os.Exit(1)
	}

	password := os.Args[1]

	// Generar hash con cost 12 (igual que auth_tools.go)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Println("Error generando hash:", err)
		os.Exit(1)
	}

	fmt.Println("Contraseña:", password)
	fmt.Println("Hash bcrypt:", string(hash))
	fmt.Println("\nSQL de ejemplo:")
	fmt.Printf("INSERT INTO tickets_usuarios (nombres, apellidos, email, password, rol, activo, created_at, updated_at)\n")
	fmt.Printf("VALUES ('Test', 'Usuario', 'test@sgi.com', '%s', 'admin', true, NOW(), NOW());\n", string(hash))
}
