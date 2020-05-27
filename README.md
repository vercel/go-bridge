# go-bridge

Bridge for `@vercel/go`.

## Usage

```go
package main

import (
	"net/http"
	bridge "github.com/vercel/go-bridge/go/bridge"
)

func main() {
	bridge.Start(http.HandlerFunc(__NOW_HANDLER_FUNC_NAME))
}

```

See [PR 3976](https://github.com/vercel/vercel/pull/3976) to see why this repo exists.
