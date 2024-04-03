# A rate limiter with a fixed window

## quick start

```bash
go get github.com/aide-family/rate
```

```go
package main

import (
	"fmt"
	"github.com/aide-family/rate"
	"time"
)

func main() {
	// This limiter allows for the generation of 3 tokens within a 10 second window with a minimum interval of 2 seconds.
	limiter := rate.NewLimiter(3, 10*time.Second, 2*time.Second)
	if limiter.Allow() {
		fmt.Println("allow")
	}else{
		fmt.Println("not allow")
    }
}
```

