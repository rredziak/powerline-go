package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentTmux(p *powerline) []pwl.Segment {
	return []pwl.Segment{{
		Name:       "tmux",
		Content:    "TMUX: stub",
		Foreground: p.theme.TmuxFg,
		Background: p.theme.TmuxBg,
	}}
}
