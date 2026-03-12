package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"git.sr.ht/~rehandaphedar/genanki-go-utils/pkg/qul"
)

func loadJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return nil
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(data), nil
}
func compareInstances(a, b Instance) bool {
	compareInstancesErrorMessage := "error while compare instances %+v and %+v: decode verse key %s: %v"

	chapterA, verseA, err := qul.DecodeVerseKey(a.VerseKey)
	if err != nil {
		log.Println(compareInstancesErrorMessage, a, b, a.VerseKey, err)
		return false
	}

	chapterB, verseB, err := qul.DecodeVerseKey(b.VerseKey)
	if err != nil {
		log.Println(compareInstancesErrorMessage, a, b, b.VerseKey, err)
		return false
	}

	if chapterA != chapterB {
		return chapterA < chapterB
	}
	if verseA != verseB {
		return verseA < verseB
	}
	return a.InstanceInVerse < b.InstanceInVerse
}
