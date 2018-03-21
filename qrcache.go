package qrcachego

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

//Supported Size
const (
	SIZE100 = "100x100"
	SIZE200 = "200x200"
	SIZE300 = "300x300"
	SIZE400 = "400x400"
	SIZE500 = "500x500"
	DEFAULT = "."
)

type qrCache struct {
	root      http.Dir
	maxLength int
	folderMap map[string]int
	gcCh      chan string
}

func NewQRCache(root http.Dir, maxLength int, expiry time.Duration) *qrCache {
	if _, err := os.Stat(string(root)); os.IsNotExist(err) {
		fmt.Println("Root folder is not existed!")
		err = os.Mkdir(string(root), 0755)
		if err != nil {
			fmt.Println("Error occurred when trying to create root folder:", err)
			panic(err)
		}
	}

	folderMap := make(map[string]int)
	folderMap[SIZE100] = 100
	folderMap[SIZE200] = 200
	folderMap[SIZE300] = 300
	folderMap[SIZE400] = 400
	folderMap[SIZE500] = 500
	folderMap[DEFAULT] = 200

	for d := range folderMap {
		sd := path.Join(string(root), d)
		if _, err := os.Stat(sd); os.IsNotExist(err) {
			fmt.Printf("Folder %s is not existed!\n", d)
			err = os.Mkdir(string(sd), 0755)
			if err != nil {
				fmt.Println("Error occurred when trying to create default folder:", err)
				panic(err)
			}
		}
	}

	var gcCh chan string

	if expiry > 0 {
		gcCh = make(chan string)
		go func() {
			gcMap := make(map[string]time.Time)
			for {
				f := <-gcCh

				gcMap[f] = time.Now()
				for k, v := range gcMap {
					if time.Now().Sub(v) > expiry {
						os.Remove(k)
						fmt.Println("Remove file:", k)
						delete(gcMap, k)
					}
				}
			}
		}()
	}

	return &qrCache{
		root,
		maxLength,
		folderMap,
		gcCh,
	}
}

func (qrc *qrCache) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)
	fmt.Println("Req.URL:", req.URL)

	if req.Method != "GET" {
		fmt.Println("Unsupported method type!")
		http.NotFound(w, req)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	p := path.Clean(req.URL.Path)

	folder := path.Base(path.Dir(p))
	fmt.Println("Folder:", folder)

	if _, ok := qrc.folderMap[folder]; !ok {
		fmt.Println("Requested image size is not supported!")
		http.NotFound(w, req)
		return
	}

	width := qrc.folderMap[folder]
	height := width

	fileName := path.Base(p)
	v := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fmt.Println("Value:", v)

	if v == "" {
		fmt.Println("Empty value.")
		http.NotFound(w, req)
		return
	}

	if len(v) > qrc.maxLength {
		fmt.Println("Exceeded maximum value length.")
		http.NotFound(w, req)
		return
	}

	if filepath.Ext(fileName) != ".jpg" {
		fmt.Println("Unsupported file format.")
		http.NotFound(w, req)
		return
	}

	fmt.Println("File Path:", p)

	dstPath := path.Join(string(qrc.root), p)
	absPath, _ := filepath.Abs(dstPath)

	fmt.Println("absPath: ", absPath)

	fInfo, err := qrc.root.Open(p)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error occurred when trying to open file:", err)
			http.NotFound(w, req)
			return
		}

		fmt.Println("File is not existed. Try to generate new qrcode image file!")

		decoded, err := base64.URLEncoding.DecodeString(v)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		qrCode, err := qr.Encode(string(decoded), qr.M, qr.Auto)
		if err != nil {
			fmt.Println("Error occurred when trying to encode qrcode:", err)
			http.NotFound(w, req)
			return
		}

		// bounds := qrCode.Bounds()
		// orgLength := bounds.Max.X - bounds.Min.X
		// fmt.Println("orgLength:", orgLength)

		qrCode, err = barcode.Scale(qrCode, width, height)
		if err != nil {
			fmt.Println("Error occurred when trying to scale qrcode:", err)
			http.NotFound(w, req)
			return
		}

		file, err := os.Create(absPath)
		if err != nil {
			fmt.Println("Error occurred when trying to create file:", err)
			http.NotFound(w, req)
			return
		}
		defer file.Close()

		jpeg.Encode(io.MultiWriter(w, file), qrCode, nil)
	} else {
		defer fInfo.Close()
		fmt.Println("Try to read existed qrcode image file!")

		_, err = io.Copy(w, fInfo)
		if err != nil {
			fmt.Println("Error occurred when trying to return content:", err)
			return
		}
	}

	if qrc.gcCh != nil {
		qrc.gcCh <- absPath
	}
}
