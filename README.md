## sendgrid_api_helper is help you to sending email with sendgrid API v3(json data)


Just set personalization (destination and variables) and source, then either subject and content, or template ID.

The sample code is below.

``` sample.go
    j := json.New(json.From("from@gmail.com"))
	j.SetCategories("category")

	p := personalization.New(
		personalization.To("aaa@gmail.com"),
		map[string]interface{}{
			"lastname":  "toyohashi",
			"firstname": "6140",
		},
	)
	err := j.AddPersonalizations(p)

	// エラー処理: リクエストを送り、personalizationをリセットし、追加できなかった要素を登録し直す
	if err != nil {
		err = j.RequestPost(os.Getenv("SENDGRID_API_KEY"))
		if err != nil {
			fmt.Println(err.Error()) // リクエスト失敗
		}
		j.ResetPersonalization()
		j.AddPersonalizations(p)
	}
```
Up to 10 categories can be set in Set Categories.
Specifying the category is optional.
