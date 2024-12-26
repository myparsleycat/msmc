// src/browser/browser.go
package browser

import (
	"fmt"
	"log"
	"msmc/src/arcalive"
	"msmc/src/database"
)

func Start() {
	chrome := NewChromeBrowser(handleMessage)
	defer chrome.Close()

	chrome.SetupWebSocketListener()
	cookies := GetArcaCookies()

	errCh := make(chan error, 1)

	go func() {
		errCh <- chrome.Start("https://arca.live/b/genshinskinmode", cookies)
	}()

	if err := <-errCh; err != nil {
		log.Printf("브라우저 오류 발생: %v", err)
		return
	}
	log.Printf("브라우저 시작 완료")
}

func handleMessage(payload string) {
	fmt.Println("메시지 받음", payload)
	if payload == "na" {
		crawler := arcalive.NewCrawler(database.GetDB())
		created, updated, err := crawler.GetPost("https://arca.live/b/genshinskinmode")
		if err != nil {
			log.Fatalf("크롤링 실패: %v", err)
		}
		fmt.Printf("크롤링 완료! 총계: %d개 생성, %d개 업데이트됨\n", created, updated)
	}
}
