package main

import (
	"bytes"
	"fmt"
	"github.com/apaladiychuk/qrgen/serverapi"
	"github.com/skip2/go-qrcode"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type QrData struct {
	BrandName string
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
		fmt.Printf(`usage  %s "параметры qr"  "полное имя файла" `, os.Args[0])
		fmt.Println(`----------------`)
		fmt.Println(`Формат данных для Qr `)
		fmt.Println(`название бренда#уникальный код#Модель#Верх#Подошва`)
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
	var qrParam QrData
	if len(paramArr) < 5 {
		qrParam = QrData{
			ModelId:   paramArr[0],
			ModelName: paramArr[1],
			Top:       paramArr[2],
			Sole:      paramArr[3],
			BrandName: "Reckless",
		}
	} else {
		qrParam = QrData{
			BrandName: paramArr[0],
			ModelId:   paramArr[1],
			ModelName: paramArr[2],
			Top:       paramArr[3],
			Sole:      paramArr[4],
		}
	}

	GenerateQr(qrParam, path, qrType)

}

func GenerateQr(info QrData, path string, qrType string) {

	var url string
	if info.BrandName == "Reckless" {
		url = "http://www.reckless.me/"
	} else {
		url = "http://esente.com.ua/"
	}
	keyString := fmt.Sprintf(
		`%s 
Модель: %s 
Верх: %s
Подошва: %s
%s
#%s#
`,
		info.BrandName, info.ModelName, info.Top, info.Sole, url, info.ModelId)

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
	serverapi.UploadInventory(info.ModelId, info.ModelName)

	os.Exit(0)
}
