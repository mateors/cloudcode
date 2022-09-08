package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

var workingDirectory string
var staticResourceRelativePath string

//go:embed assets/*
var assetsDir embed.FS

var localizer *i18n.Localizer
var bundle *i18n.Bundle

func init() {

	workingDirectory, _ = os.Getwd()

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	//bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//bundle.MustLoadMessageFile("active.es.toml")
	bundle.LoadMessageFile("locals/en.json")
	bundle.LoadMessageFile("locals/bn.json")

	//fmt.Println(Localize("en", "LocRun"))
	//os.Exit(1)
}

func main() {

	//staticMethod()
	//loadFileMethod()
	//dynamicFileMethod("en")
	//dynamicFileMethod("bn")
	//dynamicFileMethod(language.Bengali.String())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	//r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	assetPath := filepath.Join(workingDirectory, "assets")
	//http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(assetPath))))//regular http Handle
	r.Handle("/resources/*", http.StripPrefix("/resources/", http.FileServer(http.Dir(assetPath))))
	staticResourceRelativePath = "/resources" //resources/css/resizer.css

	//r.Handle("/resources/*", http.StripPrefix("/resources/", http.FileServer(http.FS(assetsDir))))
	//staticResourceRelativePath = "/resources/assets"

	r.HandleFunc("/", indexHandler)
	//http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8081", r)
}

func GetBaseURL(r *http.Request) string {

	var baseurl, proto string
	if strings.Contains(r.Proto, "HTTP") {
		proto = "http"
	} else {
		proto = "https"
	}
	baseurl = fmt.Sprintf("%s://%s", proto, r.Host)
	return baseurl
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	var baseurl string = GetBaseURL(r)

	lang := r.FormValue("lang")
	// accept := r.Header.Get("Accept-Language")
	// localizer = i18n.NewLocalizer(bundle, lang, accept)
	//localizer := i18n.NewLocalizer(bundle, languageTag) //8

	// localizeConfig := i18n.LocalizeConfig{
	// 	MessageID: "messages",
	// 	TemplateData: map[string]interface{}{
	// 		"Name":  "Mostain",
	// 		"Count": 0,
	// 	},
	// }

	// message, err := localizer.Localize(&localizeConfig)
	// if err != nil {
	// 	log.Println(err)
	// }

	tdata := make(map[string]interface{})
	tdata["Name"] = "Mostain"
	tdata["Count"] = 0
	message := LocalizeTemplate(lang, "messages", tdata)

	// LocRun, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "LocRun"})
	// if err != nil {
	// 	log.Println(err)
	// }

	LocRun := Localize(lang, "LocRun")

	// txtInvite, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "btnInvite"})
	// if err != nil {
	// 	log.Println(err)
	// }

	LocInvite := Localize(lang, "LocInvite")

	tmplt, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		log.Println(err)
		return
	}

	//fmt.Println("r.RemoteAddr", r.RemoteAddr)
	// local := map[string]interface{}{
	// 	"Title": "name",
	// 	"Paragraphs": []string{
	// 		message,
	// 		message,
	// 	},
	// }

	// data := struct {
	// 	Base   string
	// 	Title  string
	// 	Static string
	// 	Local  map[string]interface{}
	// }{
	// 	Base:   baseurl,
	// 	Title:  "CC",
	// 	Static: staticResourceRelativePath,
	// 	Local:  local,
	// }

	data := map[string]interface{}{
		"Title": "CC",
		"Paragraphs": []string{
			message,
			message,
		},
		"Base":      baseurl,
		"Static":    staticResourceRelativePath,
		"LocInvite": LocInvite,
		"LocRun":    LocRun,
	}

	err = tmplt.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func Localize(lang, messageID string) string {

	//*i18n.Localizer{}
	localizer = i18n.NewLocalizer(bundle, lang)
	transalation, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		log.Println(err)
	}
	return transalation
}

func LocalizeTemplate(lang, messageID string, data map[string]interface{}) string {

	localizer = i18n.NewLocalizer(bundle, lang)
	transalation, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID, TemplateData: data})
	if err != nil {
		log.Println(err)
	}
	return transalation
}
