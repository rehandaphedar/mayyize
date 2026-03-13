package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"git.sr.ht/~rehandaphedar/genanki-go-utils/v2/pkg/qul"
	"github.com/npcnixel/genanki-go"
	"go.yaml.in/yaml/v4"
)

func main() {
	modelIdPhrase := flag.Int64("model-id-phrase", int64(1805162761), "ID of the phrase model")
	modelNamePhrase := flag.String("model-name-phrase", "mayyize - Phrase", "Name of the phrase model")

	deckIdPhrase := flag.Int64("deck-id-phrase", int64(1748329869), "ID of the phrase deck")
	deckNamePhrase := flag.String("deck-name-phrase", "mayyize::Phrase Recognition", "Name of the phrase peck")
	deckDescriptionPhrase := flag.String("deck-description-phrase", "Recall all instances of the phrase.", "Description of the phrase deck")

	modelIdVerse := flag.Int64("model-id-verse", int64(1357701653), "ID of the verse model")
	modelNameVerse := flag.String("model-name-verse", "mayyize - Verse", "Name of the verse model")

	deckIdVerse := flag.Int64("deck-id-verse", int64(1850992899), "ID of the verse deck")
	deckNameVerse := flag.String("deck-name-verse", "mayyize::Verse Completion", "Name of the verse deck")
	deckDescriptionVerse := flag.String("deck-description-verse", "Recall the correct instance of the phrase to complete the verse.", "Description of the Verse Deck")

	outputPath := flag.String("output", "out/mayyize.apkg", "Output filepath")

	templateHtmlPath := flag.String("template-html", "templates/index.gohtml", "Path to template file")
	templateCssPath := flag.String("template-css", "templates/style.css", "Path to CSS file")

	templatePhraseFrontName := flag.String("template-phrase-front", "phrase_front", "Name of the phrase front template")
	templatePhraseBackName := flag.String("template-phrase-back", "phrase_back", "Name of the phrase back template")
	templateVerseFrontName := flag.String("template-verse-front", "verse_front", "Name of the verse front template")
	templateVerseBackName := flag.String("template-verse-back", "verse_back", "Name of the verse back template")

	wordsPath := flag.String("words", "data/qpc-hafs-word-by-word.json", "Path to words data")
	phrasesPath := flag.String("phrases", "data/phrases.json", "Path to phrases data")
	layoutPath := flag.String("layout", "data/qpc-v4-tajweed-15-lines.db", "Path to layout data")
	metadataAyahPath := flag.String("metadata-ayah", "data/quran-metadata-ayah.json", "Path to ayah metadata")
	metadataJuzPath := flag.String("metadata-juz", "data/quran-metadata-juz.json", "Path to juz metadata")
	metadataHizbPath := flag.String("metadata-hizb", "data/quran-metadata-hizb.json", "Path to hizb metadata")
	metadataRubPath := flag.String("metadata-rub", "data/quran-metadata-rub.json", "Path to rub metadata")
	metadataManzilPath := flag.String("metadata-manzil", "data/quran-metadata-manzil.json", "Path to manzil metadata")
	metadataRukuPath := flag.String("metadata-ruku", "data/quran-metadata-ruku.json", "Path to ruku metadata")

	var tagFormat qul.TagFormat

	tagFormat.Chapter = flag.String("tag-format-chapter", "quran::chapter::%03d", "Format of the chapter tag. %d is replaced with the chapter number.")
	tagFormat.Verse = flag.String("tag-format-verse", "quran::verse::%s", "Format of the verse tag. %s is replaced with the zero padded verse key (Example: 001:001).")
	tagFormat.Page = flag.String("tag-format-page", "quran::page::%03d", "Format of the page tag. %d is replaced with the page number.")
	tagFormat.Juz = flag.String("tag-format-juz", "quran::juz::%02d", "Format of the juz tag. %d is replaced with the juz number.")
	tagFormat.Hizb = flag.String("tag-format-hizb", "quran::hizb::%02d", "Format of the hizb tag. %d is replaced with the hizb number.")
	tagFormat.Rub = flag.String("tag-format-rub", "quran::rub::%03d", "Format of the rub tag. %d is replaced with the rub number.")
	tagFormat.Manzil = flag.String("tag-format-manzil", "quran::manzil::%d", "Format of the manzil tag. %d is replaced with the manzil number.")
	tagFormat.Ruku = flag.String("tag-format-ruku", "quran::ruku::%03d", "Format of the ruku tag. %d is replaced with the ruku number.")

	mediaConfigPath := flag.String("media-config", "media/config.yaml", "Path to media config")

	flag.Parse()

	var words map[string]qul.Word
	var phrases map[string]qul.Phrase
	var metadataAyah map[string]qul.MetadataAyah

	var metadataDivision qul.MetadataDivision

	err := loadJSON(*wordsPath, &words)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*phrasesPath, &phrases)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataAyahPath, &metadataAyah)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataJuzPath, &metadataDivision.Juz)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataHizbPath, &metadataDivision.Hizb)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataRubPath, &metadataDivision.Rub)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataManzilPath, &metadataDivision.Manzil)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataRukuPath, &metadataDivision.Ruku)
	if err != nil {
		log.Fatal(err)
	}

	index, err := qul.BuildIndex(*layoutPath, words, metadataDivision, tagFormat)
	if err != nil {
		log.Fatalf("build index: %v", err)
	}

	metadataAyahByVerseKey := make(map[string]qul.MetadataAyah)
	for _, metadataAyahEntry := range metadataAyah {
		metadataAyahByVerseKey[metadataAyahEntry.VerseKey] = metadataAyahEntry
	}

	css, err := readFile(*templateCssPath)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(*templateHtmlPath)
	if err != nil {
		log.Fatalf("parse template files: %v", err)
	}
	var buf bytes.Buffer

	modelPhrase := genanki.NewModel(*modelIdPhrase, *modelNamePhrase).
		SetCSS(css).
		AddField(genanki.Field{Name: "PhraseID"}).
		AddField(genanki.Field{Name: "Count"}).
		AddField(genanki.Field{Name: "Chapters"}).
		AddField(genanki.Field{Name: "Front"}).
		AddField(genanki.Field{Name: "Back"}).
		AddTemplate(genanki.Template{
			Name: "Phrase Recognition",
			Qfmt: "{{Front}}",
			Afmt: "{{Back}}",
		})
	deckPhrase := genanki.NewDeck(*deckIdPhrase, *deckNamePhrase, *deckDescriptionPhrase)

	modelVerse := genanki.NewModel(*modelIdVerse, *modelNameVerse).
		SetCSS(css).
		AddField(genanki.Field{Name: "PhraseInstanceID"}).
		AddField(genanki.Field{Name: "VerseKey"}).
		AddField(genanki.Field{Name: "PhraseID"}).
		AddField(genanki.Field{Name: "InstanceInVerse"}).
		AddField(genanki.Field{Name: "Count"}).
		AddField(genanki.Field{Name: "Chapters"}).
		AddField(genanki.Field{Name: "Front"}).
		AddField(genanki.Field{Name: "Back"}).
		AddTemplate(genanki.Template{
			Name: "Verse Completion",
			Qfmt: "{{Front}}",
			Afmt: "{{Back}}",
		})
	deckVerse := genanki.NewDeck(*deckIdVerse, *deckNameVerse, *deckDescriptionVerse)

	for phraseId, phrase := range phrases {
		instances := renderInstances(index.Word, metadataAyahByVerseKey, phrase)

		templateDataPhrase := TemplateDataPhrase{
			Count:     phrase.Count,
			Chapters:  phrase.Surahs,
			Phrase:    renderRange(index.Word, phrase.Source),
			Instances: instances,
		}

		templateErrorMessage := "error while executing template %s with data %+v: %v"

		err := tmpl.ExecuteTemplate(&buf, *templatePhraseFrontName, templateDataPhrase)
		if err != nil {
			log.Printf(templateErrorMessage, *templatePhraseFrontName, templateDataPhrase, err)
		}
		phraseFront := buf.String()
		buf.Reset()

		err = tmpl.ExecuteTemplate(&buf, *templatePhraseBackName, templateDataPhrase)
		if err != nil {
			log.Printf(templateErrorMessage, *templatePhraseBackName, templateDataPhrase, err)
		}
		phraseBack := buf.String()
		buf.Reset()

		notePhrase := genanki.NewNote(
			modelPhrase.ID,
			[]string{
				phraseId,
				strconv.Itoa(templateDataPhrase.Count),
				strconv.Itoa(templateDataPhrase.Chapters),
				phraseFront,
				phraseBack,
			},
			qul.BuildTagsForPhrase(index, phrase),
		)

		noteIdBasePhrase := fmt.Sprintf("%d_%s", modelPhrase.ID, phraseId)
		notePhrase.ID = qul.GenerateID(noteIdBasePhrase)
		deckPhrase.AddNote(notePhrase)

		for verseKey := range phrase.Ayah {
			instancesInVerse := phrase.Ayah[verseKey]

			for instanceInVerseIndex, instanceInVerse := range instancesInVerse {
				from := instanceInVerse[0]
				to := instanceInVerse[1]

				paddedVerseKey, err := qul.PadVerseKey(verseKey)
				if err != nil {
					log.Printf("error while padding verse key %s: %v", verseKey, err)
					continue
				}
				instanceInVerseNumber := instanceInVerseIndex + 1
				phraseInstanceId := fmt.Sprintf("%s_%s_%02d", paddedVerseKey, phraseId, instanceInVerseNumber)

				templateDataVerse := TemplateDataVerse{
					VerseKey:        verseKey,
					InstanceInVerse: instanceInVerseNumber,
					Count:           phrase.Count,
					Chapters:        phrase.Surahs,
					Phrase:          renderRange(index.Word, qul.Source{Key: verseKey, From: from, To: to}),
					Context:         renderContext(index.Word, metadataAyahByVerseKey, verseKey, from),
					Continuation:    renderContinuation(index.Word, metadataAyahByVerseKey, verseKey, to),
					Instances:       instances,
				}
				err = tmpl.ExecuteTemplate(&buf, *templateVerseFrontName, templateDataVerse)
				if err != nil {
					log.Printf(templateErrorMessage, *templateVerseFrontName, templateDataVerse, err)
				}
				verseFront := buf.String()
				buf.Reset()

				err = tmpl.ExecuteTemplate(&buf, *templateVerseBackName, templateDataVerse)
				if err != nil {
					log.Printf(templateErrorMessage, *templateVerseBackName, templateDataVerse, err)
				}
				verseBack := buf.String()
				buf.Reset()

				noteVerse := genanki.NewNote(
					modelVerse.ID,
					[]string{
						phraseInstanceId,
						verseKey,
						phraseId,
						strconv.Itoa(instanceInVerseNumber),
						strconv.Itoa(phrase.Count),
						strconv.Itoa(phrase.Surahs),
						verseFront,
						verseBack,
					},
					index.Tag.Verse[verseKey],
				)

				noteIdBaseVerse := fmt.Sprintf("%d_%s", modelVerse.ID, phraseInstanceId)
				noteVerse.ID = qul.GenerateID(noteIdBaseVerse)
				deckVerse.AddNote(noteVerse)
			}
		}
	}

	pkg := genanki.NewPackage([]*genanki.Deck{deckPhrase, deckVerse}).AddModel(modelPhrase).AddModel(modelVerse)

	if *mediaConfigPath != "" {
		mediaConfigDir := filepath.Dir(*mediaConfigPath)

		mediaConfigData, err := os.ReadFile(*mediaConfigPath)
		if err != nil {
			log.Fatalf("read media config: %v", err)
		}

		var mediaEntries []MediaEntry
		if err := yaml.Unmarshal(mediaConfigData, &mediaEntries); err != nil {
			log.Fatalf("parse media config: %v", err)
		}

		for _, mediaEntry := range mediaEntries {
			src := filepath.Join(mediaConfigDir, mediaEntry.Src)
			as := mediaEntry.As
			mediaEntryData, err := os.ReadFile(src)
			if err != nil {
				log.Fatalf("read media entry %s: %v", src, err)
			}
			pkg.AddMedia(as, mediaEntryData)
		}
	}

	if err := pkg.WriteToFile(*outputPath); err != nil {
		log.Fatalf("write package to %s: %v", *outputPath, err)
	}
}
