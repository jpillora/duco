package deploy

import (
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

func Command() opts.Opts {
	return opts.
		New(&deploy{
			l:       lambda.New(session.New()),
			Role:    "arn:aws:iam::652507618334:role/audifree-role",
			Timeout: 5 * time.Minute,
		}).
		Name("deploy")
}

type deploy struct {
	l        *lambda.Lambda
	Name     string        `opts:"help=function name"`
	Role     string        `opts:"help=function role ARN"`
	Memory   int64         `opts:"help=function memory"`
	Timeout  time.Duration `opts:"help=function timeout"`
	AppDir   string        `opts:"mode=arg, help=target <app> to compile"`
	Recreate bool
	//
}

func (d *deploy) Run() error {
	if err := d.run(); err != nil {
		return err
	}
	return nil
}

func (d *deploy) run() error {
	appDir, err := filepath.Abs(d.AppDir)
	if err != nil {
		return err
	}
	d.AppDir = appDir
	if d.Name == "" {
		d.Name = filepath.Base(appDir)
	}
	if d.Memory == 0 {
		d.Memory = 128
	}
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
		FunctionName: aws.String(d.Name),
	})
	exists := err == nil
	deployed := false
	if c := out.Configuration; !d.Recreate && c != nil {
		existingHash := *c.CodeSha256
		deployed = zipHash == existingHash
	}
	//differs? re-deploy
	if deployed {
		log.Printf("function already deployed")
		return nil
	}
	if d.Recreate {
		if err := d.destroy(); err != nil {
			return err
		}
		exists = false
	}
	if exists {
		return d.update(z)
	}
	return d.create(z)
}

func (d *deploy) create(z []byte) error {
	log.Printf("creating function...")
	fn := &lambda.CreateFunctionInput{
		Code:         &lambda.FunctionCode{ZipFile: z},
		FunctionName: aws.String(d.Name),
		Handler:      aws.String("myhandler"),
		Role:         aws.String(d.Role),
		Runtime:      aws.String("provided"),
		Publish:      aws.Bool(true),
		MemorySize:   aws.Int64(d.Memory),
		Timeout:      aws.Int64(int64(d.Timeout.Seconds())),
	}
	if d.Name == "convert" {
		fn.Layers = []*string{
			aws.String("arn:aws:lambda:us-west-1:652507618334:layer:ffmpeg:1"),
			aws.String("arn:aws:lambda:us-west-1:652507618334:layer:ffprobe:1"),
			aws.String("arn:aws:lambda:us-west-1:652507618334:layer:rcrack:1"),
		}
	}
	conf, err := d.l.CreateFunction(fn)
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
		FunctionName: aws.String(d.Name),
		Publish:      aws.Bool(true),
	})
	if err != nil {
		return err
	}
	log.Printf("updated: %+v", conf)
	return nil
}

func (d *deploy) destroy() error {
	log.Printf("deleting function...")
	_, err := d.l.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String(d.Name),
	})
	return err
}
