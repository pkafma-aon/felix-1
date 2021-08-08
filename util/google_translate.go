package util

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GoogleTranslate(text string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 60}

	info := whatlanggo.Detect(text)
	toLang := "zh-cn"
	sourceLang := "en"
	if info.Lang == whatlanggo.Cmn {
		toLang = "en"
		sourceLang = "zh-cn"
	}
	//url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=ene&tl=zh-cn&dt=t&q=%s", url.QueryEscape(text))

	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s", sourceLang, toLang, url.QueryEscape(text))
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("google translate api status code not 200:::" + url)
	}
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0], nil
}
