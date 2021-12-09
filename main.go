package main

import (
	"github.com/peter-mueller/drumlead/parser"
	"github.com/peter-mueller/drumlead/render"
)

const file = `# Leadsheet 
## Intro
R1*4_"Gitarre" R1*8_"Gesang" R1*8_"Gesang"

## Verse
\repeat volta 2 { <sn toml>8 8 8 8 r2 r1_"x4" }

## Refrain 
\repeat volta 2 { toml8 hh sn8 hh toml8 hh sn8 hh toml8 hh sn8 hh toml8 hh sn8. sn16_"x4" }
bd4 bd bd bd | bd bd bd bd

## Trompete
\repeat volta 2 { <sn toml>8 8 8 8 r2 r1_"x4" }

## Verse
\repeat volta 2 { <sn toml>8 8 8 8 r2 r1_"x4" }

## Refrain 
\repeat volta 2 { toml8 hh sn8 hh toml8 hh sn8 hh toml8 hh sn8 hh toml8 hh sn8. sn16_"x4" }
bd4 bd bd bd | bd bd bd bd


## Trompete
\repeat volta 2 { <sn toml>8 8 8 8 r2 r1_"x4" }

## Gesang
\repeat volta 2 { <sn toml>8 8 8 8 r2 r1_"x3" }

`

func main() {
	l := parser.Parse(file)
	render.SaveFile(l, "out.pdf")
}
