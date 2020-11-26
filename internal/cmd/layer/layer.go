package layer

import (
	"log"

	"github.com/jpillora/duco/internal/cmd/deploy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

func Command() opts.Opts {
	return opts.
		New(&layer{
			l: lambda.New(session.New()),
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
		LayerName: aws.String("duco-bootstrap"),
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
