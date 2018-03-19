package main

import (
	"fmt"
	"net/http"

	qcg "github.com/blazingorb/qrcachego"
	"github.com/rs/cors"
)

func main() {
	access := cors.AllowAll().Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/qrcode/", qcg.GenerateQRImage)

	err := http.ListenAndServe(":8888", access(mux))
	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
	}
}
