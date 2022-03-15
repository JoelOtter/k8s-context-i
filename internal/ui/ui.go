package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/JoelOtter/k8s-context-i/internal/k8s"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type ui struct {
	screen   tcell.Screen
	contexts []k8s.Context
	pointer  int
	quit     chan struct{}
}

func (u *ui) drawStr(x int, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		u.screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (u *ui) draw() {
	u.screen.Clear()
	u.drawStr(1, 1, tcell.StyleDefault.Bold(true), "Kubernetes contexts:")
	for i, branch := range u.contexts {
		if branch.Current {
			u.screen.SetCell(1, i+3, tcell.StyleDefault, '*')
		}
		style := tcell.StyleDefault
		if branch.Current {
			style = style.Bold(true)
		}
		if i == u.pointer {
			style = style.Reverse(true)
		}
		u.drawStr(3, i+3, style, branch.Name)
	}
	u.screen.Show()
}

func (u *ui) keyDown() {
	u.pointer = (u.pointer + 1) % len(u.contexts)
	u.draw()
}

func (u *ui) keyUp() {
	u.pointer = u.pointer - 1
	if u.pointer < 0 {
		u.pointer = len(u.contexts) - 1
	}
	u.draw()
}

func (u *ui) run(uiOut io.Writer, uiErr *error) {
	defer close(u.quit)
	for {
		ev := u.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyEnter:
				*uiErr = k8s.ChangeContext(u.contexts[u.pointer].Name, uiOut)
				return
			case tcell.KeyUp, tcell.KeyPgUp, tcell.KeyCtrlP:
				u.keyUp()
			case tcell.KeyDown, tcell.KeyPgDn, tcell.KeyCtrlN:
				u.keyDown()
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'j':
					u.keyDown()
				case 'k':
					u.keyUp()
				}
			}
		case *tcell.EventResize:
			u.screen.Sync()
			u.draw()
		}
	}
}

func ShowUI(contexts []k8s.Context) error {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to get screen: %w", err)
	}
	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to init screen: %w", err)
	}

	var uiErr error
	uiOut := &strings.Builder{}
	defer func() {
		if uiOut.Len() > 0 {
			fmt.Print(uiOut.String())
		}
	}()

	u := &ui{
		screen:   screen,
		contexts: contexts,
		pointer:  getInitialPointer(contexts),
		quit:     make(chan struct{}),
	}
	u.draw()

	defer screen.Fini()

	go u.run(uiOut, &uiErr)

	for {
		select {
		case <-u.quit:
			return uiErr
		}
	}
}

func getInitialPointer(contexts []k8s.Context) int {
	for i, ctx := range contexts {
		if ctx.Current {
			return i
		}
	}
	return 0
}
