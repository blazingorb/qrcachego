package main

import (
	"fmt"
	"net/http"

	qcg "github.com/blazingorb/qrcachego"
	"github.com/rs/cors"
)

const (
	ROOT_PATH     = "test"
	MAX_LENGTH    = 300
	EXPIRY_IN_MIN = 1
)

func main() {
	access := cors.AllowAll().Handler
	mux := http.NewServeMux()
	mux.Handle("/qrcode/", http.StripPrefix("/qrcode/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, EXPIRY_IN_MIN, true)))
	mux.Handle("/qrcode-perm/", http.StripPrefix("/qrcode-perm/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, EXPIRY_IN_MIN, false)))

	err := http.ListenAndServe(":8888", access(mux))
	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
	}
}
