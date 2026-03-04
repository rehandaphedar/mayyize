package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	qul "git.sr.ht/~rehandaphedar/genanki-go-utils/pkg/qul"
	genanki "github.com/npcnixel/genanki-go"
)

func main() {
	modelIdPhrase := flag.Int64("model-id-phrase", int64(1805162761), "ID of the Phrase Model")
	modelNamePhrase := flag.String("model-name-phrase", "Mutashābihāt - Phrase", "Name of the Phrase Model")

	deckIdPhrase := flag.Int64("deck-id-phrase", int64(1748329869), "ID of the Phrase Deck")
	deckNamePhrase := flag.String("deck-name-phrase", "Mutashābihāt::Phrase Recognition", "Name of the Phrase Deck")
	deckDescriptionPhrase := flag.String("deck-description-phrase", "Recall all instances of the phrase.", "Description of the Phrase Deck")

	modelIdVerse := flag.Int64("model-id-verse", int64(1357701653), "ID of the Verse Model")
	modelNameVerse := flag.String("model-name-verse", "Mutashābihāt - Verse", "Name of the Verse Model")

	deckIdVerse := flag.Int64("deck-id-verse", int64(1850992899), "ID of the Verse Deck")
	deckNameVerse := flag.String("deck-name-verse", "Mutashābihāt::Verse Completion", "Name of the Verse Deck")
	deckDescriptionVerse := flag.String("deck-description-verse", "Recall the correct instance of the phrase to complete the verse.", "Description of the Verse Deck")

	outputPath := flag.String("output", "out/quran_mutashabihat.apkg", "Output filepath")

	cssPath := flag.String("css", "templates/style.css", "Path to CSS file")
	QfmtPhrasePath := flag.String("qfmt-phrase", "templates/qfmt_phrase.html", "Path to Phrase Qfmt HTML Template")
	AfmtPhrasePath := flag.String("afmt-phrase", "templates/afmt_phrase.html", "Path to Phrase Afmt HTML Template")
	QfmtVersePath := flag.String("qfmt-verse", "templates/qfmt_verse.html", "Path to Verse Qfmt HTML Template")
	AfmtVersePath := flag.String("afmt-verse", "templates/afmt_verse.html", "Path to Verse Afmt HTML Template")

	wordsPath := flag.String("words", "metadata/qpc-hafs-word-by-word.json", "Path to words JSON")
	phrasesPath := flag.String("phrases", "metadata/phrases.json", "Path to phrases JSON")
	layoutPath := flag.String("layout", "metadata/qpc-v4-tajweed-15-lines.db", "Path to Mushaf DB")
	metadataAyahPath := flag.String("metadata-ayah", "metadata/quran-metadata-ayah.json", "Path to Ayah Metadata")
	metadataJuzPath := flag.String("metadata-juz", "metadata/quran-metadata-juz.json", "Path to Juz Metadata")
	metadataHizbPath := flag.String("metadata-hizb", "metadata/quran-metadata-hizb.json", "Path to Hizb Metadata")
	metadataRubPath := flag.String("metadata-rub", "metadata/quran-metadata-rub.json", "Path to Rub Metadata")
	metadataManzilPath := flag.String("metadata-manzil", "metadata/quran-metadata-manzil.json", "Path to Manzil Metadata")
	metadataRukuPath := flag.String("metadata-ruku", "metadata/quran-metadata-ruku.json", "Path to Ruku Metadata")

	var tagFormat qul.TagFormat

	tagFormat.Chapter = flag.String("tag-format-chapter", "quran::chapter::%03d", "Format of the chapter tag. %d is replaced with the chapter number.")
	tagFormat.Verse = flag.String("tag-format-verse", "quran::verse::%s", "Format of the verse tag. %s is replaced with the zero padded verse key (Example: 001:001).")
	tagFormat.Page = flag.String("tag-format-page", "quran::page::%03d", "Format of the page tag. %d is replaced with the page number.")
	tagFormat.Juz = flag.String("tag-format-juz", "quran::juz::%02d", "Format of the juz tag. %d is replaced with the juz number.")
	tagFormat.Hizb = flag.String("tag-format-hizb", "quran::hizb::%02d", "Format of the hizb tag. %d is replaced with the hizb number.")
	tagFormat.Rub = flag.String("tag-format-rub", "quran::rub::%03d", "Format of the rub tag. %d is replaced with the rub number.")
	tagFormat.Manzil = flag.String("tag-format-manzil", "quran::manzil::%d", "Format of the manzil tag. %d is replaced with the manzil number.")
	tagFormat.Ruku = flag.String("tag-format-ruku", "quran::ruku::%03d", "Format of the ruku tag. %d is replaced with the ruku number.")

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

	css, err := readFile(*cssPath)
	if err != nil {
		log.Fatal(err)
	}
	qfmtPhrase, err := readFile(*QfmtPhrasePath)
	if err != nil {
		log.Fatal(err)
	}
	afmtPhrase, err := readFile(*AfmtPhrasePath)
	if err != nil {
		log.Fatal(err)
	}
	qfmtVerse, err := readFile(*QfmtVersePath)
	if err != nil {
		log.Fatal(err)
	}
	afmtVerse, err := readFile(*AfmtVersePath)
	if err != nil {
		log.Fatal(err)
	}

	modelPhrase := genanki.NewModel(*modelIdPhrase, *modelNamePhrase).
		SetCSS(css).
		AddField(genanki.Field{Name: "PhraseID"}).
		AddField(genanki.Field{Name: "Chapters"}).
		AddField(genanki.Field{Name: "Count"}).
		AddField(genanki.Field{Name: "Phrase"}).
		AddField(genanki.Field{Name: "AllInstances"}).
		AddTemplate(genanki.Template{
			Name: "Phrase Recognition",
			Qfmt: qfmtPhrase,
			Afmt: afmtPhrase,
		})
	deckPhrase := genanki.NewDeck(*deckIdPhrase, *deckNamePhrase, *deckDescriptionPhrase)

	modelVerse := genanki.NewModel(*modelIdVerse, *modelNameVerse).
		SetCSS(css).
		AddField(genanki.Field{Name: "NoteID"}).
		AddField(genanki.Field{Name: "VerseKey"}).
		AddField(genanki.Field{Name: "PhraseID"}).
		AddField(genanki.Field{Name: "Chapters"}).
		AddField(genanki.Field{Name: "Count"}).
		AddField(genanki.Field{Name: "Phrase"}).
		AddField(genanki.Field{Name: "Context"}).
		AddField(genanki.Field{Name: "Continuation"}).
		AddField(genanki.Field{Name: "AllInstances"}).
		AddTemplate(genanki.Template{
			Name: "Verse Completion",
			Qfmt: qfmtVerse,
			Afmt: afmtVerse,
		})
	deckVerse := genanki.NewDeck(*deckIdVerse, *deckNameVerse, *deckDescriptionVerse)

	for phraseId, phrase := range phrases {
		allInstances := renderAllInstances(index.Word, metadataAyahByVerseKey, phrase)

		notePhrase := genanki.NewNote(
			modelPhrase.ID,
			[]string{
				phraseId,
				strconv.Itoa(phrase.Surahs),
				strconv.Itoa(phrase.Count),
				renderRange(index.Word, phrase.Source, "phrase"),
				allInstances,
			},
			qul.BuildTagsForPhrase(index, phrase),
		)
		deckPhrase.AddNote(notePhrase)

		for verseKey := range phrase.Ayah {
			instances := phrase.Ayah[verseKey]

			for instanceIndex, instance := range instances {
				from := instance[0]
				to := instance[1]

				paddedVerseKey, err := qul.PadVerseKey(verseKey)
				if err != nil {
					log.Printf("error while padding verse key %s: %v", verseKey, err)
					continue
				}
				instanceNumber := instanceIndex + 1
				noteId := fmt.Sprintf("%s_%s_%02d", paddedVerseKey, phraseId, instanceNumber)

				ayahNote := genanki.NewNote(
					modelVerse.ID,
					[]string{
						noteId,
						verseKey,
						phraseId,
						strconv.Itoa(phrase.Surahs),
						strconv.Itoa(phrase.Count),
						renderRange(index.Word, qul.PhraseSource{Key: verseKey, From: from, To: to}, "phrase"),
						renderContext(index.Word, metadataAyahByVerseKey, verseKey, from),
						renderContinuation(index.Word, metadataAyahByVerseKey, verseKey, to),
						allInstances,
					},
					index.Tag.Verse[verseKey],
				)
				deckVerse.AddNote(ayahNote)
			}
		}
	}

	pkg := genanki.NewPackage([]*genanki.Deck{deckPhrase, deckVerse}).AddModel(modelPhrase).AddModel(modelVerse)
	if err := pkg.WriteToFile(*outputPath); err != nil {
		log.Fatalf("write package to %s: %v", *outputPath, err)
	}
}
