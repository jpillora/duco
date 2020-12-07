package invoke

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jpillora/opts"
)

func Command() opts.Opts {
	return opts.
		New(&invoke{
			l: lambda.New(session.New()),
		}).
		Name("invoke")
}

type invoke struct {
	l       *lambda.Lambda
	Name    string `opts:"mode=arg,help=function name"`
	Async   bool   `opts:""`
	Payload string `opts:"help=invoke payload (defaults to stdin)"`
}

func (i *invoke) Run() error {
	if i.Payload == "" {
		b, _ := ioutil.ReadAll(os.Stdin)
		i.Payload = string(b)
	}
	itype := lambda.InvocationTypeRequestResponse
	if i.Async {
		itype = lambda.InvocationTypeEvent
	}
	t0 := time.Now()
	out, err := i.l.Invoke(&lambda.InvokeInput{
		LogType:        aws.String("Tail"),
		FunctionName:   aws.String(i.Name),
		InvocationType: aws.String(itype),
		Payload:        []byte(i.Payload),
		ClientContext:  testContext(),
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
		fmt.Println(string(p))
	}
	return nil
}

func testContext() *string {
	obj := map[string]string{"foo": "bar"}
	b, _ := json.Marshal(obj)
	s := base64.StdEncoding.EncodeToString(b)
	return aws.String(s)
}
