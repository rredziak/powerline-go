package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

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

		if !countSum {
			continue
		}

		val, _ := strconv.Atoi(l)
		sum = +val
	}

	return cntr, sum
}

func getClientNum() (int, error) {
	output, err := runTmuxCmd("list-clients", "-F", "#{session_windows}")

	if err != nil {
		return 0, err
	}

	clientNum, _ := countLines(output, false)

	return clientNum, nil
}

func getSessionNum() (int, int, error) {

	output, err := runTmuxCmd("list-sessions", "-F", "#{session_attached}")

	if err != nil {
		return 0, 0, err
	}

	sessionNum, sessionsAttached := countLines(output, true)
	return sessionNum, sessionsAttached, nil
}

func segmentTmux(p *powerline) []pwl.Segment {

	if _, exists := os.LookupEnv("TMUX"); exists {
		return []pwl.Segment{}
	}

	clientNum, err := getClientNum()

	if err != nil {
		return []pwl.Segment{}
	}

	sessionNum, sessionsAttached, err := getSessionNum()

	if err != nil {
		return []pwl.Segment{}
	}

	tmuxContent := fmt.Sprintf("%s c:%2d s:%2d a:%2d", p.symbols.Tmux, clientNum, sessionNum, sessionsAttached)

	return []pwl.Segment{{
		Name:       "tmux",
		Content:    tmuxContent,
		Foreground: p.theme.TmuxFg,
		Background: p.theme.TmuxBg,
	}}
}
