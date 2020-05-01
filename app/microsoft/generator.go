package microsoft

import (
	"github.com/connectors-test-client/app/common"
	"github.com/sirupsen/logrus"
)

var microsoftUsers = []string{
	"c2d82a71-98d8-4645-9e88-e5531ed1ca02", // Greg
	"627f153d-7dfc-44f5-9f66-4393de2bfe54", // Haarika
	"2f45be53-47dd-46df-9e47-8c442aaf48ed", // Tim
}

// GenerateTest - generates the test for microsoft
func GenerateTest(settings *common.GeneratorSettings) {
	// Impliment this!
	logrus.Printf("Acct:\t%s\n", settings.Account)
	logrus.Printf("Type:\t%s\n", settings.StreamType)
	logrus.Printf("Count:\t%v\n", settings.TransCount)
}
