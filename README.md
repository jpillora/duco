# duco

⚠️ work in progress

### Development

```go
package main

func main() {
    g := duco.New()
    g.Add(myhandler.New())
    g.Start()
}
```

### Deployment

```sh
#!/bin/bash
duco deploy <path-to-my-handler>
```

### CF

https://godoc.org/github.com/awslabs/goformation/cloudformation/lambda#Function
