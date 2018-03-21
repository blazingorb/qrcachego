package main

import (
	"fmt"
	"net/http"
	"time"

	qcg "github.com/blazingorb/qrcachego"
)

const (
	ROOT_PATH  = "example"
	MAX_LENGTH = 300
	EXPIRY     = 1 * time.Minute
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/qrcode/", http.StripPrefix("/qrcode/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, EXPIRY)))
	mux.Handle("/qrcode-perm/", http.StripPrefix("/qrcode-perm/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, EXPIRY)))

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
	}
}
