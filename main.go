package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func staticMethod() {

	// Step 1: Create bundle
	bundle := i18n.NewBundle(language.English)

	//Step 2: Define messages
	messageEn := &i18n.Message{
		ID:          "Email",
		Description: "The number of unread emails a user has",
		One:         "{{.Name}} has {{.Count}} email.",
		Other:       "{{.Name}} has {{.Count}} emails.",
	}

	messageBn := &i18n.Message{
		ID:          "Email",
		Description: "ব্যবহারকারীর অপঠিত ইমেলের সংখ্যা",
		One:         "{{.Name}} {{.Count}} টি ইমেইল",
		Other:       "{{.Name}} {{.Count}} টি ইমেইল",
	}

	//Step-3: Add messages
	bundle.AddMessages(language.English, messageEn)
	bundle.AddMessages(language.Bengali, messageBn)

	//custom localizer
	localizer := i18n.NewLocalizer(bundle, language.English.String()) //set your desired language tag

	messagesCount := 2
	translation, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "Email",
		//DefaultMessage: messageBn,
		TemplateData: map[string]interface{}{
			"Name":  "Mostain",
			"Count": messagesCount,
		},
		PluralCount: messagesCount,
	})

	fmt.Println(err, translation)

}

func loadFileMethod() {

	// Step 1: Create bundle
	bundle := i18n.NewBundle(language.Bengali)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	//bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	bundle.LoadMessageFile("locals/en.json")
	bundle.LoadMessageFile("locals/bn.json")

	localizer := i18n.NewLocalizer(bundle, language.Bengali.String()) //8

	messagesCount := 2

	localizeConfig := i18n.LocalizeConfig{ //5
		MessageID: "messages",
		TemplateData: map[string]interface{}{
			"Name":  "Mostain",
			"Count": messagesCount,
		},
	}

	localization, err := localizer.Localize(&localizeConfig) //set language which you want to translate
	fmt.Println(err, localization)
}

func dynamicFileMethod(languageTag string) {

	// Step 1: Create bundle
	bundle := i18n.NewBundle(language.Bengali)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	//bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	bundle.LoadMessageFile("locals/en.json")
	bundle.LoadMessageFile("locals/bn.json")

	localizer := i18n.NewLocalizer(bundle, languageTag) //8

	messagesCount := 2

	localizeConfig := i18n.LocalizeConfig{ //5
		MessageID: "messages",
		TemplateData: map[string]interface{}{
			"Name":  "Mostain",
			"Count": messagesCount,
		},
	}

	localization, err := localizer.Localize(&localizeConfig) //set language which you want to translate
	fmt.Println(err, localization)
}

var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
  <title>Golang i18n</title>
 <body>
	<h1>{{.Title}}</h1>
	{{range .Paragraphs}}<p>{{.}}</p>{{end}}
 </body>
</html>`))

func main() {

	//staticMethod()
	//loadFileMethod()
	//dynamicFileMethod("en")
	//dynamicFileMethod("bn")
	//dynamicFileMethod(language.Bengali.String())

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8081", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	//bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	//bundle.MustLoadMessageFile("active.es.toml")
	bundle.LoadMessageFile("locals/en.json")
	bundle.LoadMessageFile("locals/bn.json")

	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")

	localizer := i18n.NewLocalizer(bundle, lang, accept)
	//localizer := i18n.NewLocalizer(bundle, languageTag) //8

	name := "Wania"
	messagesCount := 2

	localizeConfig := i18n.LocalizeConfig{ //5
		MessageID: "messages",
		TemplateData: map[string]interface{}{
			"Name":  "Mostain",
			"Count": messagesCount,
		},
	}

	message, err := localizer.Localize(&localizeConfig)
	if err != nil {
		log.Println(err)
	}

	data := map[string]interface{}{
		"Title": name,
		"Paragraphs": []string{
			message,
			message,
		},
	}

	err = page.Execute(w, data)
	fmt.Println(err)
}
