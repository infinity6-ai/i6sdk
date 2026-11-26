//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest apisgen .
package apicollectionz

import (
	"context"
	"reflect"
	"sync"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapicustomerportalcreate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetlist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetsolutions"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidomainsolutionget"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiingestiontokencreate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiingestiontokenslist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiloadresultbatchsummary"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiloadresultcatalog"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapimodelget"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapipaymentautorenewupdate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapipaymentlist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiplanlist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapirecsystokencreate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapirecsystokenslist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisignedurl"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorecreate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistoredisplaycreate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistoredisplaylist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorelist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistripecheckout"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisubscriptionget"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisubscriptionplanupdate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitableget"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitablelist"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitokendelete"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapibatchsummary"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapieventcatalog"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapieventevent"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapigeturl"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapistreamcatalog"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapistreamevent"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezget"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezlist"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartbatchingest"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartmain"
	"github.com/infinity6-ai/i6sdk/apipipez/scheduleractivatebatchingest"
	"github.com/infinity6-ai/i6sdk/apipipez/scheduleractivatemain"
	"github.com/infinity6-ai/i6sdk/apipipez/schedulerdeactivatebatchingest"
	"github.com/infinity6-ai/i6sdk/apipipez/schedulerdeactivatemain"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelmetadata"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelpredict"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicecode"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicetoken"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthexchange"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthgoogle"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthme"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordreset"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordsignin"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordsignup"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthrefresh"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiconfsettings"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccessattenuate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccesscreate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefresh"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshcreate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshdelete"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshlist"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
)

type InvokerFunc func(ctx context.Context, api apiz.Api) (int, error)

type SdkApis struct {
	dszapicustomerportalcreate.DszCreateCustomerPortalId
	dszapidatasetlist.DszListDatasetsId
	dszapidatasetsolutions.DszDatasetSolutionsId
	dszapidomainsolutionget.DszDomainGetId
	dszapiingestiontokencreate.DszCreateIngestionTokenId
	dszapiingestiontokenslist.DszListIngestionTokensId
	dszapiloadresultbatchsummary.DszResultLastBatchSummaryId
	dszapiloadresultcatalog.DszResultRelevanceFashionCatalogId
	dszapimodelget.DszModelGetId
	dszapipaymentautorenewupdate.DszUpdatePaymentAutoRenewId
	dszapipaymentlist.DszListPaymentsId
	dszapiplanlist.DszListPlansId
	dszapirecsystokencreate.DszCreateRecsysTokenId
	dszapirecsystokenslist.DszListRecsysTokensId
	dszapisignedurl.DszSignedUrlId
	dszapistorecreate.DszCreateStoreId
	dszapistoredisplaycreate.DszStoreDisplayCreateId
	dszapistoredisplaylist.DszStoreDisplayListId
	dszapistorelist.DszListStoresId
	dszapistripecheckout.DszStripeCheckoutId
	dszapisubscriptionget.DszGetSubscriptionId
	dszapisubscriptionplanupdate.DszUpdateSubscriptionPlanId
	dszapitableget.DszTableGetId
	dszapitablelist.DszTableListId
	dszapitokendelete.DszDeleteTokenId
	ingestzapibatchsummary.IngestzLastBatchSummaryId
	ingestzapieventcatalog.IngestzEventRelevanceFashionCatalogId
	ingestzapieventevent.IngestzEventRelevanceFashionEventId
	ingestzapigeturl.IngestzGetUrlId
	ingestzapistreamcatalog.IngestzStreamRelevanceFashionCatalogId
	ingestzapistreamevent.IngestzStreamRelevanceFashionEventId
	pipezget.PipezGetId
	pipezlist.PipezListId
	pipezstartbatchingest.PipezStartRelevanceFashionRelevanceFashionBatchIngestId
	pipezstartmain.PipezStartRelevanceFashionRelevanceFashionMainId
	scheduleractivatebatchingest.PipezSchedulerActivateRelevanceFashionRelevanceFashionBatchIngestId
	scheduleractivatemain.PipezSchedulerActivateRelevanceFashionRelevanceFashionMainId
	schedulerdeactivatebatchingest.PipezSchedulerDeactivateRelevanceFashionRelevanceFashionBatchIngestId
	schedulerdeactivatemain.PipezSchedulerDeactivateRelevanceFashionRelevanceFashionMainId
	tfzapimodelmetadata.TfzModelMetadataId
	tfzapimodelpredict.TfzModelPredictId
	userapiauthdevicecode.UserAuthDeviceCodeId
	userapiauthdevicetoken.UserAuthDeviceTokenId
	userapiauthexchange.UserAuthExchangeId
	userapiauthgoogle.UserAuthGoogleId
	userapiauthlogout.UserAuthLogoutId
	userapiauthme.UserAuthMeId
	userapiauthpasswordreset.UserAuthPasswordResetId
	userapiauthpasswordsignin.UserAuthPasswordSigninId
	userapiauthpasswordsignup.UserAuthPasswordSignupId
	userapiauthrefresh.UserAuthRefreshId
	userapiconfsettings.UserConfSettingsId
	userapidatasettokenaccessattenuate.UserDatasetTokenAccessAttenuateId
	userapidatasettokenaccesscreate.UserDatasetTokenAccessCreateId
	userapidatasettokenrefresh.UserDatasetTokenRefreshId
	userapidatasettokenrefreshcreate.UserDatasetTokenRefreshCreateId
	userapidatasettokenrefreshdelete.UserDatasetTokenRefreshDeleteId
	userapidatasettokenrefreshlist.UserDatasetTokenRefreshListId
	userapitokencheck.UserTokenCheckId

	Invoker InvokerFunc
}

func (a *SdkApis) RoutezSampleFraction(ctx context.Context, api *routezsamplefraction.FractionApi) (int, error) {
	return a.Invoker(ctx, api)
}

// func (a *SdkApis) TfzModelMetadata(ctx context.Context, api *tfzapimodelmetadata.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

// func (a *SdkApis) TfzModelPredict(ctx context.Context, api *tfzapimodelpredict.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

func (a *SdkApis) List() []apiz.Api {
	var list []apiz.Api
	typ := reflect.TypeOf(a)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.Name == "List" || m.Name == "Map" {
			continue
		}
		if m.Type.NumIn() != 3 {
			continue
		}
		apiType := m.Type.In(2)
		if apiType.Kind() == reflect.Ptr {
			apiType = apiType.Elem()
		}
		apiVal, ok := reflect.New(apiType).Interface().(apiz.Api)
		if ok {
			list = append(list, apiVal)
		}
	}
	return list
}

var defaultApis = &SdkApis{}
var cachedMap map[string]apiz.Api
var once sync.Once

func (a *SdkApis) Map() map[string]apiz.Api {
	once.Do(func() {
		list := a.List()
		cachedMap = make(map[string]apiz.Api, len(list))
		for _, item := range list {
			cachedMap[item.ApiSpec().Id] = item
		}
	})
	return cachedMap
}

func ApisList() []apiz.Api {
	return defaultApis.List()
}

func Apis() []apiz.Api {
	return defaultApis.List()
}

func Api(id string) apiz.Api {
	ret, ok := defaultApis.Map()[id]
	if !ok {
		panic("api not found: " + id)
	}
	return ret
}
