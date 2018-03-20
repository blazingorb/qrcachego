package qrcachego_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	qcg "github.com/blazingorb/qrcachego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const (
	TestValue  = "https://rinkeby.etherscan.io/address/0xf0d65479732eedc406c00ffb29bc9dd426780ee4"
	TestWidth  = "200"
	TestHeight = "200"
	TestURL1   = "http://localhost:8888/qrcode/200x200/aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==.jpg"

	TestRoute     = "/qrcode/"
	TestRoot      = "."
	TestMaxLength = 300
)

var (
	TestEncoded = base64.URLEncoding.EncodeToString([]byte(TestValue))
	TestAnswer  = testAnswer()
)

func testAnswer() []byte {
	qrCode, _ := qr.Encode(TestValue, qr.M, qr.Auto)
	width, _ := strconv.Atoi(TestWidth)
	height, _ := strconv.Atoi(TestHeight)
	qrCode, _ = barcode.Scale(qrCode, width, height)

	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, qrCode, nil)

	return buffer.Bytes()
}

func TestCacheSuccess(t *testing.T) {

	req, err := http.NewRequest("GET", TestURL1, nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}

	fmt.Println(req.URL.String())

	reqr := httptest.NewRecorder()

	cache := qcg.NewQRCache(http.Dir(TestRoot), TestMaxLength)

	cache.ServeHTTP(reqr, req)
	if status := reqr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusOK, status)
	}

	resp := reqr.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if bytes.Compare(body, TestAnswer) != 0 {
		t.Errorf("Wrong Result.")
	}

}
