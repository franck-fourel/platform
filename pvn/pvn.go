package main

import (
	"encoding/json"
	"log"

	"github.com/tidepool-org/platform/pvn/data"
	"github.com/tidepool-org/platform/pvn/data/context"
	"github.com/tidepool-org/platform/pvn/data/normalizer"
	"github.com/tidepool-org/platform/pvn/data/parser"
	"github.com/tidepool-org/platform/pvn/data/types"
	"github.com/tidepool-org/platform/pvn/data/validator"
)

var rawJSON = `
[
  {
    "type": "sample",
    "boolean": false,
    "integer": 15,
    "float": 1.2345,
    "string": "a string",
    "stringArray": [
      "one",
      "two",
      "three"
    ],
    "object": {
      "one": 1,
      "two": "two",
      "three": {
        "a": "apple"
      }
    },
    "objectArray": [
      {
        "alpha": "a"
      },
      {
        "bravo": "b"
      }
    ],
    "interface": "yes",
    "interfaceArray": [
      "alpha",
      {
        "alpha": "a"
      },
      {
        "bravo": "b"
      },
      -999
    ]
  },
  {
    "type": "sample",
    "subType": "sub",
    "boolean": true,
    "integer": 15,
    "float": 1.2345,
    "string": "a string",
    "stringArray": [
      "one",
      "two",
      "three"
    ],
    "object": {
      "one": 1,
      "two": "two",
      "three": {
        "a": "apple"
      }
    },
    "objectArray": [
      {
        "alpha": "a"
      },
      {
        "bravo": "b"
      }
    ],
    "interface": "yes",
    "interfaceArray": [
      "alpha",
      {
        "alpha": "a"
      },
      {
        "bravo": "b"
      },
      -999
    ],
    "innerStruct": {
      "one": "yep, a string",
      "twos": [
        "2",
        "22",
        "222"
      ]
    },
    "innerStructArray": [
      {
        "one": "what, more strings?",
        "twos": [
        ]
      },
      {
        "one": "1"
      }
    ]
  },
  {
    "type": "sample",
    "subType": "sub",
    "boolean": "true",
    "integer": "15",
    "float": "1.2345",
    "string": 45,
    "stringArray": [
      1,
      2,
      3
    ],
    "object": 14,
    "objectArray": [
      "alpha",
      "beta"
    ],
    "innerStruct": {
      "one": 1,
      "twos": 2
    },
    "innerStructArray": [
      {
        "one": "what, more strings?",
        "twos": [
          1,
          false
        ]
      },
      {
        "one": false
      }
    ]
  },
  {
    "type": "sample",
    "subType": "sub"
  }
]
`

func main() {

	log.Printf("==> Loading JSON...")

	rawObjects := []interface{}{}
	if err := json.Unmarshal([]byte(rawJSON), &rawObjects); err != nil {
		log.Fatal("ERROR: Failure parsing JSON: ", err.Error())
	}

	log.Printf("==> Loaded %d objects!", len(rawObjects))

	context := context.NewStandard()

	log.Printf("==> Parsing objects...")

	standardArrayParser := parser.NewStandardArray(context, &rawObjects)

	parsedObjects := []data.Datum{}
	for index := range *standardArrayParser.Array() {
		log.Printf("--- Parsing object #%d", index)
		standardObjectParser := standardArrayParser.NewChildObjectParser(index)
		log.Printf("Raw: %v", *standardObjectParser.Object())
		parsedObject := types.Parse(standardObjectParser)
		log.Printf("Parsed: %s", stringifyObject(parsedObject))
		parsedObjects = append(parsedObjects, parsedObject)
	}

	log.Printf("==> Parsed objects!")

	log.Printf("==> Validating objects...")

	standardValidator := validator.NewStandard(context)

	for index, parsedObject := range parsedObjects {
		log.Printf("--- Validating object #%d", index)
		parsedObject.Validate(standardValidator.NewChildValidator(index))
	}

	log.Printf("==> Validated objects!")

	log.Printf("==> Normalizing objects...")

	standardNormalizer := normalizer.NewStandard(context)

	for index, parsedObject := range parsedObjects {
		log.Printf("--- Normalizing object #%d", index)
		parsedObject.Normalize(standardNormalizer.NewChildNormalizer(index))
	}

	log.Printf("==> Data added during normalization")

	for _, datum := range standardNormalizer.Data() {
		log.Printf("Added: %s", stringifyObject(datum))
	}

	log.Printf("==> Normalized objects!")

	if context.HasErrors() {
		log.Printf("There were %d errors:", len(context.GetErrors()))
		for _, err := range context.GetErrors() {
			log.Print(stringifyObject(err))
		}
	} else {
		log.Print("There were no errors.")
	}

	log.Print("Done!")
}

func stringifyObject(object interface{}) string {
	bytes, err := json.Marshal(object)
	if err != nil {
		return "ERROR: Failure stringifying object"
	}
	return string(bytes)
}
