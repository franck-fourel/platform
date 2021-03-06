package data

import (
	"github.com/tidepool-org/platform/structure"
)

type Datum interface {
	Meta() interface{}

	Parse(parser ObjectParser) error
	Validate(validator structure.Validator)
	Normalize(normalizer Normalizer)

	IdentityFields() ([]string, error)

	GetPayload() *Blob

	SetUserID(userID *string)
	SetDataSetID(dataSetID *string)
	SetActive(active bool)
	SetDeviceID(deviceID *string)
	SetCreatedTime(createdTime *string)
	SetCreatedUserID(createdUserID *string)
	SetModifiedTime(modifiedTime *string)
	SetModifiedUserID(modifiedUserID *string)
	SetDeletedTime(deletedTime *string)
	SetDeletedUserID(deletedUserID *string)

	DeduplicatorDescriptor() *DeduplicatorDescriptor
	SetDeduplicatorDescriptor(deduplicatorDescriptor *DeduplicatorDescriptor)
}

func DatumAsPointer(datum Datum) *Datum {
	return &datum
}
