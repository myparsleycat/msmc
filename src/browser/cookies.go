// src/browser/cookies.go
package browser

import (
	"os"

	"github.com/chromedp/cdproto/network"
)

func GetArcaCookies() []*network.CookieParam {
	as0 := os.Getenv("ARCA_SESSION2")
	as1 := os.Getenv("ARCA_SESSION2_SIG")

	return []*network.CookieParam{
		{
			Name:   "arca.session2",
			Value:  as0,
			Domain: "arca.live",
			Path:   "/",
		},
		{
			Name:   "arca.session2.sig",
			Value:  as1,
			Domain: "arca.live",
			Path:   "/",
		},
	}
}
