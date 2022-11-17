package dpfm_api_input_reader

import (
	"data-platform-api-storage-location-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToStorageLocation() *requests.StorageLocation {
	data := sdc.StorageLocation
	return &requests.StorageLocation{
		BusinessPartner: data.BusinessPartner,
		Plant:           data.Plant,
		StorageLocation: data.StorageLocation,
	}
}
