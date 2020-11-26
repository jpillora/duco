# gambda

### Development

```go
package main

func main() {
    g := gambda.New()
    g.Add(myhandler.New())
    g.Start()
}
```

### Deployment

```sh
#!/bin/bash
gambda deploy <path-to-my-handler>
```

