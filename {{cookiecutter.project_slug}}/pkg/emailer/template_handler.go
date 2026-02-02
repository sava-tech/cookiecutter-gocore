package emailer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

// GenerateHTML loads an HTML file and injects the given data into it.
// Can be reused for any HTML email or webpage.
func GenerateHTML(templatePath string, data interface{}) (string, error) {
	// Ensure template exists
	absPath, err := filepath.Abs(templatePath)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", err
	}

	tmpl, err := template.ParseFiles(absPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}