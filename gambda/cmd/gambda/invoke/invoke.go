package invoke

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

func Command() opts.Opts {
	return opts.
		New(&invoke{
			l:       lambda.New(session.New()),
			Name:    "go-raw-runtime",
			Payload: `{"hello":"world"}`,
		}).
		Name("invoke")
}

type invoke struct {
	l       *lambda.Lambda
	Name    string `opts:"help=function name"`
	Payload string `opts:"help=invoke payload"`
}

func (i *invoke) Run() error {
	t0 := time.Now()
	out, err := i.l.Invoke(&lambda.InvokeInput{
		LogType:       aws.String("Tail"),
		FunctionName:  aws.String(i.Name),
		Payload:       []byte(i.Payload),
		ClientContext: testContext(),
	})
	if err != nil {
		return err
	}
	log.Printf("invoked: %s (took %s)", i.Name, time.Since(t0))
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

func testContext() *string {
	obj := map[string]string{"foo": "bar"}
	b, _ := json.Marshal(obj)
	s := base64.StdEncoding.EncodeToString(b)
	return aws.String(s)
}
