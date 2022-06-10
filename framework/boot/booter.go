package boot

import (
	"github.com/sirupsen/logrus"
	"strings"
)

type bootComponent struct {
	bootFunc          func()
	amountOfInstances int
}

var _phases = make([]string, 0)
var _componentsToBoot = make(map[string][]bootComponent)

func SetPhases(phases ...string) {
	_phases = toUpper(phases)
}

func RegisterComponent(phase string, bootFunc func(), times int) {
	bootConf := bootComponent{bootFunc, times}
	upperPhase := strings.ToUpper(phase)

	if !phaseExists(upperPhase) {
		logrus.Warnf("registering component for yet unknown phase %s", upperPhase)
	}

	_componentsToBoot[upperPhase] = append(_componentsToBoot[upperPhase], bootConf)
}

func Boot() {
	for _, phase := range _phases {
		logrus.Infof("starting booting phase %s", phase)
		for _, component := range _componentsToBoot[phase] {
			for i := 0; i < component.amountOfInstances; i++ {
				component.bootFunc()
			}
		}
	}
}

func phaseExists(phase string) bool {
	for _, p := range _phases {
		if p == phase {
			return true
		}
	}

	return false
}

func toUpper(phases []string) []string {
	upper := make([]string, 0)

	for _, phase := range phases {
		upper = append(upper, strings.ToUpper(phase))
	}

	return upper
}
