package deploy

import (
	"encoding/json"
	"io"

	"golang.org/x/sync/errgroup"
)

type goList struct {
	Name       string
	ImportPath string
	Module     struct {
		Path  string
		Main  bool
		GoMod string
	}
}

func list(goDir string) (l goList, err error) {
	r, w := io.Pipe()
	eg := errgroup.Group{}
	eg.Go(func() error {
		return json.NewDecoder(r).Decode(&l)
	})
	eg.Go(func() error {

		err := goExec(w, goDir, "list", "-json")
		w.Close()
		return err
	})
	err = eg.Wait()
	return
}

// "Module": {
// 	"Path": "github.com/jpillora/duco",
// 	"Main": true,
// 	"Dir": "/Users/jpillora/Code/Go/src/github.com/jpillora/duco",
// 	"GoMod": "/Users/jpillora/Code/Go/src/github.com/jpillora/duco/go.mod",
// 	"GoVersion": "1.13"
// }

/*
{
        "Dir": "/Users/jpillora/Code/Go/src/github.com/jpillora/duco-example/tmp/user",
        "ImportPath": "myhandler/tmp/user",
        "Name": "main",
        "Target": "/Users/jpillora/Code/Go/bin/linux_amd64/user",
        "Root": "/Users/jpillora/Code/Go",
        "Module": {
                "Path": "myhandler",
                "Main": true,
                "Dir": "/Users/jpillora/Code/Go/src/github.com/jpillora/duco-example",
                "GoMod": "/Users/jpillora/Code/Go/src/github.com/jpillora/duco-example/go.mod",
                "GoVersion": "1.15"
        },
        "Match": [
                "."
        ],
        "Stale": true,
        "StaleReason": "stale dependency: github.com/jpillora/duco/runtime",
        "GoFiles": [
                "main.go"
        ],
        "Imports": [
                "github.com/jpillora/duco/runtime",
                "myhandler"
        ],
        "Deps": [
                "bufio",
                "bytes",
                "compress/flate",
                "compress/gzip",
                "container/list",
                "context",
                "crypto",
                "crypto/aes",
                "crypto/cipher",
                "crypto/des",
                "crypto/dsa",
                "crypto/ecdsa",
                "crypto/ed25519",
                "crypto/ed25519/internal/edwards25519",
                "crypto/elliptic",
                "crypto/hmac",
                "crypto/internal/randutil",
                "crypto/internal/subtle",
                "crypto/md5",
                "crypto/rand",
                "crypto/rc4",
                "crypto/rsa",
                "crypto/sha1",
                "crypto/sha256",
                "crypto/sha512",
                "crypto/subtle",
                "crypto/tls",
                "crypto/x509",
                "crypto/x509/pkix",
                "encoding",
                "encoding/asn1",
                "encoding/base64",
                "encoding/binary",
                "encoding/hex",
                "encoding/json",
                "encoding/pem",
                "errors",
                "fmt",
                "github.com/jpillora/duco",
                "github.com/jpillora/duco/runtime",
                "hash",
                "hash/crc32",
                "internal/bytealg",
                "internal/cpu",
                "internal/fmtsort",
                "internal/nettrace",
                "internal/oserror",
                "internal/poll",
                "internal/race",
                "internal/reflectlite",
                "internal/singleflight",
                "internal/syscall/execenv",
                "internal/syscall/unix",
                "internal/testlog",
                "internal/unsafeheader",
                "io",
                "io/ioutil",
                "log",
                "math",
                "math/big",
                "math/bits",
                "math/rand",
                "mime",
                "mime/multipart",
                "mime/quotedprintable",
                "myhandler",
                "net",
                "net/http",
                "net/http/httptrace",
                "net/http/internal",
                "net/textproto",
                "net/url",
                "os",
                "path",
                "path/filepath",
                "reflect",
                "runtime",
                "runtime/internal/atomic",
                "runtime/internal/math",
                "runtime/internal/sys",
                "sort",
                "strconv",
                "strings",
                "sync",
                "sync/atomic",
                "syscall",
                "time",
                "unicode",
                "unicode/utf16",
                "unicode/utf8",
                "unsafe",
                "vendor/golang.org/x/crypto/chacha20",
                "vendor/golang.org/x/crypto/chacha20poly1305",
                "vendor/golang.org/x/crypto/cryptobyte",
                "vendor/golang.org/x/crypto/cryptobyte/asn1",
                "vendor/golang.org/x/crypto/curve25519",
                "vendor/golang.org/x/crypto/hkdf",
                "vendor/golang.org/x/crypto/internal/subtle",
                "vendor/golang.org/x/crypto/poly1305",
                "vendor/golang.org/x/net/dns/dnsmessage",
                "vendor/golang.org/x/net/http/httpguts",
                "vendor/golang.org/x/net/http/httpproxy",
                "vendor/golang.org/x/net/http2/hpack",
                "vendor/golang.org/x/net/idna",
                "vendor/golang.org/x/sys/cpu",
                "vendor/golang.org/x/text/secure/bidirule",
                "vendor/golang.org/x/text/transform",
                "vendor/golang.org/x/text/unicode/bidi",
                "vendor/golang.org/x/text/unicode/norm"
        ]
}
*/
