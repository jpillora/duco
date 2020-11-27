package deploy

import (
	"errors"
	"fmt"
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
	if !l.IsMain() {
		tempTarget, err := tempMain(target, l.ImportPath)
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
		"-ldflags", "-s -w",
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

func tempMain(target, importPath string) (string, error) {
	s, err := os.Stat(target)
	if err != nil {
		panic("should not get here")
	}
	const mainTemplate = `package main
	import (
		fn "%s"

		"github.com/jpillora/duco/runtime"
	)
	func main() {
		g := runtime.New()
		g.Add(fn.New())
		g.Start()
	}`
	mainGoFile := fmt.Sprintf(mainTemplate, importPath)
	tempDir := filepath.Join(target, "tmp", "bootstrap")
	if err := os.MkdirAll(tempDir, s.Mode().Perm()); err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	tempMain := filepath.Join(tempDir, "main.go")
	if err := ioutil.WriteFile(tempMain, []byte(mainGoFile), s.Mode().Perm()); err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	log.Printf("tmp: %s", tempMain)
	return tempDir, nil
}
