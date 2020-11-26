package deploy

import (
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

func Command() opts.Opts {
	return opts.
		New(&deploy{
			l:    lambda.New(session.New()),
			Role: "arn:aws:iam::652507618334:role/lambda-role",
		}).
		Name("deploy")
}

type deploy struct {
	l      *lambda.Lambda
	Role   string `opts:"help=role name"`
	AppDir string `opts:"mode=arg, help=target <app> to compile"`
	//
	fnName string
}

func (d *deploy) Run() error {
	appDir, err := filepath.Abs(d.AppDir)
	if err != nil {
		return err
	}
	d.AppDir = appDir
	d.fnName = filepath.Base(appDir)
	log.Printf("compiling %s", d.AppDir)
	//package app into a zip file
	z, err := compileZip(d.AppDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("compiled and zipped %s", d.AppDir)
	//see if we need to deploy
	zipHash := hash(z)
	out, err := d.l.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(d.fnName),
	})
	exists := err == nil
	deployed := false
	if c := out.Configuration; c != nil {
		existingHash := *c.CodeSha256
		deployed = zipHash == existingHash
	}
	//differs? re-deploy
	if deployed {
		return nil
	}
	if exists {
		return d.update(z)
	}
	return d.create(z)
}

func (d *deploy) create(z []byte) error {
	log.Printf("creating function...")
	conf, err := d.l.CreateFunction(&lambda.CreateFunctionInput{
		Code:         &lambda.FunctionCode{ZipFile: z},
		FunctionName: aws.String(d.fnName),
		Handler:      aws.String("myhandler"),
		Role:         aws.String(d.Role),
		Runtime:      aws.String("provided"),
		Publish:      aws.Bool(true),
		MemorySize:   aws.Int64(128),
		Timeout:      aws.Int64(5),
		// Layers: []*string{
		// 	aws.String("arn:aws:lambda:ap-southeast-2:652507618334:layer:duco-bootstrap:1"),
		// },
	})
	if err != nil {
		return err
	}
	log.Printf("created: %+v", conf)
	return nil
}

func (d *deploy) update(z []byte) error {
	log.Printf("updating function code...")
	conf, err := d.l.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		ZipFile:      z,
		FunctionName: aws.String(d.fnName),
		Publish:      aws.Bool(true),
	})
	if err != nil {
		return err
	}
	log.Printf("updated: %+v", conf)
	return nil
}
