package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/toyohashi6140/helpsendgridapi/personalization"
	"github.com/toyohashi6140/helpsendgridapi/request"
)

// APIParam SendGrid APIv3 メール送信用struct
type requestJson struct {
	personalizations [1000]*personalization.Personalization // 宛先や送信内容の親(必須)
	from             map[string]interface{}                 // 送信元(必須)
	subject          string                                 //件名(テンプレートIDといずれか必須)
	content          []map[string]interface{}               //内容(テンプレートIDといずれか必須)
	templateID       string                                 //テンプレートID(subject & contentといずれか必須)
	categories       []string                               //送信カテゴリ(任意)
}

// New インスタンス生成
func New(from map[string]interface{}) requestJson {
	return requestJson{
		personalizations: [1000]*personalization.Personalization{},
		from:             from,
	}
}

// From fromにセットするアドレスをmapに整形して返す
func From(email string) map[string]interface{} {
	return map[string]interface{}{"email": email}
}

// SetSubject 件名セット
func (r *requestJson) SetSubject(subject string) {
	r.subject = subject
}

// SetContent 内容セット
func (r *requestJson) SetContent(contentType, value string) {
	r.content = []map[string]interface{}{
		{
			"type":  contentType,
			"value": value,
		},
	}
}

// AddPersonalization personalizationを追加。1000個要素が入っていたらerrorを返す
func (r *requestJson) AddPersonalizations(p *personalization.Personalization) error {
	for i := 0; i < len(r.personalizations); i++ {
		if r.personalizations[i] == nil {
			r.personalizations[i] = p
			return nil
		}
	}
	return errors.New("all element is not nil")
}

// ResetPersonalization requestJsonのpersonalizationsを初期化(APIの制約で上限1,000件なのでそれを超えそうな場合に使用)
func (r *requestJson) ResetPersonalization() {
	r.personalizations = [1000]*personalization.Personalization{}
}

// SetTemplateID template_idを設定(e介護Withコンテンツ用)
func (r *requestJson) SetTemplateID(id string) {
	// テンプレート選定
	r.templateID = id
}

// SetCategories categoriesを設定(最大10カテゴリ)
func (r *requestJson) SetCategories(categories ...string) error {
	if len(categories) > 10 {
		return errors.New("設定したカテゴリの数が多すぎます")
	}
	r.categories = categories
	return nil
}

// makeJson 設定されたパラメータ(requestJson)をもとにAPIリクエスト用のjsonにするmapを作成
func (r *requestJson) makeMap() map[string]interface{} {
	// personalizationの処理
	maps := []map[string]interface{}{}
	for _, p := range r.personalizations {
		// 1000件以下でリクエストが送られる場合に、配列の要素がnilとなる場合がある
		if p != nil {
			m := map[string]interface{}{
				"to":                    p.To(),
				"dynamic_template_data": p.DynamicTemplateData(),
			}
			maps = append(maps, m)
		}
	}

	jsonMap := map[string]interface{}{
		"personalizations": maps,
		"from":             r.from,
	}

	// テンプレートIDはcontentsと選択なので後から制御する
	if r.templateID != "" {
		jsonMap["template_id"] = r.templateID
	} else {
		jsonMap["subject"] = r.subject
		jsonMap["content"] = r.content
	}

	// カテゴリは任意パラメータなので後から制御する
	if r.categories != nil {
		jsonMap["categories"] = r.categories
	}
	return jsonMap
}

// RequestPost POSTリクエストを送信する
func (r requestJson) RequestPost(apiKey string) error {
	req := request.New(apiKey).Post()

	jsonByte, err := json.Marshal(r.makeMap())
	if err != nil {
		return err
	}
	req.Body = jsonByte

	res, err := req.Do()
	if err != nil {
		return err
	}

	if res.Body != "" {
		fmt.Println(res.Body)
	}
	for key, vals := range res.Headers {
		fmt.Printf("%s: %s\n", key, strings.Join(vals, ","))
	}
	return nil
}
