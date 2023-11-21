package messages

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	patternHygge  = regexp.MustCompile(`^(?P<seq>\d+)\|(?P<humi>-?\d+\.\d{2})\|(?P<temp>-?\d+\.\d{2})\|(?P<batt>-?\d+\.\d{2})$`)
	templateHygge = "seq=$seq::humi=$humi::temp=$temp::batt=$batt"
)

// Hygge is the message type coming from the humidor. It consists of a packet sequence number (PacketSequence), humidity
// (Humidity) as % RH, temperature (Temperature) in Celsius, and battery voltage (Battery) in volts.
type Hygge struct {
	PacketSequence uint64
	Humidity       float64
	Temperature    float64
	Battery        float64
}

func ParseHygge(line string) (*Hygge, error) {
	line = strings.TrimSpace(line)
	submatch := patternHygge.FindStringSubmatchIndex(line)
	if len(submatch) == 0 {
		return nil, fmt.Errorf("failed to parse hygge from line: %s", line)
	}
	capLine := string(patternHygge.ExpandString(nil, templateHygge, strings.TrimSpace(line), submatch))
	clParts := strings.Split(capLine, "::")

	hg := &Hygge{}

	for _, kvPair := range clParts {
		kv := strings.Split(kvPair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid kvpair: %s", kvPair)
		}

		switch kv[0] {
		case "seq":
			seq, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid seq: %s", kv[1])
			}
			hg.PacketSequence = seq
		case "humi":
			humi, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid humi: %s", kv[1])
			}
			hg.Humidity = humi
		case "temp":
			temp, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid temp: %s", kv[1])
			}
			hg.Temperature = temp
		case "batt":
			batt, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid batt: %s", kv[1])
			}
			hg.Battery = batt
		default:
			return nil, fmt.Errorf("invalid kvpair: %s", kvPair)
		}
	}

	return hg, nil
}

func (hg *Hygge) String() string {
	return fmt.Sprintf("%d\t%1.2f\t%1.2f\t%1.2f", hg.PacketSequence, hg.Humidity, hg.Temperature, hg.Battery)
}
