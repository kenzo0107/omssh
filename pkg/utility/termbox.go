package utility

import "github.com/nsf/termbox-go"

// TermboxKeys : return type string on termbox as []termbox.Event
func TermboxKeys(str string) []termbox.Event {
	s := []rune(str)
	e := make([]termbox.Event, 0, len(s))
	for _, r := range s {
		e = append(e, termbox.Event{Type: termbox.EventKey, Ch: r})
	}
	return e
}
