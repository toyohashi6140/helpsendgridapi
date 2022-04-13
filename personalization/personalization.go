package personalization

type to [1]map[string]interface{} // 2022.3.5のインシデントを受け、2つ以上宛先が入らないように長さ1の配列に設定

func To(email string) to {
	return [1]map[string]interface{}{{"email": email}}
}

type dynaminTemplateData map[string]interface{} // 動的テンプレート変数もpersonalizationに紐づくため定義

func DynaminTemplateData(m map[string]interface{}) dynaminTemplateData {
	return m
}

// Personalization personalizationパラメータを定義する構造体
type Personalization struct {
	to     to
	dytmpl dynaminTemplateData
}

// New 必須の宛先情報と、最大長1,000件を付与してpersonalizationインスタンスを生成
func New(to to, d dynaminTemplateData) *Personalization {
	return &Personalization{to, d}
}

// To toのgetter
func (p *Personalization) To() to {
	return p.to
}

// SetTo to info をセット。 Newで初期化せずに構造体を定義する人用
func (p *Personalization) SetTo(to to) {
	p.to = to
}

// DynamicTemplateData dynaminTemplateDataのgetter
func (p *Personalization) DynamicTemplateData() dynaminTemplateData {
	return p.dytmpl
}

// SetDynamicTemplateData dynaminTemplateData(任意パラメータ)をセット
func (p *Personalization) SetDynamicTemplateData(d dynaminTemplateData) {
	p.dytmpl = d
}
