package main

import (
	"fmt"
	"log"
	"strings"

	qul "git.sr.ht/~rehandaphedar/genanki-go-utils/pkg/qul"
)

func renderAllInstances(wordIndex qul.WordIndex, metadataAyahByVerseKey map[string]qul.MetadataAyah, phrase qul.Phrase) string {
	var stringsBuilder strings.Builder
	stringsBuilder.WriteString(`<div class="instances">`)

	for verseKey := range phrase.Ayah {
		ranges := phrase.Ayah[verseKey]
		for instanceIndex, instance := range ranges {
			from := instance[0]
			to := instance[1]
			context := renderContext(wordIndex, metadataAyahByVerseKey, verseKey, from)
			phrase := renderPhrase(wordIndex, verseKey, from, to)
			continuation := renderContinuation(wordIndex, metadataAyahByVerseKey, verseKey, to)

			stringsBuilder.WriteString(`<div class="instance">`)
			fmt.Fprintf(&stringsBuilder, `<div class="verse-key">%s - instance %d in chapter</div>`, verseKey, instanceIndex+1)
			stringsBuilder.WriteString(`<div class="quran-text">`)
			stringsBuilder.WriteString(context)
			stringsBuilder.WriteString(phrase)
			stringsBuilder.WriteString(continuation)
			stringsBuilder.WriteString(`</div>`)
			stringsBuilder.WriteString(`</div>`)
		}
	}

	stringsBuilder.WriteString("</div>")
	return stringsBuilder.String()
}

func renderPhrase(wordIndex qul.WordIndex, verseKey string, from, to int) string {
	return renderRange(wordIndex, qul.PhraseSource{Key: verseKey, From: from, To: to}, "phrase")
}

func renderContext(wordIndex qul.WordIndex, metadataAyahByVersekey map[string]qul.MetadataAyah, verseKey string, from int) string {
	if from == 1 {
		previousVerseKey, found := qul.GetPreviousVerseKey(metadataAyahByVersekey, verseKey)
		if found {
			return renderVerseFrom(wordIndex, metadataAyahByVersekey, previousVerseKey, 1, "context")
		}
		return `<span class="opening">[The Opening of the Qurʾān]</span>`
	}
	return renderRange(wordIndex, qul.PhraseSource{Key: verseKey, From: 1, To: from - 1}, "context")
}

func renderContinuation(wordIndex qul.WordIndex, metadataAyahByVerseKey map[string]qul.MetadataAyah, verseKey string, to int) string {
	words := wordIndex.VerseWords[verseKey]

	if to+1 == len(words) {
		nextVerseKey, found := qul.GetNextVerseKey(metadataAyahByVerseKey, verseKey)
		if found {
			return renderVerseFrom(wordIndex, metadataAyahByVerseKey, nextVerseKey, 1, "continuation")
		}
		return `<span class="conclusion">[The Conclusion of the Qurʾān]</span>`
	}
	return renderVerseFrom(wordIndex, metadataAyahByVerseKey, verseKey, to+1, "continuation")
}

func renderVerseFrom(wordIndex qul.WordIndex, metadataAyahByVerseKey map[string]qul.MetadataAyah, verseKey string, from int, class string) string {
	return renderRange(wordIndex, qul.PhraseSource{Key: verseKey, From: from, To: metadataAyahByVerseKey[verseKey].WordsCount}, class)
}

func renderRange(wordIndex qul.WordIndex, source qul.PhraseSource, class string) string {
	words := wordIndex.VerseWords[source.Key]

	if source.To >= len(words) {
		log.Printf("invalid range %+v, silently fixing", source)
		source.To = len(words)
	}

	if (source.To + 1) == len(words) {
		source.To++
	}

	var parts []string
	for i := source.From - 1; i < source.To; i++ {
		parts = append(parts, words[i])
	}
	return fmt.Sprintf(`<span class="%s">%s</div>`, class, strings.Join(parts, " "))
}
