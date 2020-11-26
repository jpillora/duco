package deploy

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func compile(target string, dest io.Writer) error {
	//check target
	if s, err := os.Stat(target); err != nil {
		return err
	} else if !s.IsDir() {
		return errors.New("target is not a dir")
	}
	//see if target is a main packge
	l, err := list(target)
	if err != nil {
		return err
	}
	log.Printf("golist: %#v", l.Module)
	//not a main package, write temp main
	if !l.Module.Main {
		tempTarget, err := tempMain(target)
		if err != nil {
			return err
		}
		// defer os.RemoveAll(tempTarget)
		target = tempTarget
	}
	//compile target directly into the buffer
	if err := goExec(
		dest,
		target,
		"build",
		"-trimpath",
		"-ldflags", "-v -s -w",
		"-v",
		"-o", "/dev/stdout",
		target,
	); err != nil {
		return errors.New("target build failed")
	}
	log.Printf("compiled: %s", target)
	return nil
}

func goExec(stdout io.Writer, wd string, args ...string) error {
	c := exec.Command("go", args...)
	c.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0")
	c.Stderr = os.Stderr
	c.Stdout = stdout
	c.Dir = wd
	if err := c.Run(); err != nil {
		return errors.New("target build failed")
	}
	return nil
}

func tempMain(target string) (string, error) {
	const maingo = `package main
	import (
		fn "%s"

		"github.com/jpillora/duco/runtime"
	)
	func main() {
		g := runtime.New()
		g.Add(fn.New())
		g.Start()
	}`
	tempDir := filepath.Join(target, "tmp", "bootstrap")
	if err := os.MkdirAll(tempDir, 0644); err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	tempMain := filepath.Join(tempDir, "main.go")
	if err := ioutil.WriteFile(tempMain, []byte(maingo), 0644); err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	log.Printf("tmp: %s", tempMain)
	return tempDir, nil
}
