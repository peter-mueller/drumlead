package leadsheet

type Leadsheet struct {
	Name  string
	Parts []Part
}

type LilypondDrummode string

type Part struct {
	Name    string
	Phrases []Phrase
}

type Phrase struct {
	Voices []LilypondDrummode
}
