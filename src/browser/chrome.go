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
	url            string
	cookies        []*network.CookieParam
}

func NewChromeBrowser(parentCtx context.Context, messageHandler func(string)) *ChromeBrowser {
	// log.Printf("Chrome 옵션 설정")
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

	log.Printf("Chrome 인스턴스 생성")
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
		case *network.EventWebSocketClosed:
			log.Printf("웹소켓 연결이 종료됨")
			// 페이지 새로고침
			go c.reconnect()
		}
	})
}

func (c *ChromeBrowser) reconnect() {
	log.Printf("새로고침 시작")
	err := chromedp.Run(c.ctx,
		chromedp.Reload(),
		chromedp.WaitVisible("body"),
	)

	if err != nil {
		if err == context.Canceled {
			return
		}
		log.Printf("페이지 새로고침 실패: %v", err)
		go c.reconnect()
	} else {
		log.Printf("페이지 새로고침 완료")
	}
}

func (c *ChromeBrowser) Start(url string, cookies []*network.CookieParam) error {
	log.Printf("브라우저 시작 - URL: %s", url)
	c.url = url
	c.cookies = cookies

	err := chromedp.Run(c.ctx,
		network.SetCookies(cookies),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),
	)

	if err != nil {
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
