package blob

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/tidepool-org/platform/crypto"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/id"
	"github.com/tidepool-org/platform/net"
	"github.com/tidepool-org/platform/page"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/request"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/user"
)

const (
	ErrorCodeDigestsNotEqual = "digests-not-equal"

	StatusAvailable = "available"
	StatusCreated   = "created"
)

func ErrorDigestsNotEqual(value string, calculated string) error {
	return errors.Preparedf(ErrorCodeDigestsNotEqual, "digests not equal", "digest %q does not equal calculated digest %q", value, calculated)
}

func Statuses() []string {
	return []string{
		StatusAvailable,
		StatusCreated,
	}
}

type Client interface {
	List(ctx context.Context, userID string, filter *Filter, pagination *page.Pagination) (Blobs, error)
	Create(ctx context.Context, userID string, create *Create) (*Blob, error)
	Get(ctx context.Context, id string) (*Blob, error)
	GetContent(ctx context.Context, id string) (*Content, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type Filter struct {
	MediaType *[]string `json:"mediaType,omitempty"`
	Status    *[]string `json:"status,omitempty"`
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Parse(parser structure.ObjectParser) {
	if value := parser.StringArray("mediaType"); value != nil {
		f.MediaType = value
	}
	if value := parser.StringArray("status"); value != nil {
		f.Status = value
	}
}

func (f *Filter) Validate(validator structure.Validator) {
	validator.StringArray("mediaType", f.MediaType).NotEmpty().Each(func(stringValidator structure.String) { stringValidator.Using(net.MediaTypeValidator) }).EachUnique()
	validator.StringArray("status", f.Status).NotEmpty().EachOneOf(Statuses()...).EachUnique()
}

func (f *Filter) MutateRequest(req *http.Request) error {
	parameters := map[string][]string{}
	if f.MediaType != nil {
		parameters["mediaType"] = *f.MediaType
	}
	if f.Status != nil {
		parameters["status"] = *f.Status
	}
	return request.NewArrayParametersMutator(parameters).MutateRequest(req)
}

type Create struct {
	Body      io.Reader
	DigestMD5 *string
	MediaType *string
}

func NewCreate() *Create {
	return &Create{}
}

func (c *Create) Validate(validator structure.Validator) {
	if c.Body == nil {
		validator.WithReference("body").ReportError(structureValidator.ErrorValueNotExists())
	}
	validator.String("digestMD5", c.DigestMD5).Using(crypto.Base64EncodedMD5HashValidator)
	validator.String("mediaType", c.MediaType).Exists().Using(net.MediaTypeValidator)
}

type Content struct {
	Body      io.ReadCloser
	DigestMD5 *string
	MediaType *string
	Size      *int
}

func NewContent() *Content {
	return &Content{}
}

func (c *Content) Validate(validator structure.Validator) {
	if c.Body == nil {
		validator.WithReference("body").ReportError(structureValidator.ErrorValueNotExists())
	}
	validator.String("digestMD5", c.DigestMD5).Exists().Using(crypto.Base64EncodedMD5HashValidator)
	validator.String("mediaType", c.MediaType).Exists().Using(net.MediaTypeValidator)
	validator.Int("size", c.Size).Exists().GreaterThanOrEqualTo(0)
}

type Blob struct {
	ID           *string    `json:"id,omitempty" bson:"id,omitempty"`
	UserID       *string    `json:"userId,omitempty" bson:"userId,omitempty"`
	DigestMD5    *string    `json:"digestMD5,omitempty" bson:"digestMD5,omitempty"`
	MediaType    *string    `json:"mediaType,omitempty" bson:"mediaType,omitempty"`
	Size         *int       `json:"size,omitempty" bson:"size,omitempty"`
	Status       *string    `json:"status,omitempty" bson:"status,omitempty"`
	CreatedTime  *time.Time `json:"createdTime,omitempty" bson:"createdTime,omitempty"`
	ModifiedTime *time.Time `json:"modifiedTime,omitempty" bson:"modifiedTime,omitempty"`
}

func (b *Blob) Parse(parser structure.ObjectParser) {
	b.ID = parser.String("id")
	b.UserID = parser.String("userId")
	b.DigestMD5 = parser.String("digestMD5")
	b.MediaType = parser.String("mediaType")
	b.Size = parser.Int("size")
	b.Status = parser.String("status")
	b.CreatedTime = parser.Time("createdTime", time.RFC3339)
	b.ModifiedTime = parser.Time("modifiedTime", time.RFC3339)
}

func (b *Blob) Validate(validator structure.Validator) {
	validator.String("id", b.ID).Exists().Using(IDValidator)
	validator.String("userId", b.UserID).Exists().Using(user.IDValidator)
	validator.String("digestMD5", b.DigestMD5).Exists().Using(crypto.Base64EncodedMD5HashValidator)
	validator.String("mediaType", b.MediaType).Exists().Using(net.MediaTypeValidator)
	validator.Int("size", b.Size).Exists().GreaterThanOrEqualTo(0)
	validator.String("status", b.Status).Exists().OneOf(Statuses()...)
	validator.Time("createdTime", b.CreatedTime).Exists().NotZero().BeforeNow(time.Second)
	validator.Time("modifiedTime", b.ModifiedTime).After(pointer.ToTime(b.CreatedTime)).BeforeNow(time.Second)
}

type Blobs []*Blob

func NewID() string {
	return id.Must(id.New(16))
}

func IsValidID(value string) bool {
	return ValidateID(value) == nil
}

func IDValidator(value string, errorReporter structure.ErrorReporter) {
	errorReporter.ReportError(ValidateID(value))
}

func ValidateID(value string) error {
	if value == "" {
		return structureValidator.ErrorValueEmpty()
	} else if !idExpression.MatchString(value) {
		return ErrorValueStringAsIDNotValid(value)
	}
	return nil
}

func ErrorValueStringAsIDNotValid(value string) error {
	return errors.Preparedf(structureValidator.ErrorCodeValueNotValid, "value is not valid", "value %q is not valid as blob id", value)
}

var idExpression = regexp.MustCompile("^[0-9a-z]{32}$")
