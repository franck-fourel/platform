package device

import (
	. "github.com/tidepool-org/platform/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/tidepool-org/platform/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/tidepool-org/platform/data/_fixtures"
	"github.com/tidepool-org/platform/data/types"
)

var _ = Describe("DeviceEvent", func() {

	var helper *types.TestingHelper

	BeforeEach(func() {
		helper = types.NewTestingHelper()
	})

	Context("status", func() {

		var deviceEventObj = fixtures.TestingDatumBase()
		deviceEventObj["type"] = "deviceEvent"
		deviceEventObj["subType"] = "status"
		deviceEventObj["status"] = "suspended"
		deviceEventObj["reason"] = map[string]string{"suspended": "automatic"}

		It("returns a Status if the obj is valid", func() {
			Expect(helper.ValidDataType(Build(deviceEventObj, helper.ErrorProcessing))).To(BeNil())
		})

		Context("validation", func() {})
	})
})
