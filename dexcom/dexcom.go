package dexcom

import (
	"regexp"
	"time"
)

const (
	AlertNameFixedLow   = "fixedLow"
	AlertNameLow        = "low"
	AlertNameHigh       = "high"
	AlertNameRise       = "rise"
	AlertNameFall       = "fall"
	AlertNameOutOfRange = "outOfRange"

	EventCarbs    = "carbs"
	EventExercise = "exercise"
	EventHealth   = "health"
	EventInsulin  = "insulin"

	ExerciseLight  = "light"
	ExerciseMedium = "medium"
	ExerciseHeavy  = "heavy"

	HealthIllness      = "illness"
	HealthStress       = "stress"
	HealthHighSymptoms = "highSymptoms"
	HealthLowSymptoms  = "lowSymptoms"
	HealthCycle        = "cycle"
	HealthAlcohol      = "alcohol"

	ModelG5MobileApp         = "G5 Mobile App"
	ModelG5Receiver          = "G5 Receiver"
	ModelG4WithShareReceiver = "G4 with Share Receiver"
	ModelG4Receiver          = "G4 Receiver"
	ModelUnknown             = "Unknown"

	StatusHigh             = "high"
	StatusLow              = "low"
	StatusOK               = "ok"
	StatusOutOfCalibration = "outOfCalibration"
	StatusSensorNoise      = "sensorNoise"

	TrendDoubleUp       = "doubleUp"
	TrendSingleUp       = "singleUp"
	TrendFortyFiveUp    = "fortyFiveUp"
	TrendFlat           = "flat"
	TrendFortyFiveDown  = "fortyFiveDown"
	TrendSingleDown     = "singleDown"
	TrendDoubleDown     = "doubleDown"
	TrendNone           = "none"
	TrendNotComputable  = "notComputable"
	TrendRateOutOfRange = "rateOutOfRange"

	UnitMinutes = "minutes"
	UnitGrams   = "grams"
	UnitUnits   = "units"

	UnitMgdL     = "mg/dL"
	UnitMmolL    = "mmol/L"
	UnitMgdLMin  = "mg/dL/min"
	UnitMmolLMin = "mmol/L/min"

	EGVValueMinMgdL = 40
	EGVValueMaxMgdL = 400

	DeviceIDMultiple           = "multiple"
	DeviceIDUnknown            = "unknown"
	DeviceModelMultiple        = "multiple"
	DeviceSerialNumberMultiple = "multiple"

	DateTimeFormat = "2006-01-02T15:04:05"
	NowThreshold   = 24 * time.Hour

	TransmitterIDExpressionString = "^[0-9A-Z]{5,6}$"
)

func AlertNames() []string {
	return []string{
		AlertNameFixedLow,
		AlertNameLow,
		AlertNameHigh,
		AlertNameRise,
		AlertNameFall,
		AlertNameOutOfRange,
	}
}

var transmitterIDExpression = regexp.MustCompile(TransmitterIDExpressionString)
