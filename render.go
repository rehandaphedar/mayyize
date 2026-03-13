package main

import (
	"log"
	"sort"

	qul "git.sr.ht/~rehandaphedar/genanki-go-utils/v2/pkg/qul"
)

func renderInstances(wordIndex qul.WordIndex, metadataAyahByVerseKey map[string]qul.MetadataAyah, phrase qul.Phrase) []Instance {
	var instances []Instance

	for verseKey := range phrase.Ayah {
		ranges := phrase.Ayah[verseKey]
		for instanceInVerseIndex, instanceInVerse := range ranges {
			from := instanceInVerse[0]
			to := instanceInVerse[1]
			instanceInVerseNumber := instanceInVerseIndex + 1

			instances = append(instances, Instance{
				VerseKey:        verseKey,
				InstanceInVerse: instanceInVerseNumber,
				Phrase:          renderPhrase(wordIndex, verseKey, from, to),
				Context:         renderContext(wordIndex, metadataAyahByVerseKey, verseKey, from),
				Continuation:    renderContinuation(wordIndex, metadataAyahByVerseKey, verseKey, to),
			})
		}
	}
	sort.Slice(instances, func(i, j int) bool {
		return compareInstances(instances[i], instances[j])
	})

	return instances
}

func renderPhrase(wordIndex qul.WordIndex, verseKey string, from, to int) []string {
	return renderRange(wordIndex, qul.Source{Key: verseKey, From: from, To: to})
}

func renderContext(wordIndex qul.WordIndex, metadataAyahByVersekey map[string]qul.MetadataAyah, verseKey string, from int) []string {
	if from == 1 {
		previousVerseKey, found := qul.GetPreviousVerseKey(metadataAyahByVersekey, verseKey)
		if found {
			return renderVerseFrom(wordIndex, metadataAyahByVersekey, previousVerseKey, 1)
		}
		// TODO: Better Context, Continuation edge cases?
		return []string{}
	}
	return renderRange(wordIndex, qul.Source{Key: verseKey, From: 1, To: from - 1})
}

func renderContinuation(wordIndex qul.WordIndex, metadataAyahByVerseKey map[string]qul.MetadataAyah, verseKey string, to int) []string {
	words := wordIndex.VerseWords[verseKey]

	if to+1 == len(words) {
		nextVerseKey, found := qul.GetNextVerseKey(metadataAyahByVerseKey, verseKey)
		if found {
			return renderVerseFrom(wordIndex, metadataAyahByVerseKey, nextVerseKey, 1)
		}
		return []string{}
	}
	return renderVerseFrom(wordIndex, metadataAyahByVerseKey, verseKey, to+1)
}

func renderVerseFrom(wordIndex qul.WordIndex, metadataAyahByVerseKey map[string]qul.MetadataAyah, verseKey string, from int) []string {
	return renderRange(wordIndex, qul.Source{Key: verseKey, From: from, To: metadataAyahByVerseKey[verseKey].WordsCount})
}

func renderRange(wordIndex qul.WordIndex, source qul.Source) []string {
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
	return parts
}
