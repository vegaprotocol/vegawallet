package printer

import (
	"fmt"
	"io"
	"runtime"

	"github.com/muesli/termenv"
)

type InteractivePrinter struct {
	writer       io.Writer
	profile      termenv.Profile
	checkMark    string
	bangMark     string
	crossMark    string
	questionMark string
	arrow        termenv.Style
}

func fancyPrinter(w io.Writer, profile termenv.Profile) *InteractivePrinter {
	return &InteractivePrinter{
		writer:       w,
		profile:      profile,
		checkMark:    termenv.String("✓ ").Foreground(profile.Color("2")).String(),
		bangMark:     termenv.String("! ").Foreground(profile.Color("1")).String(),
		questionMark: termenv.String("? ").Foreground(profile.Color("2")).String(),
		crossMark:    termenv.String("✗ ").Foreground(profile.Color("1")).String(),
		arrow:        termenv.String("➜ ").Foreground(profile.Color("2")),
	}
}

func ansiPrinter(w io.Writer, profile termenv.Profile) *InteractivePrinter {
	return &InteractivePrinter{
		writer:    w,
		profile:   profile,
		checkMark: termenv.String("* ").Foreground(profile.Color("2")).String(),
		bangMark:  termenv.String("! ").Foreground(profile.Color("1")).String(),
		crossMark: termenv.String("x ").Foreground(profile.Color("1")).String(),
		arrow:     termenv.String("> ").Foreground(profile.Color("2")),
	}
}

func NewInteractivePrinter(w io.Writer) *InteractivePrinter {
	enableLegacyWindowsANSI()
	profile := termenv.EnvColorProfile()
	if runtime.GOOS == "windows" {
		return ansiPrinter(w, profile)
	}
	return fancyPrinter(w, profile)
}

func (p *InteractivePrinter) GreenArrow() *InteractivePrinter {
	p.printOut(p.arrow.Foreground(p.profile.Color("2")).String())
	return p
}

func (p *InteractivePrinter) RedArrow() *InteractivePrinter {
	p.printOut(p.arrow.Foreground(p.profile.Color("1")).String())
	return p
}

func (p *InteractivePrinter) BlueArrow() *InteractivePrinter {
	p.printOut(p.arrow.Foreground(p.profile.Color("6")).String())
	return p
}

func (p *InteractivePrinter) CheckMark() *InteractivePrinter {
	p.printOut(p.checkMark)
	return p
}

func (p *InteractivePrinter) BangMark() *InteractivePrinter {
	p.printOut(p.bangMark)
	return p
}

func (p *InteractivePrinter) QuestionMark() *InteractivePrinter {
	p.printOut(p.questionMark)
	return p
}

func (p *InteractivePrinter) CrossMark() *InteractivePrinter {
	p.printOut(p.crossMark)
	return p
}

func (p *InteractivePrinter) SuccessText(t string) *InteractivePrinter {
	p.printOut(termenv.String(t).Foreground(p.profile.Color("2")).String())
	return p
}

func (p *InteractivePrinter) InfoText(t string) *InteractivePrinter {
	p.printOut(termenv.String(t).Foreground(p.profile.Color("6")).String())
	return p
}

func (p *InteractivePrinter) WarningText(t string) *InteractivePrinter {
	p.printOut(termenv.String(t).Foreground(p.profile.Color("3")).String())
	return p
}

func (p *InteractivePrinter) DangerText(t string) *InteractivePrinter {
	p.printOut(termenv.String(t).Foreground(p.profile.Color("1")).String())
	return p
}

func (p *InteractivePrinter) NextLine() *InteractivePrinter {
	p.printOut("\n")
	return p
}

func (p *InteractivePrinter) NextSection() *InteractivePrinter {
	p.printOut("\n\n")
	return p
}

func (p *InteractivePrinter) Text(s string) *InteractivePrinter {
	p.printOut(s)
	return p
}

func (p *InteractivePrinter) Code(s string) *InteractivePrinter {
	p.printOut(fmt.Sprintf("    $ %s", s))
	return p
}

func (p *InteractivePrinter) Bold(name string) *InteractivePrinter {
	p.printOut(termenv.String(name).Bold().String())
	return p
}

func (p *InteractivePrinter) DangerBold(name string) *InteractivePrinter {
	p.printOut(termenv.String(name).Bold().Foreground(p.profile.Color("1")).String())
	return p
}

func (p *InteractivePrinter) SuccessBold(name string) *InteractivePrinter {
	p.printOut(termenv.String(name).Bold().Foreground(p.profile.Color("2")).String())
	return p
}

func (p *InteractivePrinter) Underline(name string) *InteractivePrinter {
	p.printOut(termenv.String(name).Underline().String())
	return p
}

func (p *InteractivePrinter) printOut(s string) {
	if _, err := fmt.Fprint(p.writer, s); err != nil {
		panic(fmt.Sprintf("couldn't write to %v: %v", p.writer, err))
	}
}
