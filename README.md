# data-platform-api-storage-location-exconf-rmq-kube
data-platform-api-storage-location-exconf-rmq-kube は、データ連携基盤において、API で 保管場所の存在性チェックを行うためのマイクロサービスです。

## 動作環境
・ OS: LinuxOS  
・ CPU: ARM/AMD/Intel  

## 存在確認先テーブル名
以下のsqlファイルに対して、保管場所の存在確認が行われます。

* data-platform-plant-sql-strage-location-data.sql（データ連携基盤 プラント - 保管場所データ）

## caller.go による存在性確認
Input で取得されたファイルに基づいて、caller.go で、 API がコールされます。
caller.go の 以下の箇所が、指定された API をコールするソースコードです。

```
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
```

## Input
data-platform-api-storage-location-exconf-rmq-kube では、以下のInputファイルをRabbitMQからJSON形式で受け取ります。  

```
{
	"connection_key": "request",
	"result": true,
	"redis_key": "abcdefg",
	"api_status_code": 200,
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"business_partner": 201,
	"filepath": "/var/lib/aion/Data/rededge_sdc/abcdef.json",
	"service_label": "PRODUCT_MASTER",
	"StorageLocation": {
		"BusinessPartner": 101,
		"Plant": "AB01",
		"StorageLocation": "0001"
	},
	"api_schema": "DPFMProductMasterCreates",
	"accepter": ["All"],
	"order_id": null,
	"deleted": false
}
```

## Output
data-platform-api-storage-location-exconf-rmq-kube では、[golang-logging-library-for-data-platform](https://github.com/latonaio/golang-logging-library-for-data-platform) により、Output として、RabbitMQ へのメッセージを JSON 形式で出力します。保管場所の対象値が存在する場合 true、存在しない場合 false、を返します。"cursor" ～ "time"は、golang-logging-library-for-data-platform による 定型フォーマットの出力結果です。

```
{
	"connection_key": "request",
	"result": true,
	"redis_key": "abcdefg",
	"filepath": "/var/lib/aion/Data/rededge_sdc/abcdef.json",
	"api_status_code": 200,
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"business_partner": 201,
	"service_label": "PRODUCT_MASTER",
	"StorageLocation": {
		"BusinessPartner": 201,
		"Plant": "AB01",
		"StorageLocation": "0001",
		"ExistenceConf": false
	},
	"api_schema": "DPFMProductMasterCreates",
	"accepter": ["All"],
	"order_id": null,
	"deleted": false
}
```