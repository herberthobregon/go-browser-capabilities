# browser-capabilities

A Go Module that detects `Browser capabilities` and `Client type` from a user agent string.

The following keywords are supported. See [main.ts](https://github.com/herberthobregon/browser-capabilities/blob/master/main.go) for the latest browser support matrix.

> Inspired by [https://www.npmjs.com/package/browser-capabilities](https://www.npmjs.com/package/browser-capabilities) NodeJS package.
## 🎯 Features
### Browser Capabilities
| Keyword       | Description
| :----         | :----
| push          | [HTTP/2 Server Push](https://developers.google.com/web/fundamentals/performance/http2/#server-push)
| serviceworker | [Service Worker API](https://developers.google.com/web/fundamentals/getting-started/primers/service-workers)
| modules       | [JavaScript Modules](https://www.chromestatus.com/feature/5365692190687232) (including dynamic `import()` and `import.meta`)
| es2015        | [ECMAScript 2015 (aka ES6)](https://developers.google.com/web/shows/ttt/series-2/es2015)
| es2016        | ECMAScript 2016
| es2017        | ECMAScript 2017
| es2018        | ECMAScript 2018

### Client Information

```go
type Client struct {
    Browser        string // type: "firefox" | "edge" | "chrome" | "facebook" | "google_app" | "ie" | "safari" | "safari_mobile" | "other" | "vivaldi"
    BrowserVersion float64
    IsMobile       bool
    OS             string // type: "ios" | "android" | "linux" | "mac" | "windows" | "playstation" | "other"
    OSVersion      float64
}
```

## ⚡️ Quickstart
```go
package main

import (
    "fmt"
    "github.com/herberthobregon/browser-capabilities"
)

func main() {
    // CheckES5
    ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36"
    capabilities := browser_capabilities.BrowserCapabilities(ua)
    client := browser_capabilities.GetClient(ua)
    if capabilities["es2015"] || client.Browser == "other" {
        fmt.Println("success ✅")
    } else {
        fmt.Println("Please use Google Chrome > 54 ❌")
    }
}
```

## ⚙️ Installation
Make sure you have Go installed [download](https://golang.org/dl/). Version 1.14 or higher is required.

Initialize your project by creating a folder and then running `go mod init github.com/your/repo` [learn more](https://blog.golang.org/using-go-modules) inside the folder. Then install Fiber with the go get command:
```bash
go get -u github.com/herberthobregon/browser-capabilities
```

