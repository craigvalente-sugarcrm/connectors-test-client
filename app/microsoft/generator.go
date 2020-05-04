package microsoft

import (
	"encoding/json"
	"io/ioutil"
	"math"

	"github.com/connectors-test-client/app/common"
	"github.com/sirupsen/logrus"
)

// Account stores test account data

var msAccounts = []common.Account{
	{
		ID:       "c2d82a71-98d8-4645-9e88-e5531ed1ca02",
		Email:    "glevine@sugarcrm.com",
		Calendar: "SugarCRM",
	},
	{
		ID:       "627f153d-7dfc-44f5-9f66-4393de2bfe54",
		Email:    "hmaddipatla@sugarcrm.com",
		Calendar: "SugarCRM",
	},
	{
		ID:       "2f45be53-47dd-46df-9e47-8c442aaf48ed",
		Email:    "twolf@sugarcrm.com",
		Calendar: "SugarCRM",
	},
}

// GenerateTest - generates the test for microsoft
func GenerateTest(settings *common.GeneratorSettings) {
	// Impliment this!
	logrus.Printf("Acct:\t%s\n", settings.Account)
	logrus.Printf("Type:\t%s\n", settings.StreamType)
	logrus.Printf("Count:\t%v\n", settings.TransCount)

	testData := []*common.Action{}
	for _, a := range msAccounts {
		if settings.Account == "all" || a.Email == settings.Account {
			acct := &common.Action{
				Account: a,
				Items:   []*common.Projection{},
			}
			testData = append(testData, acct)
		}
	}

	transPerAcct := int(math.Ceil(float64(settings.TransCount) / float64(len(testData))))
	createsPerAcct := int(math.Ceil(float64(transPerAcct) * 0.20))
	updatesPerAcct := transPerAcct - createsPerAcct

	for i := 0; i < len(testData); i++ {
		action := testData[i]
		projMap := make(map[string]*common.Projection)
		for j := 0; j < createsPerAcct; j++ {
			proj := action.GenerateProjection()
			projMap[proj.CaseID] = proj
			action.Items = append(action.Items, proj)
		}
		for j := 0; j < updatesPerAcct; j += createsPerAcct {
			for _, proj := range projMap {
				updated := action.UpdateProjection(proj)
				action.Items = append(action.Items, updated)
			}
		}
	}

	jsonData, _ := json.MarshalIndent(testData, "", "  ")
	// jsonData, _ := json.Marshal(testData)
	_ = ioutil.WriteFile("test.json", jsonData, 0644)
}
