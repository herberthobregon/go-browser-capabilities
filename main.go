package browser_capabilities

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Client struct {
	Browser        string // type: "firefox" | "edge" | "chrome" | "facebook" | "google_app" | "ie" | "safari" | "safari_mobile" | "other" | "vivaldi"
	BrowserVersion float64
	IsMobile       bool
	OS             string // type: "ios" | "android" | "linux" | "mac" | "windows" | "playstation" | "other"
	OSVersion      float64
}

type OS struct {
	Regex string
	OS    string // type: "ios" | "android" | "linux" | "mac" | "windows" | "playstation" | "other"
	Apply func(c string) string
}

type Browser struct {
	Browser string // type: "firefox" | "edge" | "chrome" | "facebook" | "google_app" | "ie" | "safari" | "safari_mobile" | "other" | "vivaldi"
	Regex   string
}
/*
Get Client information from a user agent string.
*/
func GetClient(ua string) Client {
	o := Client{
		Browser:        "other",
		BrowserVersion: 0,
		IsMobile:       false,
		OS:             "other",
		OSVersion:      0,
	}
	a := []Browser{
		Browser{Browser: "facebook", Regex: `FBAV\/([0-9\.]+)`},
		{Browser: "facebook", Regex: `FBSV\/([0-9\.]+)`},
		{Browser: "google_app", Regex: `GSA\/([0-9\.]+)`},
		{Browser: "chrome", Regex: `Chrome\/([0-9\.]+)`},
		{Browser: "chrome", Regex: `CriOS\/([0-9\.]+)`},
		{Browser: "firefox", Regex: `Firefox\/([0-9\.]+)`},
		{Browser: "firefox", Regex: `FxiOS\/([0-9\.]+)`},
		{Browser: "safari", Regex: `Version\/([0-9\.]+).+Safari`},
		{Browser: "edge", Regex: `Edge\/([0-9\.]+)`},
		{Browser: "ie", Regex: `Trident\/.+rv[: ]?([0-9]+)`},
		{Browser: "vivaldi", Regex: `Vivaldi\/([0-9\.]+)`},
	}
	for i := 0; i < len(a); i++ {
		re := regexp.MustCompile(a[i].Regex)
		if ok := re.FindAllStringSubmatch(ua, -1); ok != nil {
			o.Browser = a[i].Browser
			f, err := strconv.ParseInt(strings.Split(ok[0][1], ".")[0], 10, 32)
			if err == nil {
				o.BrowserVersion = float64(f)
			}
			break
		}
	}

	if _, err := regexp.MatchString(`Mobile`, ua); err != nil {
		o.IsMobile = true
	}

	opSys := []OS{
		OS{
			OS:    "ios",
			Regex: `([0-9_]+) like Mac OS X`,
			Apply: versionFixer,
		},
		OS{
			OS:    "ios",
			Regex: `CPU like Mac OS X`,
			Apply: func(c string) string {
				return "0"
			},
		},
		OS{
			OS:    "android",
			Regex: `Android ([0-9\.]+)`,
			Apply: func(v string) string { return v },
		},
		OS{
			OS:    "linux",
			Regex: `Linux (x86_64|x86)`,
			Apply: func(v string) string {
				if v == "x86_64" {
					return "64"
				} else {
					return "32"
				}
			},
		},
		OS{
			OS:    "mac",
			Regex: `Macintosh.+Mac OS X ([0-9_]+)`,
			Apply: versionFixer,
		},
		OS{
			OS:    "windows",
			Regex: `Windows NT ([0-9\.]+)`,
			Apply: versionFixer,
		},
		OS{
			OS:    "playstation",
			Regex: `PlayStation ([0-9\.]+)`,
			Apply: versionFixer,
		},
	}

	for i := 0; i < len(opSys); i++ {
		re := regexp.MustCompile(opSys[i].Regex)
		if ok := re.FindAllStringSubmatch(ua, -1); ok != nil {
			o.OS = opSys[i].OS
			num, err := strconv.ParseFloat(opSys[i].Apply(ok[0][1]), 64)
			if err == nil {
				o.OSVersion = num
			} else {
				o.OSVersion = 0
			}
			break
		}
	}
	return o
}

type FeautureMap map[string]func(ua Client) bool
type BrowserMap map[string]FeautureMap

func versionFixer(v string) string {
	parts := strings.Split(strings.Replace(v, "_", ".", -1), ".")
	if len(parts) >= 2 {
		return parts[0] + `.` + parts[1]
	} else if len(parts) == 1 {
		return parts[0]
	}
	return "0"
}

var chrome = FeautureMap{
	"es2015":        since(49),
	"es2016":        since(58),
	"es2017":        since(58),
	"es2018":        since(64),
	"push":          since(41),
	"serviceworker": since(45),
	"modules":       since(64),
}
var browserPredicates = BrowserMap{
	"chrome": chrome,
	"opera": FeautureMap{
		"es2015":        since(36),
		"es2016":        since(45),
		"es2017":        since(45),
		"es2018":        since(51),
		"push":          since(28),
		"serviceworker": since(32),
		"modules":       since(48),
	},
	"vivaldi": FeautureMap{
		"es2015":        since(1),
		"es2016":        since(1, 14),
		"es2017":        since(1, 14),
		"es2018":        since(1, 14),
		"push":          since(1),
		"serviceworker": since(1),
		"modules":       since(1, 14),
	},
	"safari_mobile": FeautureMap{
		"es2015":        sinceOS(10),
		"es2016":        sinceOS(10, 3),
		"es2017":        sinceOS(10, 3),
		"es2018":        func(ua Client) bool { return false },
		"push":          sinceOS(9, 2),
		"serviceworker": sinceOS(11, 3),
		"modules":       sinceOS(11, 3),
	},
	"safari": FeautureMap{
		"es2015": since(10),
		"es2016": since(10, 1),
		"es2017": since(10, 1),
		"es2018": func(ua Client) bool { return false },
		"push": func(ua Client) bool {
			return versionAtLeast([]float64{9}, parseVersion(ua.BrowserVersion)) &&
				// HTTP/2 on desktop Safari requires macOS 10.11 according to
				// caniuse.com.
				versionAtLeast([]float64{10, 11}, parseVersion(ua.OSVersion))
		},
		// https://webkit.org/status/#specification-service-workers
		"serviceworker": since(11, 1),
		"modules":       since(11, 1),
	},
	"edge": FeautureMap{
		// Edge versions before 15.15063 may contain a JIT bug affecting ES6
		// constructors (https://github.com/Microsoft/ChakraCore/issues/1496).
		// Since this bug was fixed after es2016 and 2017 support, all these
		// versions are the same.
		"es2015": since(15, 15063),
		"es2016": since(15, 15063),
		"es2017": since(15, 15063),
		"es2018": func(ua Client) bool { return false },
		"push":   since(12),
		// https://developer.microsoft.com/en-us/microsoft-edge/platform/status/serviceworker/
		"serviceworker": func(ua Client) bool { return false },
		"modules":       func(ua Client) bool { return false },
	},
	"firefox": FeautureMap{
		"es2015": since(51),
		"es2016": since(52),
		"es2017": since(52),
		"es2018": since(58),
		// Firefox bug - https://bugzilla.mozilla.org/show_bug.cgi?id=1409570
		"push":          since(63),
		"serviceworker": since(44),
		"modules":       since(67),
	},
}

/*
Detects browser capabilities from a user agent string.
*/
func BrowserCapabilities(userAgent string) map[string]bool {
	ua := GetClient(userAgent)
	var capabilities map[string]bool = map[string]bool{}
	browserName := ua.Browser
	if ua.OS == "ios" {
		// if iOS is really safari_mobile.
		browserName = "safari_mobile"
	} else if ua.Browser == "facebook" {
		browserName = "chrome"
	}
	predicates := browserPredicates[browserName]
	for k, capability := range predicates {
		if capability(ua) {
			capabilities[k] = true
		}
	}
	return capabilities
}

/*
Return whether `version` is at least as high as `atLeast`.
*/
func versionAtLeast(atLeast []float64, version []float64) bool {
	for i := float64(0); i < float64(len(atLeast)); i++ {
		r := atLeast[int32(i)]
		var v float64 = 0
		if float64(len(version)) > i {
			v = version[int32(i)]
		}
		if v > r {
			return true
		}
		if v < r {
			return false
		}
	}
	return true
}

func parseVersion(version float64) []float64 {
	if version <= 0 {
		return []float64{}
	}
	rsp := []float64{}
	s := strings.Split(fmt.Sprintf("%f", version), ".")
	for _, v := range s {
		p, err := strconv.ParseFloat(v, 64)
		if err != nil {
			rsp = append(rsp, -1)
		} else {
			rsp = append(rsp, p)
		}
	}
	return rsp
}

/*
Make a predicate that checks if the browser version is at least this high.
*/
func since(atLeast ...float64) func(ua Client) bool {
	return func(ua Client) bool { return versionAtLeast(atLeast, parseVersion(ua.BrowserVersion)) }
}

/*
Make a predicate that checks if the OS version is at least this high.
*/
func sinceOS(atLeast ...float64) func(ua Client) bool {
	return func(ua Client) bool { return versionAtLeast(atLeast, parseVersion(ua.OSVersion)) }
}
