package qrcachego_test

import (
	"bytes"
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

	TestRoot      = "test"
	TestMaxLength = 300
	TestExpiry    = 1
)

var (
	TestAnswer = testAnswer()
)

type cacheTest struct {
	def    string // test case definition
	method string // http method
	url    string // http URL
	status int    // the expected status code
	want   []byte // the expected output
}

var cacheTests = []cacheTest{
	{
		def:    "Wrong request method",
		method: "POST",
		url:    "",
		status: http.StatusNotFound,
	},
	{
		def:    "Unsupported image size",
		method: "GET",
		url:    "http://localhost/123x123/aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==.jpg",
		status: http.StatusNotFound,
	},
	{
		def:    "Empty value",
		method: "GET",
		url:    "http://localhost/200x200/.jpg",
		status: http.StatusNotFound,
	},
	{
		def:    "Exceeded maximum value length.",
		method: "GET",
		url:    "http://localhost/200x200/aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==.jpg",
		status: http.StatusNotFound,
	},
	{
		def:    "Unsupported file format.",
		method: "GET",
		url:    "http://localhost/200x200/aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==.png",
		status: http.StatusNotFound,
	},
	{
		def:    "Success",
		method: "GET",
		url:    "http://localhost/200x200/aHR0cHM6Ly9yaW5rZWJ5LmV0aGVyc2Nhbi5pby9hZGRyZXNzLzB4ZjBkNjU0Nzk3MzJlZWRjNDA2YzAwZmZiMjliYzlkZDQyNjc4MGVlNA==.jpg",
		status: http.StatusOK,
		want:   TestAnswer,
	},
}

func testAnswer() []byte {
	qrCode, _ := qr.Encode(TestValue, qr.M, qr.Auto)
	width, _ := strconv.Atoi(TestWidth)
	height, _ := strconv.Atoi(TestHeight)
	qrCode, _ = barcode.Scale(qrCode, width, height)

	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, qrCode, nil)

	return buffer.Bytes()
}

func TestCache(t *testing.T) {
	for _, test := range cacheTests {
		t.Run(test.def, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Error("Request Creation Failed: ", err)
				return
			}

			reqr := httptest.NewRecorder()

			cache := qcg.NewQRCache(http.Dir(TestRoot), TestMaxLength, TestExpiry, false)

			cache.ServeHTTP(reqr, req)
			if status := reqr.Code; status != test.status {
				t.Errorf("Status code differs. Expected %d \n Got %d", test.status, status)
				return
			}

			if test.want == nil {
				return
			}

			resp := reqr.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if bytes.Compare(body, test.want) != 0 {
				t.Errorf("Wrong Result.")
			}
		})
	}
}
