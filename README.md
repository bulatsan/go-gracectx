# Usage

```go
package main

import (
	"context"
	"errors"
	"github.com/bulatsan/go-gracectx"
	"net/http"
)

func main() {
    ctx, cancel := gracectx.New()
    defer cancel()

	srv := &http.Server{
		Addr: ":8000",
	}
	go func() {
		err := srv.ListenAndServe()
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
}
```