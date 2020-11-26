package deploy

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompileZip(target string, plugin bool) ([]byte, error) {
	//prepare buffer holding a zip "file"
	b := bytes.Buffer{}
	zw := zip.NewWriter(&b)
	fh := &zip.FileHeader{
		Name:   filepath.Base(target),
		Method: zip.Deflate,
	}
	fh.SetMode(os.ModePerm) //777
	fz, err := zw.CreateHeader(fh)
	if err != nil {
		return nil, err
	}
	//check target
	if s, err := os.Stat(target); err != nil {
		return nil, err
	} else if !s.IsDir() {
		return nil, errors.New("target is not a dir")
	}

	mode := "default"
	if plugin {
		mode = "plugin"
	}
	//compile target directly into the buffer
	args := []string{"build",
		"-buildmode=" + mode,
		"-trimpath",
		"-ldflags", "-v -s -w",
		"-v",
		"-o", "/dev/stdout",
		target,
	}
	log.Printf("go %s", strings.Join(args, " "))
	c := exec.Command("go", args...)
	c.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0")
	c.Stderr = os.Stderr
	c.Stdout = fz
	c.Dir = target
	if err := c.Run(); err != nil {
		return nil, errors.New("target build failed")
	}
	log.Printf("compiled: %s", fh.Name)
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
