package toolbox

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string][]string{
	"image":    {".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"},
	"video":    {".mp4", ".mov", ".avi", ".mkv", ".webm", ".flv", ".wmv"},
	"document": {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"},
}

var allowedMIMETypes = map[string][]string{
	"image":    {"image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp"},
	"video":    {"video/mp4", "video/quicktime", "video/x-msvideo", "video/x-matroska", "video/webm", "video/x-flv", "video/x-ms-wmv"},
	"document": {"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "text/plain"},
}

// Firmas de archivo (magic numbers)
var fileSignatures = map[string][][]byte{
	".jpg":  {{0xFF, 0xD8, 0xFF}},
	".jpeg": {{0xFF, 0xD8, 0xFF}},
	".png":  {{0x89, 0x50, 0x4E, 0x47}},
	".gif":  {{0x47, 0x49, 0x46, 0x38}},
	".pdf":  {{0x25, 0x50, 0x44, 0x46}},
	".mp4":  {{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}, {0x00, 0x00, 0x00, 0x1C, 0x66, 0x74, 0x79, 0x70}, {0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70}},
}

func ValidateUploadedFile(filename string, fileContent []byte, category string) error {
	if err := ValidateFileExtension(filename, category); err != nil {
		return err
	}
	if err := ValidateFileMIME(filename, fileContent); err != nil {
		return err
	}
	return ValidateFileSignature(filename, fileContent)
}

func ValidateFileExtension(filename string, category string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	extensions, ok := allowedExtensions[category]
	if !ok {
		return fmt.Errorf("categoria de archivo no valida: %s", category)
	}
	for _, allowed := range extensions {
		if ext == allowed {
			return nil
		}
	}
	return fmt.Errorf("extension de archivo no permitida: %s", ext)
}

func ValidateFileMIME(filename string, fileContent []byte) error {
	detectedMIME := http.DetectContentType(fileContent)
	ext := strings.ToLower(filepath.Ext(filename))

	for category, extensions := range allowedExtensions {
		for _, allowedExt := range extensions {
			if ext == allowedExt {
				mimes := allowedMIMETypes[category]
				for _, mime := range mimes {
					if strings.HasPrefix(detectedMIME, mime) {
						return nil
					}
				}
				return fmt.Errorf("tipo MIME no coincide con la extension: detectado %s para %s", detectedMIME, ext)
			}
		}
	}
	return fmt.Errorf("extension no reconocida: %s", ext)
}

func ValidateFileSignature(filename string, fileContent []byte) error {
	ext := strings.ToLower(filepath.Ext(filename))
	signatures, ok := fileSignatures[ext]
	if !ok {
		// Si no hay firma registrada, se permite
		return nil
	}
	for _, sig := range signatures {
		if len(fileContent) >= len(sig) {
			match := true
			for i, b := range sig {
				if fileContent[i] != b {
					match = false
					break
				}
			}
			if match {
				return nil
			}
		}
	}
	return fmt.Errorf("la firma del archivo no coincide con la extension: %s", ext)
}
