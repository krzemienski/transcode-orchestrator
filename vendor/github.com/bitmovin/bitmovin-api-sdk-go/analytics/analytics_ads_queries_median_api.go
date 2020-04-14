package analytics
import (
    "github.com/bitmovin/bitmovin-api-sdk-go/common"
    
    "github.com/bitmovin/bitmovin-api-sdk-go/model"
)

type AnalyticsAdsQueriesMedianApi struct {
    apiClient *common.ApiClient
}

func NewAnalyticsAdsQueriesMedianApi(configs ...func(*common.ApiClient)) (*AnalyticsAdsQueriesMedianApi, error) {
	apiClient, err := common.NewApiClient(configs...)
	if err != nil {
		return nil, err
	}

    api := &AnalyticsAdsQueriesMedianApi{apiClient: apiClient}


	if err != nil {
		return nil, err
	}

	return api, nil
}

func (api *AnalyticsAdsQueriesMedianApi) Create(adAnalyticsMedianQueryRequest model.AdAnalyticsMedianQueryRequest) (*model.AnalyticsResponse, error) {
    reqParams := func(params *common.RequestParams) {
    }

    var responseModel *model.AnalyticsResponse
    err := api.apiClient.Post("/analytics/ads/queries/median", &adAnalyticsMedianQueryRequest, &responseModel, reqParams)
    return responseModel, err
}

