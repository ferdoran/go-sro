package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

const ProgressWidth = 62

func PrintSection(section string) {
	log.Infoln("################################################################")
	var padDistance = (60 - len(section)) / 2
	sb := strings.Builder{}
	sb.WriteString("# ")
	for i := 0; i < padDistance; i++ {
		sb.WriteRune(' ')
	}
	sb.WriteString(section)
	for i := 0; i < padDistance; i++ {
		sb.WriteRune(' ')
	}
	sb.WriteString(" #")
	log.Infoln(sb.String())
	log.Infoln("################################################################")
}

func PrintProgress(current, max int) {
	count := int(float32(current) / float32(max) * ProgressWidth)
	sb := strings.Builder{}
	for i := 0; i < count; i++ {
		sb.WriteString("=")
	}
	if count < ProgressWidth {
		sb.WriteString(">")
	}
	fmt.Printf("\r[%-62s] %d/%d", sb.String(), current, max)
	if current >= max {
		fmt.Printf("\n")
	}
}
