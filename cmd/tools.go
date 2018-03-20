package main

import (
	"fmt"
	"net/http"

	qcg "github.com/blazingorb/qrcachego"
	"github.com/rs/cors"
)

const (
	ROOT_PATH      = "test"
	MAX_LENGTH     = 300
	MAX_QUEUE_SIZE = 250
)

func main() {
	access := cors.AllowAll().Handler
	mux := http.NewServeMux()
	mux.Handle("/qrcode/", http.StripPrefix("/qrcode/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, MAX_QUEUE_SIZE)))

	err := http.ListenAndServe(":8888", access(mux))
	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
	}
}
