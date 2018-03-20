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
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const (
	ROUTE_PREFIX = "/qrcode/"
)

const (
	SIZE100 = "100x100"
	SIZE200 = "200x200"
	SIZE300 = "300x300"
	SIZE400 = "400x400"
	SIZE500 = "500x500"
)

type qrCache struct {
	root      http.Dir
	maxLength int
}

func NewQRCache(root http.Dir, maxLength int) *qrCache {
	return &qrCache{
		root,
		maxLength,
	}
}

func (qrc *qrCache) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)
	fmt.Println("Req.URL:", req.URL)

	if _, err := os.Stat(string(qrc.root)); os.IsNotExist(err) {
		fmt.Println("Root folder is not existed!")
		panic(err)
	}

	if req.Method != "GET" {
		fmt.Println("Unsupported method type!")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	p := path.Clean(req.URL.Path)
	p = strings.TrimPrefix(p, ROUTE_PREFIX)

	if strings.Count(p, "/") > 1 {
		fmt.Println("Unsupported folder structure!")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	folder := path.Base(path.Dir(p))
	fmt.Println("Folder:", folder)

	var width, height int

	switch folder {
	case SIZE100:
		width = 100
		height = 100
	case SIZE200:
		width = 200
		height = 200
	case SIZE300:
		width = 300
		height = 300
	case SIZE400:
		width = 400
		height = 400
	case SIZE500:
		width = 500
		height = 500
	case ".":
		width = 200
		height = 200
	default:
		fmt.Println("Requested image size is not supported!")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	fileName := path.Base(p)
	v := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fmt.Println("Value:", v)

	if v == "" || len(v) > qrc.maxLength {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	dstPath := path.Join(folder, fileName)
	fmt.Println("File Path:", dstPath)

	fInfo, err := qrc.root.Open(dstPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error occurred when trying to open file:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fmt.Println("File is not existed. Try to generate new qrcode image file!")

		decoded, err := base64.URLEncoding.DecodeString(v)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		qrCode, err := qr.Encode(string(decoded), qr.M, qr.Auto)
		if err != nil {
			fmt.Println("Error occurred when trying to encode qrcode:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		qrCode, err = barcode.Scale(qrCode, width, height)
		if err != nil {
			fmt.Println("Error occurred when trying to scale qrcode:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		dstPath = path.Join(string(qrc.root), dstPath)
		absPath, _ := filepath.Abs(dstPath)

		folder = path.Join(string(qrc.root), folder)
		absFolderPath, _ := filepath.Abs(folder)
		fmt.Println("absPath: ", absPath)
		fmt.Println("absFolderPath: ", absFolderPath)

		err = os.MkdirAll(absFolderPath, 0755)
		if err != nil {
			fmt.Println("Error occurred when trying to create folder:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		file, err := os.Create(absPath)
		if err != nil {
			fmt.Println("Error occurred when trying to create file:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer file.Close()

		mw := io.MultiWriter(w, file)
		jpeg.Encode(mw, qrCode, nil)
	} else {
		defer fInfo.Close()
		fmt.Println("Try to read existed qrcode image file!")

		existingImg, err := ioutil.ReadAll(fInfo)
		if err != nil {
			fmt.Println("Error occurred when trying to read file:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Write(existingImg)
	}
}
