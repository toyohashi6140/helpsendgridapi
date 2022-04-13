package json

import (
	"reflect"
	"testing"

	"github.com/toyohashi6140/helpsendgridapi/personalization"
)

// Test_requestJson_makeMap すべてのパラメータをセットした場合、及び分岐点となるtemplate_idがない場合/ある場合/カテゴリがない場合についてテスト
func Test_requestJson_makeMap(t *testing.T) {
	type fields struct {
		personalizations [1000]*personalization.Personalization
		from             map[string]interface{}
		subject          string
		content          []map[string]interface{}
		templateID       string
		categories       []string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
		{
			name: "all parameters set",
			fields: fields{
				[1000]*personalization.Personalization{
					personalization.New(
						personalization.To("to@gmail.com"),
						personalization.DynaminTemplateData(
							map[string]interface{}{"firstname": "toyohashi", "lastname": 6140},
						),
					),
				},
				map[string]interface{}{"email": "from@gmail.com"},
				"subject",
				[]map[string]interface{}{{"type": "text/plain", "content": "メールのテストです"}},
				"aaaaaaaaaaaaaaa",
				[]string{"a", "b", "c"},
			},
			want: map[string]interface{}{
				"from":        map[string]interface{}{"email": "from@gmail.com"},
				"template_id": "aaaaaaaaaaaaaaa",
				"categories":  []string{"a", "b", "c"},
			},
		},
		{
			name: "template_id isn't set",
			fields: fields{
				personalizations: [1000]*personalization.Personalization{
					personalization.New(
						personalization.To("to@gmail.com"),
						personalization.DynaminTemplateData(
							map[string]interface{}{"firstname": "toyohashi", "lastname": 6140},
						),
					),
				},
				from:       map[string]interface{}{"email": "from@gmail.com"},
				subject:    "subject",
				content:    []map[string]interface{}{{"type": "text/plain", "content": "メールのテストです"}},
				categories: []string{"a", "b", "c"},
			},
			want: map[string]interface{}{
				"from":       map[string]interface{}{"email": "from@gmail.com"},
				"subject":    "subject",
				"content":    []map[string]interface{}{{"type": "text/plain", "content": "メールのテストです"}},
				"categories": []string{"a", "b", "c"},
			},
		},
		{
			name: "no categories",
			fields: fields{
				personalizations: [1000]*personalization.Personalization{
					personalization.New(
						personalization.To("to@gmail.com"),
						personalization.DynaminTemplateData(
							map[string]interface{}{"firstname": "toyohashi", "lastname": 6140},
						),
					),
				},
				from:       map[string]interface{}{"email": "from@gmail.com"},
				subject:    "subject",
				content:    []map[string]interface{}{{"type": "text/plain", "content": "メールのテストです"}},
				templateID: "aaaaaaaaaaaaaaa",
			},
			want: map[string]interface{}{
				"from":        map[string]interface{}{"email": "from@gmail.com"},
				"template_id": "aaaaaaaaaaaaaaa",
			},
		},
	}
	for _, tt := range tests {
		var ms []map[string]interface{}
		for _, p := range tt.fields.personalizations {
			if p != nil {
				ms = append(ms, map[string]interface{}{"to": p.To(), "dynamic_template_data": p.DynamicTemplateData()})
			}
		}
		tt.want["personalizations"] = ms
		t.Run(tt.name, func(t *testing.T) {
			r := &requestJson{
				personalizations: tt.fields.personalizations,
				from:             tt.fields.from,
				subject:          tt.fields.subject,
				content:          tt.fields.content,
				templateID:       tt.fields.templateID,
				categories:       tt.fields.categories,
			}
			if got := r.makeMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestJson.makeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_requestJson_AddPersonalizations personalizationが999個(1,000個目の追加)の時と1,000個(1,001個目の追加)の時の境界値テスト
func Test_requestJson_AddPersonalizations(t *testing.T) {
	type fields struct {
		personalizations [1000]*personalization.Personalization
		from             map[string]interface{}
		subject          string
		content          []map[string]interface{}
		templateID       string
		categories       []string
	}
	type args struct {
		p *personalization.Personalization
	}
	var p1 [1000]*personalization.Personalization
	for i := 0; i < len(p1)-1; i++ {
		p1[i] = personalization.New([1]map[string]interface{}{{"email": "to@gmail.com"}}, map[string]interface{}{})
	}
	var p2 [1000]*personalization.Personalization
	for i := 0; i < len(p2); i++ {
		p2[i] = personalization.New([1]map[string]interface{}{{"email": "to@gmail.com"}}, map[string]interface{}{})
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"999",
			fields{personalizations: p1},
			args{personalization.New([1]map[string]interface{}{{"email": "to@gmail.com"}}, map[string]interface{}{})},
			false,
		},
		{
			"1000",
			fields{personalizations: p2},
			args{personalization.New([1]map[string]interface{}{{"email": "to@gmail.com"}}, map[string]interface{}{})},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestJson{
				personalizations: tt.fields.personalizations,
				from:             tt.fields.from,
				subject:          tt.fields.subject,
				content:          tt.fields.content,
				templateID:       tt.fields.templateID,
				categories:       tt.fields.categories,
			}
			if err := r.AddPersonalizations(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("requestJson.AddPersonalizations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
