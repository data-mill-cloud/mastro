# Mastro
## Data Catalogue
Data providers can describe and publish data using a shared definition format.
Consequently, data definitions can be crawled from networked and distributed file systems, as well as directly published to a common endpoint.

### Catalogue API
A Catalogue service endpoint implements the following interface:

```go
type CatalogueService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	UpsertAssets(assets *[]Asset) (*[]Asset, *resterrors.RestErr)
	GetAssetByID(assetID string) (*Asset, *resterrors.RestErr)
	GetAssetByName(name string) (*Asset, *resterrors.RestErr)
	SearchAssetsByTags(tags []string, limit int, page int) (*Paginated[Asset], *resterrors.RestErr)
	Search(query string, limit int, page int) (*Paginated[Asset], *resterrors.RestErr)
	ListAllAssets(limit int, page int) (*Paginated[Asset], *resterrors.RestErr)
}
```

This can be easily mapped to a DAO backend:
```go
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Upsert(asset *Asset) error
	GetById(id string) (*Asset, error)
	GetByName(id string) (*Asset, error)
	SearchAssetsByTags(tags []string, limit int, page int) (*Paginated[Asset], error)
	ListAllAssets(limit int, page int) (*Paginated[Asset], error)
	Search(query string, limit int, page int) (*Paginated[Asset], error)
	CloseConnection()
}
```

Have a look at `catalogue/daos/*` for example implementations.

This is translated to the following endpoint:

| Verb        | Endpoint                | Maps to                                                         |
|-------------|-------------------------|-----------------------------------------------------------------|
| **GET**     | /healthcheck/asset      | github.com/data-mill-cloud/mastro/catalogue.Ping                |
| ~~**GET**~~ | ~~/asset/id/:asset_id~~ | ~~github.com/data-mill-cloud/mastro/catalogue.GetAssetByID~~    |
| **GET**     | /asset/name/:asset_name | github.com/data-mill-cloud/mastro/catalogue.GetAssetByName      |
| **PUT**     | /asset/                 | github.com/data-mill-cloud/mastro/catalogue.UpsertAsset         |
| **PUT**     | /assets/                | github.com/data-mill-cloud/mastro/catalogue.BulkUpsert          |
| **POST**    | /assets/tags            | github.com/data-mill-cloud/mastro/catalogue.SearchAssetsByTags  |
| **POST**    | /assets/search          | github.com/data-mill-cloud/mastro/catalogue.Search              |
| ~~**GET**~~ | ~~/assets/~~            | ~~github.com/data-mill-cloud/mastro/catalogue.ListAllAssets~~   | 

Those crossed out are meant for testing purposes and will be removed in the following releases.

### Examples

We provide a few examples below:

List all - *GET* on `localhost:8085/assets` with empty body, has result:
```json
{
    "message": "Error while retrieving asset :: mongo: no documents in result",
    "status": 404,
    "error": "not_found"
}
```

Upsert - *PUT* on `localhost:8085/asset` with body:
```json
{"last-discovered-at" : "2021-03-22T21:19:39.634Z", "published-on" : "0001-01-01T00:00:00.000Z", "name" : "example_featureset", "description" : "my first featureset pushed to the catalogue", "depends-on" : ["table.mydb.mytable"], "type" : "featureset"}
```

Bulk upsert - *PUT* on `localhost:8085/assets` with body:
```json
[
	{"last-discovered-at" : "2021-03-22T21:19:39.634Z", "published-on" : "0001-01-01T00:00:00.000Z", "name" : "example_featureset", "description" : "my first featureset pushed to the catalogue", "depends-on" : ["table.mydb.mytable"], "type" : "featureset", "tags" : ["featureset"]},
    {"last-discovered-at" : "2021-03-22T21:19:39.634Z", "published-on" : "0001-01-01T00:00:00.000Z", "name" : "example_featureset", "description" : "my first featureset pushed to the catalogue", "depends-on" : ["table.mydb.mytable"], "type" : "featureset", "tags" : ["featureset"]}    
]
```

GetByName - *GET* on `localhost:8085/asset/example_featureset` has now result:
```json
{
    "last-discovered-at": "2021-03-23T13:52:43.787Z",
    "published-on": "0001-01-01T00:00:00Z",
    "name": "example_featureset",
    "description": "my first featureset pushed to the catalogue",
    "depends-on": [
        "table.mydb.mytable"
    ],
    "type": "featureset",
	"tags": [
        "featureset"
    ]
}
```

SearchAssetsByTags - *POST* on `localhost:8085/assets/tags` passing a Json body of kind:
```json
{
    "tags" : ["something"],
	"limit": 4,
    "page": 1
}
```

returns an HTTP error status with a Json body of kind:
```json
{
    "message": "Error while retrieving assets using filter :: empty result set",
    "status": 404,
    "error": "not_found"
}
```

while with body:
```json
{
    "tags" : ["featureset"],
	"limit": 4,
    "page": 1
}
```

we get a paginated list of all assets having the provided tags:
```json
{
    "data": [
        {
            "last-discovered-at": "2021-03-23T13:52:43.787Z",
            "published-on": "0001-01-01T00:00:00Z",
            "name": "example_featureset",
            "description": "my first featureset pushed to the catalogue",
            "depends-on": [
                "table.mydb.mytable",
                "example_featureset1",
                "otherdep2"
            ],
            "type": "featureset",
            "labels": {
                "environment": "test",
                "instance": "2020_08",
                "project": "example"
            },
            "tags": [
                "featureset",
                "example"
            ],
            "versions": {}
        },
        {
            "last-discovered-at": "2021-03-23T13:52:43.788Z",
            "published-on": "0001-01-01T00:00:00Z",
            "name": "example_featureset1",
            "description": "my second featureset pushed to the catalogue",
            "depends-on": [
                "table.mydb.mytable",
                "example_featureset2",
                "otherdep2"
            ],
            "type": "featureset",
            "labels": {
                "environment": "test",
                "instance": "2020_08",
                "project": "example"
            },
            "tags": [
                "featureset",
                "example"
            ],
            "versions": {}
        },
    ],
    "pagination": {
        "total": 15,
        "page": 1,
        "perPage": 4,
        "prev": 0,
        "next": 2,
        "totalPage": 4
    }
}
```