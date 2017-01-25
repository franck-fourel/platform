package data

/* CHECKLIST
 * [x] Uses interfaces as appropriate
 * [x] Private package variables use underscore prefix
 * [x] All parameters validated
 * [x] All errors handled
 * [x] Reviewed for concurrency safety
 * [x] Code complete
 * [x] Full test coverage
 */

import (
	"strconv"

	"github.com/tidepool-org/platform/app"
)

type Deduplicator interface {
	Name() string

	RegisterDataset() error

	AddDatasetData(datasetData []Datum) error
	DeduplicateDataset() error

	DeleteDataset() error
}

type DeduplicatorDescriptor struct {
	Name string `bson:"name,omitempty"`
	Hash string `bson:"hash,omitempty"`
}

func NewDeduplicatorDescriptor() *DeduplicatorDescriptor {
	return &DeduplicatorDescriptor{}
}

func (d *DeduplicatorDescriptor) IsRegisteredWithAnyDeduplicator() bool {
	return d.Name != ""
}

func (d *DeduplicatorDescriptor) IsRegisteredWithNamedDeduplicator(name string) bool {
	return d.Name == name
}

func (d *DeduplicatorDescriptor) RegisterWithNamedDeduplicator(name string) error {
	if d.Name != "" {
		return app.Errorf("data", "deduplicator descriptor already registered with %s", strconv.Quote(d.Name))
	}

	d.Name = name
	return nil
}