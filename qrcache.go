package qrcachego

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func GenerateQRImage(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)
	fmt.Println("Req.URL:", req.URL)

	if req.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	p := path.Clean(req.URL.Path)
	folder := path.Base(path.Dir(p))
	fmt.Println("Folder:", folder)
	splits := strings.Split(folder, "x")

	width, _ := strconv.Atoi(splits[0])
	height, _ := strconv.Atoi(splits[1])

	if width == 0 {
		width = 200
	}

	if height == 0 {
		height = 200
	}

	fileName := path.Base(p)
	v := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fmt.Println("Value:", v)

	if v == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	absFolderPath, _ := filepath.Abs(folder)

	fPath := path.Join(folder, fileName)
	absPath, _ := filepath.Abs(fPath)
	fmt.Println("File Path:", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Println("Try to generate new qrcode image file!")

		decoded, err := base64.URLEncoding.DecodeString(v)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		qrCode, err := qr.Encode(string(decoded), qr.M, qr.Auto)
		if err != nil {
			fmt.Println("Error occurred when trying to encode qrcode:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		qrCode, err = barcode.Scale(qrCode, width, height)
		if err != nil {
			fmt.Println("Error occurred when trying to scale qrcode:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll(absFolderPath, 0755)
		if err != nil {
			fmt.Println("Error occurred when trying to create folder:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		file, err := os.Create(absPath)
		if err != nil {
			fmt.Println("Error occurred when trying to create file:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		mw := io.MultiWriter(w, file)
		jpeg.Encode(mw, qrCode, nil)

	} else {
		fmt.Println("Try to read existed qrcode image file!")

		existingImg, err := ioutil.ReadFile(absPath)
		if err != nil {
			fmt.Println("Error occurred when trying to read file:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(existingImg)
	}
}
