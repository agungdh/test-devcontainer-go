// Package docs menyimpan spesifikasi OpenAPI (embed ke binary
// supaya tidak tergantung working directory saat dijalankan).
package docs

import (
	_ "embed"
)

//go:embed openapi.yaml
var Spec []byte
