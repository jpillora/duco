# duco

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

