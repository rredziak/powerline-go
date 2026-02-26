package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

type TmuxState struct {
	clients  int
	sessions int
	attached int
	p        *powerline
}

func (ts *TmuxState) isOn() bool {
	return ts.clients+ts.sessions+ts.attached > 0
}

func (ts *TmuxState) status() string {
	var tOn string = fmt.Sprintf("%s %s", ts.p.symbols.Tmux, ts.p.symbols.TmuxOn)
	var tOff string = fmt.Sprintf("%s %s", ts.p.symbols.Tmux, ts.p.symbols.TmuxOff)

	if false { // if in full mode
		tOn = fmt.Sprintf("%s c:%2d s:%2d a:%2d", ts.p.symbols.Tmux, ts.clients, ts.sessions, ts.attached)
		tOff = ""
	}

	if ts.isOn() {
		return tOn
	}

	return tOff
}

func runTmuxCmd(parameters ...string) (string, error) {
	output, err := exec.Command("tmux", parameters...).Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func countLines(input string, countSum bool) (int, int) {
	var cntr, sum int

	lines := strings.Split(input, "\n")

	for _, l := range lines {

		if len(l) == 0 {
			continue
		}

		cntr++

		if countSum {
			val, _ := strconv.Atoi(l)
			sum = +val
		}
	}

	return cntr, sum
}

func updateClientNum(ts *TmuxState) {
	output, err := runTmuxCmd("list-clients", "-F", "#{session_windows}")

	if err == nil {
		cn, _ := countLines(output, false)
		ts.clients = cn
	}
}

func updateSessionNum(ts *TmuxState) {
	output, err := runTmuxCmd("list-sessions", "-F", "#{session_attached}")

	if err == nil {
		sn, sa := countLines(output, true)
		ts.sessions, ts.attached = sn, sa
	}
}

func segmentTmux(p *powerline) []pwl.Segment {
	ts := &TmuxState{p: p}

	if _, exists := os.LookupEnv("TMUX"); !exists {
		updateClientNum(ts)
		updateSessionNum(ts)

		if tc := ts.status(); len(tc) > 0 {
			return []pwl.Segment{{
				Name:       "tmux",
				Content:    tc,
				Foreground: p.theme.TmuxFg,
				Background: p.theme.TmuxBg,
			}}
		}
	}

	return []pwl.Segment{}
}
