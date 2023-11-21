package messages

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	patternBridgeLine  = regexp.MustCompile(`^(?P<action>\w+)\s+\[(?P<rssi>-?\d+) RSSI\] -- (?P<data>.*)$`)
	templateBridgeLine = "action=$action::rssi=$rssi::data=$data"
)

// BridgeLine represents a raw line coming in from the lora bridge. It consists of a header of two parts, and a payload.
// The header consists of an Action, typically RECV, and an RSSI value for the received lora message.
//
// Example: RECV [-52 RSSI] -- W4PHO|1|84efdf24|4876|64.73|18.17|3.98
type BridgeLine struct {
	Action  string // example: RECV
	RSSI    int    // example: -52
	Message *MEF
}

func ParseBridgeLine(line string) (*BridgeLine, error) {
	line = strings.TrimSpace(line)
	submatch := patternBridgeLine.FindStringSubmatchIndex(strings.TrimSpace(line))
	if len(submatch) == 0 {
		return nil, fmt.Errorf("failed to parse BridgeLine from line: %s", line)
	}
	capLine := string(patternBridgeLine.ExpandString(nil, templateBridgeLine, line, submatch))
	clParts := strings.Split(capLine, "::")

	bl := &BridgeLine{}

	for _, kvPair := range clParts {
		kv := strings.Split(kvPair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid kvpair: %s", kvPair)
		}

		switch kv[0] {
		case "action":
			bl.Action = kv[1]
		case "rssi":
			rssi, err := strconv.Atoi(kv[1])
			if err != nil {
				return nil, fmt.Errorf("invalid rssi: %s", kv[1])
			}
			bl.RSSI = rssi
		case "data":
			mef, err := ParseMEF(kv[1])
			if err != nil {
				return nil, fmt.Errorf("invalid MEF: %s", kv[1])
			}
			bl.Message = mef
		default:
			return nil, fmt.Errorf("invalid kvpair: %s", kvPair)
		}
	}

	return bl, nil
}
