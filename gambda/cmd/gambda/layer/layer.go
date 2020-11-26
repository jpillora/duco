package layer

import (
	"gambda/cmd/gambda/deploy"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

func Command() opts.Opts {
	return opts.
		New(&layer{
			l: lambda.New(session.New()),
			// Name: "go-raw-runtime",
			// Role: "arn:aws:iam::652507618334:role/lambda-role",
		}).
		Name("layer")
}

type layer struct {
	l             *lambda.Lambda
	BootstrapPath string
}

func (l *layer) Run() error {

	z, err := deploy.CompileZip(l.BootstrapPath, false)
	if err != nil {
		return err
	}

	out, err := l.l.PublishLayerVersion(&lambda.PublishLayerVersionInput{
		LayerName: aws.String("gambda-bootstrap"),
		Content: &lambda.LayerVersionContentInput{
			ZipFile: z,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("published layer %s (%d bytes, %s)",
		*out.Content.Location,
		*out.Content.CodeSize,
		*out.Content.CodeSha256,
	)
	return nil
}
