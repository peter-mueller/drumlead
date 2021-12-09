package parser

import (
	"log"

	"github.com/peter-mueller/drumlead/leadsheet"
)

type parser struct {
	currentHashes int

	value leadsheet.Leadsheet
}

func Parse(input string) leadsheet.Leadsheet {
	items := Lex(input)

	p := parser{}
	for it := range filterDoubleEmptyText(items) {
		switch it.typ {
		case HASHTAG:
			p.currentHashes += 1
		case TEXT:
			p.acceptText(it.value)
		}
	}
	return p.value
}

func filterDoubleEmptyText(in <-chan item) <-chan item {
	c := make(chan item)
	go func() {
		defer close(c)
		lastWasEmptyText := false
		for i := range in {
			isEmptyText := i.typ == TEXT && i.value == ""
			if isEmptyText && lastWasEmptyText {
				continue
			}
			lastWasEmptyText = isEmptyText
			c <- i
		}
	}()
	return c
}

func (p *parser) acceptText(text string) {
	switch p.currentHashes {
	case 1:
		if p.value.Name != "" {
			log.Fatal("only set top level leadsheet once")
		}
		p.value.Name = text
		p.currentHashes = 0
		return
	case 2:
		p.addNewPart(text)
		p.currentHashes = 0
		return
	default:
		if text == "" {
			p.addNewPhrase()
			return
		}
		lilyLine := leadsheet.LilypondDrummode(text)
		lastPhrase := p.lastPhrase()
		lastPhrase.Voices = append(lastPhrase.Voices, lilyLine)
	}
}

func (p *parser) lastPart() *leadsheet.Part {
	return &p.value.Parts[len(p.value.Parts)-1]
}

func (p *parser) lastPhrase() *leadsheet.Phrase {
	lastPart := p.lastPart()
	return &lastPart.Phrases[len(lastPart.Phrases)-1]
}

func (p *parser) addNewPart(name string) {
	newPart := leadsheet.Part{
		Name: name,
		Phrases: []leadsheet.Phrase{
			{Voices: make([]leadsheet.LilypondDrummode, 0)},
		},
	}
	p.value.Parts = append(p.value.Parts, newPart)
}

func (p *parser) addNewPhrase() {
	newPhrase := leadsheet.Phrase{Voices: make([]leadsheet.LilypondDrummode, 0)}
	lastPart := p.lastPart()
	lastPart.Phrases = append(lastPart.Phrases, newPhrase)
}
