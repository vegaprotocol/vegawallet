package printer

import (
	"fmt"

	"github.com/muesli/termenv"
)

type HumanPrinter struct {
	profile   termenv.Profile
	checkMark string
	bangMark  string
	crossMark string
	arrow     termenv.Style
}

func NewHumanPrinter() *HumanPrinter {
	profile := termenv.ColorProfile()
	return &HumanPrinter{
		profile:   profile,
		checkMark: termenv.String("✓ ").Foreground(profile.Color("2")).String(),
		bangMark:  termenv.String("! ").Foreground(profile.Color("1")).String(),
		crossMark: termenv.String("✗ ").Foreground(profile.Color("1")).String(),
		arrow:     termenv.String("➜ ").Foreground(profile.Color("2")),
	}
}

func (p *HumanPrinter) GreenArrow() *HumanPrinter {
	printOut(p.arrow.Foreground(p.profile.Color("2")).String())
	return p
}

func (p *HumanPrinter) RedArrow() *HumanPrinter {
	printOut(p.arrow.Foreground(p.profile.Color("1")).String())
	return p
}

func (p *HumanPrinter) BlueArrow() *HumanPrinter {
	printOut(p.arrow.Foreground(p.profile.Color("6")).String())
	return p
}

func (p *HumanPrinter) CheckMark() *HumanPrinter {
	printOut(p.checkMark)
	return p
}

func (p *HumanPrinter) BangMark() *HumanPrinter {
	printOut(p.bangMark)
	return p
}

func (p *HumanPrinter) CrossMark() *HumanPrinter {
	printOut(p.crossMark)
	return p
}

func (p *HumanPrinter) SuccessText(t string) *HumanPrinter {
	printOut(termenv.String(t).Foreground(p.profile.Color("2")).String())
	return p
}

func (p *HumanPrinter) InfoText(t string) *HumanPrinter {
	printOut(termenv.String(t).Foreground(p.profile.Color("6")).String())
	return p
}

func (p *HumanPrinter) WarningText(t string) *HumanPrinter {
	printOut(termenv.String(t).Foreground(p.profile.Color("3")).String())
	return p
}

func (p *HumanPrinter) DangerText(t string) *HumanPrinter {
	printOut(termenv.String(t).Foreground(p.profile.Color("1")).String())
	return p
}

func (p *HumanPrinter) Jump() *HumanPrinter {
	printOut("\n")
	return p
}

func (p *HumanPrinter) NJump(num int) *HumanPrinter {
	var jumps string
	for i := 0; i < num; i++ {
		jumps += "\n"
	}
	printOut(jumps)
	return p
}

func (p *HumanPrinter) Text(s string) *HumanPrinter {
	printOut(s)
	return p
}

func (p *HumanPrinter) Code(s string) *HumanPrinter {
	printOut(fmt.Sprintf("    $ %s", s))
	return p
}

func (p *HumanPrinter) Bold(name string) *HumanPrinter {
	printOut(termenv.String(name).Bold().String())
	return p
}

func (p *HumanPrinter) Underline(name string) *HumanPrinter {
	printOut(termenv.String(name).Underline().String())
	return p
}

func printOut(s string) {
	fmt.Print(s) //nolint:forbidigo
}
