package ajax

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kere/gos"
)

var (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"
)

type Ajax struct {
	prepareUrl       string
	cookies          *http.Cookie
	cookiesUpdatedAt time.Time
}

func NewAjax(url string) *Ajax {
	return &Ajax{prepareUrl: url}
}

func (a *Ajax) PreparePageCookie() error {
	if a.prepareUrl == "" {
		return nil
	}

	if a.cookies != nil && gos.NowInLocation().Before(a.cookiesUpdatedAt.AddDate(0, 0, 10)) {
		return nil
	}

	resp, err := a.Get(a.prepareUrl)
	if err != nil {
		return gos.DoError(err)
	}
	l := resp.Cookies()

	if len(l) > 0 {
		a.cookies = l[0]
		a.cookiesUpdatedAt = gos.NowInLocation()
		fmt.Println("set cookies", l[0])
	}
	return nil
}

// Wget 抓取数据
func (a *Ajax) Get(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if req == nil {
		return nil, fmt.Errorf("failed get url:%s", url)
	}
	req.Header.Add("User-Agent", userAgent)
	if a.cookies != nil {
		req.AddCookie(a.cookies)
	}
	gos.Log.Info("wget", url)
	resp, err := client.Do(req)

	if err != nil {
		return nil, gos.DoError("error:", err)
	} else if resp.Body == nil {
		return nil, gos.DoError("error: resp.Body is empty")
	}

	if resp.StatusCode != 200 {
		return resp, gos.DoError(fmt.Sprintf("Get failed:%d", resp.StatusCode))
	}
	// var OKPJKmpr={pages:0,data:[{stats:false}]}
	return resp, nil
}

func (a *Ajax) GetBody(url string) ([]byte, error) {
	resp, err := a.Get(url)
	if err != nil {
		return nil, gos.DoError(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, gos.DoError(err)
	}

	return body, nil
}

func (a *Ajax) Post(url string, val url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(val.Encode()))
	if err != nil {
		return nil, gos.DoError(err)
	}
	req.Header.Add("User-Agent", userAgent)
	if a.cookies != nil {
		req.AddCookie(a.cookies)
	}

	client := &http.Client{}
	gos.Log.Info("post", url, val)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, gos.DoError(err)
	}

	if resp.StatusCode != 200 {
		return resp, gos.DoError(fmt.Sprintf("Post failed:%d", resp.StatusCode))
	}

	return resp, nil
}

func (a *Ajax) PostBody(url string, val url.Values) ([]byte, error) {
	resp, err := a.Post(url, val)
	if err != nil {
		return nil, gos.DoError(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, gos.DoError(err)
	}

	return body, nil
}
