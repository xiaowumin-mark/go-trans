package bing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bytedance/sonic"
)

// 目前只收集到了bing的批量翻译api，此api用于edge的网页翻译功能

type BatchTranslateResp struct {
	Raw    string
	Parsed []BatchTranslateResult
}

type BatchTranslateResult struct {
	Translations []struct {
		Text    string `json:"text"`
		To      string `json:"to"`
		SentLen struct {
			SrcSentLen   []int `json:"srcSentLen"`
			TransSentLen []int `json:"transSentLen"`
		} `json:"sentLen"`
	} `json:"translations"`
}

/*
此Api是bing的批量翻译api，适用于大量的翻译任务
from 和 to参数可以传入语言代码或者语言名称，如："zh-Hans"、"zh-Hant"、"en"、"ja" "zh-cn"
*/
func BatchTranslate(text []string, from string, to string) (*BatchTranslateResp, error) {

	//https://edge.microsoft.com/translate/translatetext?from=en-us&to=zh&isEnterpriseClient=false

	bodyJson, err := sonic.Marshal(text)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://edge.microsoft.com/translate/translatetext?from="+from+"&to="+to+"&isEnterpriseClient=false", bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	req.Header.Set("x-edge-shopping-flag", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 查看状态码
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []BatchTranslateResult
	err = sonic.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &BatchTranslateResp{Raw: string(body), Parsed: result}, nil

}
