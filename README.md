# Introduction

A program to generate Anki flashcards for Mutashābihāt (similar/confusing verses) in the Qurʾān.

It is available on AnkiWeb as [mayyize](https://ankiweb.net/shared/info/1065363749).

# Installation

```sh
go install git.sr.ht/~rehandaphedar/mayyize@latest
```

The package helps interact with the [Quranic Universal Library (QUL)](https://qul.tarteel.ai/resources/quran-metadata).

# Usage

The documentation for usage and flags can be accessed by running `mayyize -h`.

- The `-phrases` data can be obtained from QUL's [Mutashabihat ul Quran - mutashabihat(Phrase)](https://qul.tarteel.ai/resources/mutashabihat/73)
- The `-words` data can be obtained from QUL's [Ayah by ayah and word by text of Quran](https://qul.tarteel.ai/resources/quran-script)
- The `-layout` data can be obtained from QUL's [Mushaf Layout Resources](https://qul.tarteel.ai/resources/mushaf-layout)
- The `-metadata-*` can be obtained from QUL's [Quran data, surahs, ayahs, words, juz etc.](https://qul.tarteel.ai/resources/quran-metadata)
- The `-media-config` is a YAML file with a list of objects with the keys `src` and `as`. The filepaths are resolved relative to the config file.

# Note Types

## Phrase

This note type produces the Phrase Recognition card type.
The front of the card shows the phrase, while the back of the card shows all instances of the phrase.

![Phrase front](https://git.sr.ht/~rehandaphedar/mayyize/blob/main/assets/phrase-front.png)
![Phrase back](https://git.sr.ht/~rehandaphedar/mayyize/blob/main/assets/phrase-back.png)

## Verse

This note type produces the Verse Recognition card type.
The front of the card shows the context + phrase, while the back of the card shows the context + phrase + continuation, as well as all instances of the phrase.

![Verse front](https://git.sr.ht/~rehandaphedar/mayyize/blob/main/assets/verse-front.png)
![Verse back](https://git.sr.ht/~rehandaphedar/mayyize/blob/main/assets/verse-back.png)

# Recommended Usage

It is recommended to mainly use the verse deck, suspending all notes first and unsuspending by tag as you memorise.

The phrase deck is mainly intended to serve as a reference. It is generally unfeasible to memorise *all* instances of a phrase. However, you can sort by the `Count` field and selectively memorise phrases with a lesser number of instances.
