package translate

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	LANG_EN = "EN"
	LANG_UA = "UK"
)

type Tranlator struct {
	host       string
	apiKey     string
	defaultUrl string
}

func NewTranlator(host, apiKey string) *Tranlator {
	return &Tranlator{
		host:       host,
		apiKey:     apiKey,
		defaultUrl: "/v2/translate",
	}
}

type tranlatePayload struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type translateResponse struct {
	Translations []translation `json:"translations"`
}

type translation struct {
	Text string `json:"text"`
}

func (t *Tranlator) Tranlate(text, targetLang string) string {
	payload := tranlatePayload{
		Text:       []string{text},
		TargetLang: targetLang,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		pkg.PrintErr(pkg.ErrTranslation, err)
		return ""
	}
	req, _ := http.NewRequest(http.MethodPost, "https://"+t.host+t.defaultUrl, bytes.NewBuffer(body))
	req.Header.Add("Authorization", "DeepL-Auth-Key "+t.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Contnt-Length", fmt.Sprintf("%d", len(body)))
	// req.Header.Add("User-Agent", )

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		pkg.PrintErr(pkg.ErrTranslation, err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		result, _ := io.ReadAll(resp.Body)
		fmt.Println(string(result))

		pkg.PrintErr(pkg.ErrTranslation, fmt.Errorf("%s: %d", "unexpected StatusCode", resp.StatusCode))
		return ""
	}

	responseTranlation := translateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseTranlation)
	if err != nil {
		pkg.PrintErr(pkg.ErrTranslation, err)
		return ""
	}

	return responseTranlation.Translations[0].Text
}
