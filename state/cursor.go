package state

import tea "github.com/charmbracelet/bubbletea"

type Cursor struct {
	Index int
}

func (c *Cursor) MoveDown(length int) {
	if length == 0 {
		c.Index = 0
		return
	}
	if c.Index < length-1 {
		c.Index += 1
	} else {
		c.Index = 0
	}
}

func (c *Cursor) MoveUp(length int) {
	if length == 0 {
		c.Index = 0
		return
	}
	if c.Index > 0 {
		c.Index -= 1
	} else {
		c.Index = length - 1
	}
}

func (c *Cursor) HandleKey(msg tea.KeyMsg, listLength int) bool {
	if listLength <= 0 {
		c.Index = 0
		return false
	}

	original := c.Index

	switch msg.String() {
	case "j", tea.KeyDown.String():
		if c.Index < listLength-1 {
			c.Index += 1
		} else {
			c.Index = 0
		}
	case "k", tea.KeyUp.String():
		if c.Index > 0 {
			c.Index -= 1
		} else {
			c.Index = listLength - 1
		}
	case "h", tea.KeyLeft.String():
		c.Index = 0
	case "l", tea.KeyRight.String():
		c.Index = listLength - 1
	}

	return c.Index != original
}
