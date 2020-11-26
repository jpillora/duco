package deploy

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
)

func compileZip(target string) ([]byte, error) {
	//prepare buffer holding a zip "file"
	b := bytes.Buffer{}
	zw := zip.NewWriter(&b)
	fh := &zip.FileHeader{
		Name:   "bootstrap",
		Method: zip.Deflate,
	}
	fh.SetMode(os.ModePerm) //777
	fz, err := zw.CreateHeader(fh)
	if err != nil {
		return nil, err
	}
	if err := compile(target, fz); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	log.Printf("zip file: %d", b.Len())
	return b.Bytes(), nil
}

func hash(z []byte) string {
	h := sha256.New()
	h.Write(z)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
