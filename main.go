package main

import (
	"bytes"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type QrData struct {
	ModelId   string
	ModelName string
	Top       string
	Sole      string
}

func main() {

	var paramString string
	var path string
	var qrType string
	if len(os.Args) < 3 {
		fmt.Println("Wrong argument ")
		fmt.Printf(`usage  %s "параметры qr"  "полное имя файла"\n `, os.Args[0])
		fmt.Println(`уникальный код#Модель#Верх#Подошва`)
		return
	}
	paramString = os.Args[1]
	path = os.Args[2]
	if len(os.Args) > 3 {
		qrType = os.Args[3]
		if qrType != "jpeg" && qrType != "png" {
			qrType = "jpeg"
		}
	} else {
		qrType = "jpeg"
	}

	paramArr := strings.Split(paramString, "#")
	if len(paramArr) < 4 {
		fmt.Errorf("неправильный формат Qr ")
		return
	}

	qrParam := QrData{
		ModelId:   paramArr[0],
		ModelName: paramArr[1],
		Top:       paramArr[2],
		Sole:      paramArr[3],
	}

	GenerateQr(qrParam, path, qrType)

}

func GenerateQr(info QrData, path string, qrType string) {

	keyString := fmt.Sprintf(
		`Reckless 
Модель: %s 
Верх: %s
Подошва: %s
http://www.reckless.me/
#%s#"
`,
		info.ModelName, info.Top, info.Sole, info.ModelId)

	if pngBuff, err := qrcode.Encode(keyString, qrcode.Medium, 256); err != nil {
		fmt.Errorf("[model] generate ", err.Error())

	} else {
		if jpg, err := png.Decode(bytes.NewReader(pngBuff)); err != nil {
			fmt.Errorf("tojpeg ", err.Error())
		} else {

			if outFile, err := os.Create(path); err != nil {
				fmt.Errorf("jpeg to butes  ", err.Error())
			} else {
				defer outFile.Close()
				if err := jpeg.Encode(outFile, jpg, &jpeg.Options{Quality: 100}); err != nil {
					fmt.Errorf("jpeg to butes  ", err.Error())
				}
			}
		}
	}

}
