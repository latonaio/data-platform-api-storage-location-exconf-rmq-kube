package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-storage-location-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-storage-location-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-storage-location-exconf-rmq-kube/database"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.StorageLocation {
	businessPartner := *input.StorageLocation.BusinessPartner
	plant := *input.StorageLocation.Plant
	storageLocation := *input.StorageLocation.StorageLocation
	notKeyExistence := make([]dpfm_api_output_formatter.StorageLocation, 0, 1)
	KeyExistence := make([]dpfm_api_output_formatter.StorageLocation, 0, 1)

	existData := &dpfm_api_output_formatter.StorageLocation{
		BusinessPartner: businessPartner,
		Plant:           plant,
		StorageLocation: storageLocation,
		ExistenceConf:   false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confStorageLocation(businessPartner, plant, storageLocation) {
			notKeyExistence = append(
				notKeyExistence,
				dpfm_api_output_formatter.StorageLocation{businessPartner, plant, storageLocation, false},
			)
			return
		}
		KeyExistence = append(KeyExistence, dpfm_api_output_formatter.StorageLocation{businessPartner, plant, storageLocation, true})
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confStorageLocation(businessPartner int, plant string, storageLocation string) bool {
	rows, err := e.db.Query(
		`SELECT StorageLocation 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_plant_storage_location_data 
		WHERE (BusinessPartner, plant, storageLocation) = (?, ?, ?);`, businessPartner, plant, storageLocation,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var businessPartner int
		var plant string
		var storageLocation string
		err := rows.Scan(&storageLocation)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if businessPartner == businessPartner {
			return true
		}
		if plant == plant {
			return true
		}
	}
	return false
}
