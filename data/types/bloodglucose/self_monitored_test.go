package bloodglucose_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fixtures "github.com/tidepool-org/platform/data/_fixtures"
	"github.com/tidepool-org/platform/data/types"
	"github.com/tidepool-org/platform/data/types/bloodglucose"
)

var _ = Describe("Selfmonitored", func() {

	var helper *types.TestingHelper
	var bgObj types.Datum
	var mmolL = "mmol/L"

	BeforeEach(func() {
		helper = types.NewTestingHelper()
	})

	Context("smbg from obj", func() {

		BeforeEach(func() {
			bgObj = fixtures.TestingDatumBase()
			bgObj["type"] = "smbg"
			bgObj["value"] = 5.5
			bgObj["units"] = "mmol/l"
		})

		It("returns a bolus if the obj is valid", func() {
			Expect(helper.ValidDataType(bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing))).To(BeNil())
		})

	})
	Context("validation", func() {

		BeforeEach(func() {
			bgObj = fixtures.TestingDatumBase()
			bgObj["type"] = "smbg"
			bgObj["value"] = 5.5
			bgObj["units"] = "mmol/l"
		})
		Context("units", func() {
			It("is required", func() {
				delete(bgObj, "units")

				Expect(
					helper.ErrorIsExpected(
						bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing),
						types.ExpectedErrorDetails{
							Path:   "0/units",
							Detail: "Must be one of mmol/L, mg/dL given '<nil>'",
						}),
				).To(BeNil())
			})

			It("can be mmol/l but saved as mmol/L", func() {
				bgObj["units"] = "mmol/l"
				smbg := bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing)
				Expect(helper.ValidDataType(smbg)).To(BeNil())
				Expect(smbg.Units).To(Equal(&mmolL))
			})

			It("can be mg/dl but saved as mmol/L", func() {
				bgObj["units"] = "mg/dl"

				smbg := bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing)
				Expect(helper.ValidDataType(smbg)).To(BeNil())
				Expect(smbg.Units).To(Equal(&mmolL))
			})

			It("cannot be anything else", func() {
				bgObj["units"] = "grams"
				Expect(
					helper.ErrorIsExpected(
						bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing),
						types.ExpectedErrorDetails{
							Path:   "0/units",
							Detail: "Must be one of mmol/L, mg/dL given 'grams'",
						}),
				).To(BeNil())
			})

		})
		Context("value", func() {
			It("is required", func() {
				delete(bgObj, "value")
				Expect(
					helper.ErrorIsExpected(
						bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing),
						types.ExpectedErrorDetails{
							Path:   "0/value",
							Detail: "Must be between 0.0 and 55.0 given '<nil>'",
						}),
				).To(BeNil())
			})

			It("if mg/dL < 1000", func() {

				bgObj["value"] = 1000.1
				bgObj["units"] = "mg/dL"

				Expect(
					helper.ErrorIsExpected(
						bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing),
						types.ExpectedErrorDetails{
							Path:   "0/value",
							Detail: "Must be between 0.0 and 1000.0 given '1000.1'",
						}),
				).To(BeNil())
			})

			It("if mmol/L < 55", func() {

				bgObj["value"] = 55.1
				bgObj["units"] = "mmol/L"

				Expect(
					helper.ErrorIsExpected(
						bloodglucose.BuildSelfMonitored(bgObj, helper.ErrorProcessing),
						types.ExpectedErrorDetails{
							Path:   "0/value",
							Detail: "Must be between 0.0 and 55.0 given '55.1'",
						}),
				).To(BeNil())
			})

		})
	})
})