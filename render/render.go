package render

import (
	"context"
	"embed"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/peter-mueller/drumlead/leadsheet"
)

//go:embed templates/*.gotmpl
var fs embed.FS

func SaveFile(l leadsheet.Leadsheet, filename string) {

	if !strings.HasSuffix(filename, ".pdf") {
		log.Fatalf("filename is not a pdf: %s", filename)
	}

	filename = filename[:len(filename)-len(".pdf")]

	ctx := context.Background()

	t, err := template.New("leadsheet.gotmpl").
		Funcs(template.FuncMap{"renderParallelVoices": renderParallelVoices}).
		ParseFS(fs, "templates/*.gotmpl")
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.CommandContext(ctx, "lilypond", "-o", filename, "-")
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(in, l)
	if err != nil {
		log.Fatalln(err)
	}
	in.Close()

	cmd.Wait()

}

func renderParallelVoices(voices []leadsheet.LilypondDrummode) (s string) {
	var (
		n = len(voices)
		b = strings.Builder{}
	)
	if n == 0 {
		return ""
	}

	estimatedVoiceLen := len(voices[0]) + 10
	b.Grow(estimatedVoiceLen * n)

	b.WriteString("<< {\n")
	b.WriteString(string(voices[0]))
	for _, voice := range voices[1:] {
		b.WriteString("\n} \\\\ {\n")
		b.WriteString(string(voice))
	}
	b.WriteString("\n} >>")
	return b.String()
}
