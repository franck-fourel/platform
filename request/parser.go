package request

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidepool-org/platform/crypto"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/net"
	"github.com/tidepool-org/platform/structure"
	structureNormalizer "github.com/tidepool-org/platform/structure/normalizer"
	structureParser "github.com/tidepool-org/platform/structure/parser"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
)

func DecodeRequestBody(req *http.Request, object interface{}) error {
	if req == nil {
		return errors.New("request is missing")
	}
	if req.Body == nil {
		return errors.New("request body is missing")
	}

	defer req.Body.Close()
	return DecodeObject(structure.NewPointerSource(), req.Body, object)
}

func DecodeResponseBody(res *http.Response, object interface{}) error {
	if res == nil {
		return errors.New("response is missing")
	}
	if res.Body == nil {
		return errors.New("response body is missing")
	}

	defer res.Body.Close()
	return DecodeObject(structure.NewPointerSource(), res.Body, object)
}

func DecodeObject(source structure.Source, reader io.Reader, object interface{}) error {
	if err := ParseStreamObject(source, reader, object); err != nil {
		return err
	}
	if err := ValidateObjects(source, object); err != nil {
		return err
	}
	return NormalizeObjects(source, object)
}

func ParseStreamObject(source structure.Source, reader io.Reader, object interface{}) error {
	if objectParsable, ok := object.(structure.ObjectParsable); ok {
		return ParseObjectParseableStreamObject(source, reader, objectParsable)
	}
	if arrayParsable, ok := object.(structure.ArrayParsable); ok {
		return ParseArrayParseableStreamObject(source, reader, arrayParsable)
	}
	return ParseSimpleStreamObject(reader, object)
}

func ParseObjectParseableStreamObject(source structure.Source, reader io.Reader, objectParsable structure.ObjectParsable) error {
	object := &map[string]interface{}{}
	if err := ParseSimpleStreamObject(reader, object); err != nil {
		return err
	}

	parser := structureParser.NewObject(object).WithSource(source)
	objectParsable.Parse(parser)
	parser.NotParsed()

	return parser.Error()
}

func ParseArrayParseableStreamObject(source structure.Source, reader io.Reader, arrayParsable structure.ArrayParsable) error {
	array := &[]interface{}{}
	if err := ParseSimpleStreamObject(reader, array); err != nil {
		return err
	}

	parser := structureParser.NewArray(array).WithSource(source)
	arrayParsable.Parse(parser)
	parser.NotParsed()

	return parser.Error()
}

func ParseSimpleStreamObject(reader io.Reader, object interface{}) error {
	if reader == nil {
		return errors.New("reader is missing")
	}
	if object == nil {
		return errors.New("object is missing")
	}

	if err := json.NewDecoder(reader).Decode(object); err != nil {
		return errors.Wrap(err, "json is malformed")
		// return ErrorJSONMalformed()
	}

	return nil
}

func ValidateObjects(source structure.Source, objects ...interface{}) error {
	validatables := []structure.Validatable{}
	for _, object := range objects {
		if validatable, ok := object.(structure.Validatable); ok {
			validatables = append(validatables, validatable)
		}
	}

	validator := structureValidator.New().WithSource(source)
	for _, validatable := range validatables {
		validatable.Validate(validator)
	}
	return validator.Error()
}

func NormalizeObjects(source structure.Source, objects ...interface{}) error {
	normalizables := []structure.Normalizable{}
	for _, object := range objects {
		if normalizable, ok := object.(structure.Normalizable); ok {
			normalizables = append(normalizables, normalizable)
		}
	}

	normalizer := structureNormalizer.New().WithSource(source)
	for _, normalizable := range normalizables {
		normalizable.Normalize(normalizer)
	}
	return normalizer.Error()
}

func DecodeRequestQuery(req *http.Request, objectParsables ...structure.ObjectParsable) error {
	if req == nil {
		return errors.New("request is missing")
	}
	if req.URL == nil {
		return errors.New("request url is missing")
	}

	values, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return errors.New("unable to parse request query")
	}

	return DecodeValues((map[string][]string)(values), objectParsables...)
}

func DecodeValues(values map[string][]string, objectParsables ...structure.ObjectParsable) error {
	objects := []interface{}{}
	for _, object := range objectParsables {
		objects = append(objects, object)
	}

	if err := ParseValuesObjects(values, objectParsables...); err != nil {
		return err
	}
	if err := ValidateObjects(structure.NewParameterSource(), objects...); err != nil {
		return err
	}
	return NormalizeObjects(structure.NewParameterSource(), objects...)
}

func ParseValuesObjects(values map[string][]string, objectParsables ...structure.ObjectParsable) error {
	parser := NewValues(&values)
	for _, objectParsable := range objectParsables {
		objectParsable.Parse(parser)
	}
	parser.NotParsed()
	return parser.Error()
}

func ParseDigestMD5Header(header http.Header, key string) (*string, error) {
	if values, ok := header[key]; ok {
		switch len(values) {
		case 0:
			return nil, nil
		case 1:
			if parts := strings.SplitN(values[0], "=", 2); len(parts) == 2 {
				if algorithm := strings.ToUpper(parts[0]); algorithm == "MD5" {
					if value := parts[1]; crypto.IsValidBase64EncodedMD5Hash(value) {
						return &value, nil
					}
				}
			}
		}
		return nil, ErrorHeaderInvalid(key)
	}
	return nil, nil
}

func ParseMediaTypeHeader(header http.Header, key string) (*string, error) {
	if values, ok := header[key]; ok {
		switch len(values) {
		case 0:
			return nil, nil
		case 1:
			if value := values[0]; net.IsValidMediaType(value) {
				return &value, nil
			}
		}
		return nil, ErrorHeaderInvalid(key)
	}
	return nil, nil
}

func ParseIntHeader(header http.Header, key string) (*int, error) {
	if values, ok := header[key]; ok {
		switch len(values) {
		case 0:
			return nil, nil
		case 1:
			if value, err := strconv.Atoi(values[0]); err == nil {
				return &value, nil
			}
		}
		return nil, ErrorHeaderInvalid(key)
	}
	return nil, nil
}
