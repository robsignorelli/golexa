package speech

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/language"
)

// NewTemplate creates a brand new template for one of your possible responses. It includes an
// English (US) translation by default, but you can use the `WithTranslation` option to define
// translations for other languages. Since the translations use standard Go templates, you can
// define custom functions that should be available in this template using `WithFunc`.
func NewTemplate(englishSpeech string, options ...TemplateOption) Template {
	t := Template{
		funcMap:      template.FuncMap{},
		translations: map[language.Tag]*template.Template{},
	}

	// Make sure that the funcMap is populated before attempting to do the translations
	// otherwise they'll likely fail if the function hasn't been added to the template yet.
	// This allows you to have any oder of WithFunc() and WithTranslation() that you want.
	sort.Slice(options, func(i, j int) bool {
		return options[i].order < options[j].order
	})
	for _, opt := range options {
		opt.apply(&t)
	}

	// Your required default translation trumps any attempts to WithTranslation it away later.
	WithTranslation(language.AmericanEnglish, englishSpeech).apply(&t)
	return t
}

// Template represents a single type of speech response that you can return from your intent handler.
// The idea is that if your intent can say one of 5 possible "things" based on the input data then
// you should create 5 of these so that you can decide which one to evaluate. A golexa template also
// supports defining translations for other languages based on the incoming request's language.
type Template struct {
	funcMap      template.FuncMap
	translations map[language.Tag]*template.Template
}

// Eval processes the parsed template using the given contextual data.
func (t Template) Eval(ctx TemplateContext) (string, error) {
	localizedTemplate := t.translationFor(ctx.Language)
	output := strings.Builder{}

	if err := localizedTemplate.Execute(&output, ctx); err != nil {
		return "", fmt.Errorf("template eval: %v", err)
	}
	return strings.TrimSpace(output.String()), nil
}

func (t Template) translationFor(lang language.Tag) *template.Template {
	if localizedTemplate := t.translations[lang]; localizedTemplate != nil {
		return localizedTemplate
	}
	// We've gone from "es-MX" to "es-241" to "es" and still no translation. We're now
	// at "und", so fall back to the English translation.
	if lang.IsRoot() {
		return t.translations[language.AmericanEnglish]
	}
	return t.translationFor(lang.Parent())
}

// WithFunc adds a named function that will be available when parsing/evaluating responses.
func WithFunc(name string, function interface{}) TemplateOption {
	return TemplateOption{
		order: 0,
		apply: func(t *Template) {
			t.funcMap[name] = function
		},
	}
}

// WithTranslation defines a version of this response that is localized for the given language. This
// can be a plain text template (parsed using standard Go templates) or an SSML template.
func WithTranslation(lang language.Tag, localizedSpeech string) TemplateOption {
	return TemplateOption{
		order: 1,
		apply: func(t *Template) {
			t.translations[lang] = template.Must(template.New(lang.String()).
				Funcs(t.funcMap).
				Parse(localizedSpeech))
		},
	}
}

// TemplateOption should not be used directly. Use WithFunc or WithTranslation to provide the
// template option of your choice.
type TemplateOption struct {
	apply func(*Template)
	order int
}

// TemplateContext represents the single piece of 'data' that you pass to a template
// when evaluating it. This gives your template access to some higher level data like
// the current timestamp and request as well as any data you generated when handling
// the intent that has an effect on the speech response.
type TemplateContext struct {
	Language language.Tag
	Now      time.Time
	Value    interface{}
}
