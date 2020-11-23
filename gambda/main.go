package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	Name    string `opts:"help=function name"`
	Role    string `opts:"help=role name"`
	App     string `opts:"mode=arg, help=target <app> to compile"`
	Payload string
}

var c = cli{
	Name: "go-raw-runtime",
	Role: "arn:aws:iam::652507618334:role/lambda-role",
}

var s = session.New()
var l = lambda.New(s)

func main() {

	opts.New(&c).Name("gambda").Parse()
	//
	z, err := pacakgeApp(c.App)
	if err != nil {
		log.Fatal(err)
	}
	zipHash := hash(z)

	out, err := l.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(c.Name),
	})
	exists := err == nil
	deployed := false
	if c := out.Configuration; c != nil {
		existingHash := *c.CodeSha256
		deployed = zipHash == existingHash
	}

	if !deployed {
		if err := deployApp(exists, z); err != nil {
			log.Fatal(err)
		}
	}

	if err := invokeApp(c.Payload); err != nil {
		log.Fatal(err)
	}
}

func pacakgeApp(target string) ([]byte, error) {
	//prepare buffer holding a zip "file"
	b := bytes.Buffer{}
	z := zip.NewWriter(&b)
	fh := &zip.FileHeader{
		Name:   "bootstrap",
		Method: zip.Deflate,
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
	c := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-v", "-o", "/dev/stdout")
	c.Env = append(os.Environ(), "GOOS=linux")
	c.Stderr = os.Stderr
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

func deployApp(exists bool, z []byte) error {
	if !exists {
		log.Printf("creating function...")
		conf, err := l.CreateFunction(&lambda.CreateFunctionInput{
			Code:         &lambda.FunctionCode{ZipFile: z},
			FunctionName: aws.String(c.Name),
			Handler:      aws.String("myhandler"),
			Role:         aws.String(c.Role),
			Runtime:      aws.String("provided"),
			Publish:      aws.Bool(true),
			MemorySize:   aws.Int64(128),
			Timeout:      aws.Int64(5),
		})
		if err != nil {
			return err
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
			return err
		}
		log.Printf("updated: %+v", conf)
	}
	return nil
}

func invokeApp(payload string) error {

	if payload == "" {
		payload = `{"hello":"world"}`
	}

	t0 := time.Now()
	out, err := l.Invoke(&lambda.InvokeInput{
		LogType:       aws.String("Tail"),
		FunctionName:  aws.String(c.Name),
		Payload:       []byte(payload),
		ClientContext: testContext(),
	})
	if err != nil {
		return err
	}
	log.Printf("invoked: %s (took %s)", c.Name, time.Since(t0))
	if r := out.LogResult; r != nil {
		if b, err := base64.StdEncoding.DecodeString(*r); err == nil {
			log.Printf("logs:\n%s", string(b))
		}
	}
	if p := out.Payload; p != nil {
		log.Printf("payload: %s", string(p))
	}
	return nil
}

func hash(z []byte) string {
	h := sha256.New()
	h.Write(z)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func testContext() *string {
	obj := map[string]string{"foo": "bar"}
	b, _ := json.Marshal(obj)
	s := base64.StdEncoding.EncodeToString(b)
	return aws.String(s)
}
