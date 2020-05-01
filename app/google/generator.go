package google

import (
	"github.com/connectors-test-client/app/common"
	"github.com/sirupsen/logrus"
)

// GenerateTest - generates the test for google
func GenerateTest(settings *common.GeneratorSettings) {
	// Impliment this!
	logrus.Printf("Acct:\t%s\n", settings.Account)
	logrus.Printf("Type:\t%s\n", settings.StreamType)
	logrus.Printf("Count:\t%v\n", settings.TransCount)
}
