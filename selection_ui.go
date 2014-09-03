package selection_ui

import "github.com/nsf/termbox-go"

const (
	defaultBg = termbox.ColorBlack
	defaultFg = termbox.ColorWhite
	headerMsg = "Select one or more"
	footerMsg = "Move cursor with standard directionals; Select with SPACE; confirm with ENTER"
)

func Prompt(options []string, itemDesc string) (selections []bool) {
	ui := NewSelectionUI(options, itemDesc)

	if err := ui.Run(); err != nil {
		panic(err)
	}

	return ui.Selections
}

type SelectionUI struct {
	Options         []string // all the available options
	ItemDescription string   // string describing what one option is
	Selections      []bool   // 1:1 list indicating which were selected

	cursorIdx    int
	errorMessage string
}

func NewSelectionUI(options []string, itemDesc string) SelectionUI {
	return SelectionUI{
		Options:         options,
		ItemDescription: itemDesc,
		Selections:      make([]bool, len(options)),
	}
}

func (s *SelectionUI) Run() error {
	if err := s.init(); err != nil {
		return err
	}
	s.printAll()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			s.errorMessage = ""
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					break mainloop

				case termbox.KeyEnter:
					if isAnyTrue(s.Selections) {
						break mainloop
					} else {
						s.errorMessage = "Make a selection before continuing!"
					}

				case termbox.KeyArrowUp, termbox.KeyCtrlP:
					s.moveCursorUp()

				case termbox.KeyArrowDown, termbox.KeyCtrlN:
					s.moveCursorDown()

				case termbox.KeySpace:
					s.toggleSelectedUnderCursor()
				}

			} else {
				switch ev.Ch {
				case 'k':
					s.moveCursorUp()

				case 'j':
					s.moveCursorDown()

				case 'q':
					break mainloop
				}
			}

		case termbox.EventError:
			return ev.Err
		}

		s.printAll()
	}

	s.deinit()

	return nil
}

func (s *SelectionUI) init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)
	return nil
}

func (s *SelectionUI) deinit() {
	termbox.Close()
}

func (s *SelectionUI) setCursor(rowIdx int) {
	s.cursorIdx = rowIdx
}

func (s *SelectionUI) moveCursorUp() {
	if s.cursorIdx > 0 {
		s.setCursor(s.cursorIdx - 1)
	}
}

func (s *SelectionUI) moveCursorDown() {
	if s.cursorIdx < len(s.Options)-1 {
		s.setCursor(s.cursorIdx + 1)
	}
}

func (s *SelectionUI) toggleSelectedUnderCursor() {
	s.toggleSelected(s.cursorIdx)
}

func (s *SelectionUI) toggleSelected(row int) {
	s.Selections[row] = !s.Selections[row]
}

func (s *SelectionUI) printHeader(x int, y int) (newY int) {
	msg := headerMsg + " " + s.ItemDescription + ":"
	newY = printText(x, y, msg, termbox.ColorCyan, defaultBg)
	return
}

func (s *SelectionUI) printFooter(x, y int) (newY int) {
	newY = printText(x, y, footerMsg, termbox.ColorCyan, defaultBg)
	return
}

func (s *SelectionUI) printAll() {
	termbox.Clear(defaultFg, defaultBg)
	termbox.HideCursor()

	var y int
	y = s.printHeader(2, 1)
	y += 1
	y = s.printOptions(2, y)
	y += 1
	y = s.printFooter(2, y)
	y += 1

	if s.errorMessage != "" {
		y = printError(2, y, s.errorMessage)
	}

	termbox.Flush()
}

func (s *SelectionUI) printOptions(x, originY int) (newY int) {
	var bgColor, fgColor termbox.Attribute
	var y int = originY

	for i, option := range s.Options {
		if s.Selections[i] {
			bgColor = defaultBg
			fgColor = termbox.ColorGreen | termbox.AttrBold
		} else {
			bgColor = defaultBg
			fgColor = defaultFg
		}

		if i == s.cursorIdx {
			termbox.SetCell(x, y, '➜', defaultFg, defaultBg)
		}

		y = printText(x+2, y, option, fgColor, bgColor)
	}

	return y
}

func printText(originX, originY int, text string, fg, bg termbox.Attribute) int {
	var x, y int = originX, originY
	var width, _ = termbox.Size()
	wrapWidth := width - originX - 3

	for i, r := range text {
		y = originY + i/wrapWidth
		x = originX + i%wrapWidth
		if y > originY {
			x += 2
		}
		termbox.SetCell(x, y, r, fg, bg)
	}

	return y + 1
}

func printError(x, y int, msg string) int {
	return printText(x, y, msg, termbox.ColorRed, defaultBg)
}

func isAnyTrue(a []bool) bool {
	for _, b := range a {
		if b {
			return true
		}
	}
	return false
}
