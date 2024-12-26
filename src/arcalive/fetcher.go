package arcalive

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type FetcherProps struct {
	URL    string
	Method string
}

func Fetch(props FetcherProps) (*http.Response, error) {
	req, err := http.NewRequest(props.Method, props.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("request 생성 실패: %v", err)
	}

	as0 := os.Getenv("ARCA_SESSION2")
	as1 := os.Getenv("ARCA_SESSION2_SIG")

	req.Header.Set("Cookie", fmt.Sprintf("arca.session2=%s; arca.session2.sig=%s", as0, as1))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://arca.live/b/genshinskinmode")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP 요청 실패 - 상태 코드: %d, 응답: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
