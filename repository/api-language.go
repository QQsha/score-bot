package repository

import (
	"fmt"
	"strconv"

	"github.com/chrisport/go-lang-detector/langdet"
	"github.com/chrisport/go-lang-detector/langdet/langdetdef"
	"github.com/jdkato/prose/v2"
)

type DetectLanguageAPI struct {
	detector langdet.Detector
}

func NewLanguageAPI() *DetectLanguageAPI {
	detector := langdetdef.NewWithDefaultLanguages()
	return &DetectLanguageAPI{
		detector: detector,
	}
}

func (d DetectLanguageAPI) EnglisheDetector(text string) (bool, error) {
	if len(text) < 10 {
		return false, nil
	}
	var notEnglish bool
	fullResults := d.detector.GetLanguages(text)
	if len(fullResults) > 0 && fullResults[0].Name != "english" && fullResults[0].Confidence > 50 {
		notEnglish = true
	}
	return notEnglish, nil
}

func (d DetectLanguageAPI) EnglishDetectorTest(text string) string {
	fullResults := d.detector.GetLanguages(text)
	res := "langdetdef: I guess this is " + fullResults[0].Name + " for " + strconv.Itoa(fullResults[0].Confidence) + "\n"
	return res
}

func (d DetectLanguageAPI) NameDetectorTest(text string) string {
	doc, _ := prose.NewDocument(text)
	res := ""
	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
		res += ent.Text + " is " + ent.Label + "\n"
		// Lebron James PERSON
		// Los Angeles GPE
	}
	return res
}
func (d DetectLanguageAPI) NameDetector(text string) bool {
	doc, _ := prose.NewDocument(text)
	for _, ent := range doc.Entities() {
		if ent.Label == "GPE" || ent.Label == "PERSON" {
			return true
		}
	}
	return false
}
