package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

type cli struct {
	Name string `opts:"help=function name"`
	Role string `opts:"help=role name"`
	App  string `opts:"mode=arg, help=target app"`
}

func main() {
	//cli
	c := cli{
		Name: "go-raw-runtime",
		Role: "arn:aws:iam::652507618334:role/lambda-role",
	}
	opts.New(&c).Name("gambda").Parse()
	//zip code
	z, err := zipApp(c.App)
	if err != nil {
		log.Fatal(err)
	}
	//lambda
	s := session.New()
	l := lambda.New(s)
	//
	_, err = l.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(c.Name),
	})

	log.Printf("get function: %v", err)

	if err != nil {
		log.Printf("creating function...")
		conf, err := l.CreateFunction(&lambda.CreateFunctionInput{
			Code:         &lambda.FunctionCode{ZipFile: z},
			FunctionName: aws.String(c.Name),
			Handler:      aws.String("handler-unused"),
			Role:         aws.String(c.Role),
			Runtime:      aws.String("provided"),
			Publish:      aws.Bool(true),
			MemorySize:   aws.Int64(128),
			Timeout:      aws.Int64(5),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created: %+v", conf)
	} else {
		log.Printf("updating function code...")
		conf, err := l.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
			ZipFile:      z,
			FunctionName: aws.String(c.Name),
			Publish:      aws.Bool(true),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("updated: %+v", conf)
	}
}

func zipApp(target string) ([]byte, error) {
	//prepare buffer holding a zip "file"
	b := bytes.Buffer{}
	z := zip.NewWriter(&b)
	fh := &zip.FileHeader{
		Name:     "bootstrap",
		Method:   zip.Deflate,
		Modified: time.Now(),
	}
	fh.SetMode(os.ModePerm) //777
	fz, err := z.CreateHeader(fh)
	if err != nil {
		return nil, err
	}
	//check target
	if s, err := os.Stat(target); err != nil {
		return nil, err
	} else if !s.IsDir() {
		return nil, errors.New("target is not a dir")
	}
	//compile target directly into the buffer
	c := exec.Command("go", "build", "-o", "/dev/stdout")
	c.Stdout = fz
	c.Dir = target
	if err := c.Run(); err != nil {
		return nil, errors.New("target build failed")
	}
	log.Printf("compiled: %s", target)
	if err := z.Close(); err != nil {
		return nil, err
	}
	log.Printf("zip file: %d", b.Len())
	return b.Bytes(), nil
}
