package main

type TemplateDataPhrase struct {
	Count     int
	Chapters  int
	Phrase    []string
	Instances []Instance
}

type TemplateDataVerse struct {
	VerseKey        string
	InstanceInVerse int
	Count           int
	Chapters        int
	Phrase          []string
	Context         []string
	Continuation    []string
	Instances       []Instance
}

type Instance struct {
	VerseKey        string
	InstanceInVerse int
	Phrase          []string
	Context         []string
	Continuation    []string
}

type MediaEntry struct {
	Src string `yaml:"src"`
	As  string `yaml:"as"`
}
