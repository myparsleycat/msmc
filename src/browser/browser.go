// src/browser/browser.go
package browser

import (
	"fmt"
	"log"
	"msmc/src/arcalive"

	"gorm.io/gorm"
)

func Start(db *gorm.DB) {
	messageHandler := func(payload string) {
		handleMessage(db, payload)
	}

	chrome := NewChromeBrowser(messageHandler)
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

func handleMessage(db *gorm.DB, payload string) {
	fmt.Println("메시지 받음", payload)
	if payload == "na" {
		crawler := arcalive.NewCrawler(db)
		created, updated, err := crawler.GetPost("https://arca.live/b/genshinskinmode")
		if err != nil {
			log.Fatalf("크롤링 실패: %v", err)
		}
		fmt.Printf("크롤링 완료! 총계: %d개 생성, %d개 업데이트됨\n", created, updated)
	}
}
