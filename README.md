# Usage

```go
package main

import (
	"fmt"

	"github.com/bulatsan/go-gracectx"
)

func main() {
	ctx, cancel := gracectx.New()
	defer cancel()

	fmt.Println("waiting for signal")
	<-ctx.Done()
}
```