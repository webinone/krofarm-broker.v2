package libs

import (
	"fmt"
	"strings"
	"encoding/json"
	"bytes"
)

// Pointer Argument !!!
func CleanJsonBody ( body *[]byte) {

	stringBody := string(*body)
	stringBody = strings.Replace(stringBody, "\n", " ", -1)
	stringBody = strings.Replace(stringBody, "\r", " ", -1)
	stringBody = strings.Replace(stringBody, "\t", " ", -1)
	stringBody = strings.Replace(stringBody, "\\\\", " ", -1)

	*body = []byte(stringBody)

	dst := new(bytes.Buffer)
	err := json.Compact(dst, *body)

	if err != nil {
		fmt.Println(err)
	}

	*body = []byte(dst.Bytes())
}
