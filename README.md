# browser-capabilities

A Go Module that detects `Browser capabilities` and `Client type` from a user agent string.

The following keywords are supported. See [main.ts](https://github.com/herberthobregon/browser-capabilities/blob/master/main.go) for the latest browser support matrix.

> Inspired by [https://www.npmjs.com/package/browser-capabilities](https://www.npmjs.com/package/browser-capabilities) NodeJS package.
## Browser Capabilities
| Keyword       | Description
| :----         | :----
| push          | [HTTP/2 Server Push](https://developers.google.com/web/fundamentals/performance/http2/#server-push)
| serviceworker | [Service Worker API](https://developers.google.com/web/fundamentals/getting-started/primers/service-workers)
| modules       | [JavaScript Modules](https://www.chromestatus.com/feature/5365692190687232) (including dynamic `import()` and `import.meta`)
| es2015        | [ECMAScript 2015 (aka ES6)](https://developers.google.com/web/shows/ttt/series-2/es2015)
| es2016        | ECMAScript 2016
| es2017        | ECMAScript 2017
| es2018        | ECMAScript 2018

## Client Information

```go
type Client struct {
	Browser        string // type: "firefox" | "edge" | "chrome" | "facebook" | "google_app" | "ie" | "safari" | "safari_mobile" | "other" | "vivaldi"
	BrowserVersion float64
	IsMobile       bool
	OS             string // type: "ios" | "android" | "linux" | "mac" | "windows" | "playstation" | "other"
	OSVersion      float64
}
```