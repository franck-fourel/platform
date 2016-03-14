package data

//Basal represents a basal device data record
type Basal struct {
	DeliveryType string          `json:"deliveryType" bson:"deliveryType" valid:"required"`
	ScheduleName string          `json:"scheduleName" bson:"scheduleName" valid:"required"`
	Rate         float64         `json:"rate" bson:"rate" valid:"omitempty,gte=0"`
	Duration     int             `json:"duration" bson:"duration" valid:"omitempty,gte=0"`
	Suppressed   *SupressedBasal `json:"suppressed" bson:"suppressed,omitempty" valid:"omitempty"`
	Base         `bson:",inline"`
}

//SupressedBasal represents a suppressed basal portion of a basal
type SupressedBasal struct {
	Type         string  `json:"type" bson:"type" valid:"required"`
	DeliveryType string  `json:"deliveryType" bson:"deliveryType" valid:"required"`
	ScheduleName string  `json:"scheduleName" bson:"scheduleName" valid:"required"`
	Rate         float64 `json:"rate" bson:"rate" valid:"gte=0"`
}

const (
	//BasalName is the given name for the type of a `Basal` datum
	BasalName = "basal"

	deliveryTypeField = "deliveryType"
	scheduleNameField = "scheduleName"
	insulinField      = "insulin"
	rateField         = "rate"
	durationField     = "duration"
)

//BuildBasal will build a Basal record
func BuildBasal(obj map[string]interface{}) (*Basal, *Error) {

	base, errs := BuildBase(obj)
	cast := NewCaster(errs)

	basal := &Basal{
		DeliveryType: cast.ToString(deliveryTypeField, obj[deliveryTypeField]),
		ScheduleName: cast.ToString(scheduleNameField, obj[scheduleNameField]),
		Base:         base,
	}

	if obj[rateField] != nil {
		basal.Rate = cast.ToFloat64(rateField, obj[rateField])
	}
	if obj[durationField] != nil {
		basal.Duration = cast.ToInt(durationField, obj[durationField])
	}

	errs.AppendError(validator.ValidateStruct(basal))
	if errs.IsEmpty() {
		return basal, nil
	}
	return basal, errs
}

//Selector will return the `unique` fields used in upserts
func (b *Basal) Selector() interface{} {

	unique := map[string]interface{}{}
	unique[deliveryTypeField] = b.DeliveryType
	unique[scheduleNameField] = b.ScheduleName
	unique[deviceTimeField] = b.Time
	unique[typeField] = b.Type
	return unique
}
