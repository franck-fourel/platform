package deduplicator

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
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/data/types/upload"
	"github.com/tidepool-org/platform/log"
)

type truncateFactory struct {
	*BaseFactory
}

type truncateDeduplicator struct {
	*BaseDeduplicator
}

const _TruncateDeduplicatorName = "truncate"

var _TruncateExpectedDeviceManufacturers = []string{"Animas"}

func NewTruncateFactory() (Factory, error) {
	baseFactory, err := NewBaseFactory(_TruncateDeduplicatorName)
	if err != nil {
		return nil, err
	}

	factory := &truncateFactory{
		BaseFactory: baseFactory,
	}
	factory.Factory = factory

	return factory, nil
}

func (t *truncateFactory) CanDeduplicateDataset(dataset *upload.Upload) (bool, error) {
	if can, err := t.BaseFactory.CanDeduplicateDataset(dataset); err != nil || !can {
		return can, err
	}

	if dataset.DeviceID == nil {
		return false, nil
	}
	if *dataset.DeviceID == "" {
		return false, nil
	}
	if dataset.DeviceManufacturers == nil {
		return false, nil
	}
	if !app.StringsContainsAnyStrings(*dataset.DeviceManufacturers, _TruncateExpectedDeviceManufacturers) {
		return false, nil
	}

	return true, nil
}

func (t *truncateFactory) NewDeduplicatorForDataset(logger log.Logger, dataStoreSession store.Session, dataset *upload.Upload) (data.Deduplicator, error) {
	baseDeduplicator, err := NewBaseDeduplicator(t.name, logger, dataStoreSession, dataset)
	if err != nil {
		return nil, err
	}

	if dataset.DeviceID == nil {
		return nil, app.Error("deduplicator", "dataset device id is missing")
	}
	if *dataset.DeviceID == "" {
		return nil, app.Error("deduplicator", "dataset device id is empty")
	}
	if dataset.DeviceManufacturers == nil {
		return nil, app.Error("deduplicator", "dataset device manufacturers is missing")
	}
	if !app.StringsContainsAnyStrings(*dataset.DeviceManufacturers, _TruncateExpectedDeviceManufacturers) {
		return nil, app.Error("deduplicator", "dataset device manufacturers does not contain expected device manufacturers")
	}

	return &truncateDeduplicator{
		BaseDeduplicator: baseDeduplicator,
	}, nil
}

func (t *truncateDeduplicator) DeduplicateDataset() error {
	// TODO: Technically, ActivateDatasetData could succeed, but DeleteOtherDatasetData fail. This would
	// result in duplicate (and possible incorrect) data. Is there a way to resolve this? Would be nice to have transactions.

	if err := t.BaseDeduplicator.DeduplicateDataset(); err != nil {
		return err
	}

	if err := t.dataStoreSession.DeleteOtherDatasetData(t.dataset); err != nil {
		return app.ExtErrorf(err, "deduplicator", "unable to remove all other data except dataset with id %s", strconv.Quote(t.dataset.UploadID))
	}

	return nil
}