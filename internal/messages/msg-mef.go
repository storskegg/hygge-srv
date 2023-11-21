package messages

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	patternMEF  = regexp.MustCompile(`^(?P<callsign>[\w\d]{4,6})\|(?P<msg_t>\d+)\|(?P<digest>[0-9a-fA-F]{8})\|(?P<payload>.*)$`)
	templateMEF = "callsign=$callsign::msg_t=$msg_t::digest=$digest::payload=$payload"
)

// MEF is the Message Encapsulation Format. It consists of a header and a payload. The header contains an amateur radio
// callsign (Callsign), a message type (MessageType), and a message digest (Digest). The payload is the remainder of the
// line being parsed, and may contain its own formatting, or it may be arbitrary.
//
// Example: W4PHO|1|84efdf24|4876|64.73|18.17|3.98
type MEF struct {
	Callsign    string // Amateur Radio Callsign
	MessageType uint32 // Message Type
	Digest      string // Hex-encoded uint32; agreed upon chunk of a 128-bit HMAC digest (this isn't critical infra, 32 bits will suffice)
	Data        *Hygge // ...
}

func ParseMEF(line string) (*MEF, error) {
	line = strings.TrimSpace(line)
	submatch := patternMEF.FindStringSubmatchIndex(line)
	if len(submatch) == 0 {
		return nil, fmt.Errorf("failed to parse MEF from line: %s", line)
	}
	capLine := string(patternMEF.ExpandString(nil, templateMEF, line, submatch))
	clParts := strings.Split(capLine, "::")

	mef := &MEF{}

	for _, kvPair := range clParts {
		kv := strings.Split(kvPair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid kvpair: %s", kvPair)
		}

		switch kv[0] {
		case "callsign":
			mef.Callsign = kv[1]
		case "msg_t":
			msgT, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid msg_t: %s", kv[1])
			}
			mef.MessageType = uint32(msgT)
		case "digest":
			mef.Digest = kv[1]
		case "payload":
			arotMsg, err := ParseHygge(kv[1]) // We only support one type of ARoT message at this time.
			if err != nil {
				return nil, fmt.Errorf("invalid Hygge: %s", kv[1])
			}
			mef.Data = arotMsg
		default:
			return nil, fmt.Errorf("invalid kvpair: %s", kvPair)
		}
	}

	return mef, nil
}
