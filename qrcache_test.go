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
	TestValue  = "TestValue"
	TestWidth  = "200"
	TestHeight = "200"
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

func TestGenerateQRImage(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}
	req.Header.Set("Content-Type", "image/jpeg")
	q := req.URL.Query()
	q.Add("value", TestEncoded)
	q.Add("width", TestWidth)
	q.Add("height", TestHeight)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	reqr := httptest.NewRecorder()

	http.HandlerFunc(qcg.GenerateQRImage).ServeHTTP(reqr, req)
	if status := reqr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusOK, status)
	}

	resp := reqr.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if bytes.Compare(body, TestAnswer) != 0 {
		t.Errorf("Wrong Result.")
	}

}
