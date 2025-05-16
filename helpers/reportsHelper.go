package helpers

import (
	"bytes"
	"encoding/base64"

	"github.com/tealeg/xlsx/v3"
)

func EncodeExcelToBase64(excelFile *xlsx.File) (string, error) {
	// Create a buffer to hold the Excel file's data in memory
	var buffer bytes.Buffer

	// Write the Excel file to the buffer
	err := excelFile.Write(&buffer)
	if err != nil {
		return "", err
	}

	// Convert the buffer's content to a byte slice
	data := buffer.Bytes()

	// Encode the byte slice to base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}
