// src/browser/chrome.go
package browser

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ChromeBrowser struct {
	ctx            context.Context
	cancel         context.CancelFunc
	messageHandler func(string)
}

func NewChromeBrowser(parentCtx context.Context, messageHandler func(string)) *ChromeBrowser {
	log.Printf("Chrome 옵션 설정")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)

	log.Printf("Chrome context 생성")
	allocCtx, cancel := chromedp.NewExecAllocator(parentCtx, opts...)
	ctx, ctxCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	combinedCancel := func() {
		ctxCancel()
		cancel()
	}

	log.Printf("Chrome 브라우저 인스턴스 생성 완료")
	return &ChromeBrowser{
		ctx:            ctx,
		cancel:         combinedCancel,
		messageHandler: messageHandler,
	}
}

func (c *ChromeBrowser) SetupWebSocketListener() {
	chromedp.ListenTarget(c.ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventWebSocketFrameReceived:
			if payload := ev.Response.PayloadData; strings.Contains(payload, "na") {
				c.messageHandler(payload)
			}
		}
	})
}

func (c *ChromeBrowser) Start(url string, cookies []*network.CookieParam) error {
	log.Printf("브라우저 시작 - URL: %s", url)
	err := chromedp.Run(c.ctx,
		network.SetCookies(cookies),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),
	)

	if err != nil {
		// context.Canceled 에러는 정상 종료로 처리
		if err == context.Canceled {
			return nil
		}
		log.Printf("브라우저 시작 실패: %v", err)
		return err
	}

	log.Printf("브라우저 페이지 로드 완료")
	return nil
}

func (c *ChromeBrowser) Close() {
	c.cancel()
}
