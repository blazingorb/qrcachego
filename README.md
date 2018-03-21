QRCache for GOlang 
==================================


## Requirements
- Go 1.5 or later.
- [barcode]

[barcode]: https://github.com/boombuler/barcode

## Usage

Sample Usage for using QRCache:
- /qrcode/: Cache qrcode image with GC
- /qr/

```go
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
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/qrcode/", http.StripPrefix("/qrcode/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, 1*time.Minute)))
	mux.Handle("/qrcode-perm/", http.StripPrefix("/qrcode-perm/", qcg.NewQRCache(http.Dir(ROOT_PATH), MAX_LENGTH, -1)))

	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
	}
}

```

For more detail sample cases, please refer to files under example folder.

## Notes
    No GC when expiry is negative