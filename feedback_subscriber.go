package gait

import (
	"context"
	"strings"

	mhist "github.com/alexmorten/mhist/proto"
)

// SubscribeToFeedback ...
func SubscribeToFeedback(c mhist.MhistClient, handler func(legName string, doneCommand string)) {
	stream, err := c.Subscribe(context.Background(), &mhist.Filter{Names: []string{"gait_feedback"}})
	if err != nil {
		panic(err)
	}

	for {
		m, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		if legName, command, ok := parseFeedbackMessage(m); ok {
			handler(legName, command)
		}
	}

}

func parseFeedbackMessage(m *mhist.MeasurementMessage) (legName, command string, ok bool) {
	if r := m.GetMeasurement().GetRaw(); r != nil {
		message := string(r.Value)
		return parseFeedback(message)
	}

	return "", "", false
}

func parseFeedback(message string) (legName, command string, ok bool) {
	splitLine := strings.SplitN(message, " ", 3)
	if len(splitLine) != 3 {
		return "", "", false
	}

	for i, part := range splitLine {
		splitLine[i] = strings.ToLower(removeNewLineChars(part))
	}

	legName = splitLine[0]
	verb := splitLine[1]
	command = splitLine[2]

	if verb != "did" {
		return "", "", false
	}

	return legName, command, true
}

func removeNewLineChars(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(s, "\r\n", ""),
		"\n",
		"",
	)
}
