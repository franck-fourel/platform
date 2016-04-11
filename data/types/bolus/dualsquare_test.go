package bolus

import (
	. "github.com/tidepool-org/platform/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/tidepool-org/platform/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/tidepool-org/platform/data/types"

	"github.com/tidepool-org/platform/data/_fixtures"
	"github.com/tidepool-org/platform/validate"
)

var _ = Describe("Bolus", func() {

	var bolusObj = fixtures.TestingDatumBase()
	bolusObj["type"] = "bolus"
	bolusObj["subType"] = "dual/square"
	bolusObj["normal"] = 2.0
	bolusObj["extended"] = 1.0
	bolusObj["duration"] = 3600000

	var processing validate.ErrorProcessing

	Context("dual/square from obj", func() {

		BeforeEach(func() {
			processing = validate.ErrorProcessing{BasePath: "0", ErrorsArray: validate.NewErrorsArray()}
		})

		It("if the obj is valid", func() {
			bolus := Build(bolusObj, processing)
			var bolusType *DualSquare
			Expect(bolus).To(BeAssignableToTypeOf(bolusType))
			Expect(processing.HasErrors()).To(BeFalse())
		})

		Context("validation", func() {

			Context("duration", func() {
				BeforeEach(func() {
					processing = validate.ErrorProcessing{BasePath: "0", ErrorsArray: validate.NewErrorsArray()}
				})

				It("is not required", func() {
					delete(bolusObj, "duration")
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeFalse())
				})

				It("invalid when less than zero", func() {
					bolusObj["duration"] = -1
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeTrue())
					Expect(processing.Errors[0].Detail).To(ContainSubstring("'Duration' failed with 'Must be greater than 0' when given '-1'"))
				})

				It("valid greater than zero", func() {
					bolusObj["duration"] = 4000
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeFalse())
				})

			})
			Context("extended", func() {
				BeforeEach(func() {
					processing = validate.ErrorProcessing{BasePath: "0", ErrorsArray: validate.NewErrorsArray()}
				})

				It("is not required", func() {
					delete(bolusObj, "extended")
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeFalse())
				})

				It("invalid when zero", func() {
					bolusObj["extended"] = -0.1
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeTrue())
					Expect(processing.Errors[0].Detail).To(ContainSubstring("'Extended' failed with 'Must be greater than 0.0' when given '-0.1'"))
				})

				It("valid when greater than zero", func() {
					bolusObj["extended"] = 0.7
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeFalse())
				})

			})
			Context("normal", func() {
				BeforeEach(func() {
					processing = validate.ErrorProcessing{BasePath: "0", ErrorsArray: validate.NewErrorsArray()}
				})

				It("is not required", func() {
					delete(bolusObj, "normal")
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeFalse())
				})

				It("invalid when zero", func() {
					bolusObj["normal"] = -0.1
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeTrue())
					Expect(processing.Errors[0].Detail).To(ContainSubstring("'Normal' failed with 'Must be greater than 0.0' when given '-0.1'"))
				})

				It("valid when greater than zero", func() {
					bolusObj["normal"] = 0.7
					bolus := Build(bolusObj, processing)
					types.GetPlatformValidator().Struct(bolus, processing)
					Expect(processing.HasErrors()).To(BeFalse())
				})

			})
		})
	})
})
