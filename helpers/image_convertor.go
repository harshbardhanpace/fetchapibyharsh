package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
)

//Given a base64 string of a PNG, encodes it into an PNG image test.png and upload it to s3
func Base64toPng(name string, id string, data string) (error, string) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println("base64toPng", bounds, formatString)

	//Encode from image format to writer
	// pngFilename := "test.png"
	pngFilename := name + id + ".png"
	file, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Error in open file %v \n", err)
		return err, ""
	}

	err = png.Encode(file, m)
	if err != nil {
		fmt.Printf("Error in png encode %v \n", err)
		return err, ""
	}
	file.Close()

	upFile, err := os.Open(pngFilename)
	if err != nil {
		fmt.Printf("Error opening file %v \n", err)
		return err, ""
	}

	// Get the file info
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	upFile.Close()
	err = os.Remove(pngFilename)
	if err != nil {
		fmt.Printf("Error in deleting image %v \n", err)
		return err, ""
	}

	var s3ImageLocation string
	err, s3ImageLocation = UploadToS3(pngFilename, bytes.NewReader(fileBuffer))
	if err != nil {
		fmt.Printf("Error in upload to s3 %v \n", err)
		return err, ""
	}

	return nil, s3ImageLocation

}
