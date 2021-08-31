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
	fmt.Printf(p.arrow.Foreground(p.profile.Color("2")).String())
	return p
}

func (p *HumanPrinter) RedArrow() *HumanPrinter {
	fmt.Printf(p.arrow.Foreground(p.profile.Color("1")).String())
	return p
}

func (p *HumanPrinter) BlueArrow() *HumanPrinter {
	fmt.Printf(p.arrow.Foreground(p.profile.Color("6")).String())
	return p
}

func (p *HumanPrinter) CheckMark() *HumanPrinter {
	fmt.Printf(p.checkMark)
	return p
}

func (p *HumanPrinter) BangMark() *HumanPrinter {
	fmt.Printf(p.bangMark)
	return p
}

func (p *HumanPrinter) CrossMark() *HumanPrinter {
	fmt.Printf(p.crossMark)
	return p
}

func (p *HumanPrinter) SuccessText(t string) *HumanPrinter {
	fmt.Printf(termenv.String(t).Foreground(p.profile.Color("2")).String())
	return p
}

func (p *HumanPrinter) InfoText(t string) *HumanPrinter {
	fmt.Printf(termenv.String(t).Foreground(p.profile.Color("6")).String())
	return p
}

func (p *HumanPrinter) WarningText(t string) *HumanPrinter {
	fmt.Printf(termenv.String(t).Foreground(p.profile.Color("3")).String())
	return p
}

func (p *HumanPrinter) DangerText(t string) *HumanPrinter {
	fmt.Printf(termenv.String(t).Foreground(p.profile.Color("1")).String())
	return p
}

func (p *HumanPrinter) Jump() *HumanPrinter {
	fmt.Println()
	return p
}

func (p *HumanPrinter) NJump(num int) *HumanPrinter {
	var jumps string
	for i := 0; i < num; i++ {
		jumps += "\n"
	}
	fmt.Printf(jumps)
	return p
}

func (p *HumanPrinter) Text(s string) *HumanPrinter {
	fmt.Printf(s)
	return p
}

func (p *HumanPrinter) Code(s string) *HumanPrinter {
	fmt.Printf("    $ %s", s)
	return p
}

func (p *HumanPrinter) Bold(name string) *HumanPrinter {
	fmt.Printf(termenv.String(name).Bold().String())
	return p
}

func (p *HumanPrinter) Underline(name string) *HumanPrinter {
	fmt.Printf(termenv.String(name).Underline().String())
	return p
}
