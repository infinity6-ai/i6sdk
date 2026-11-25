//go:generate go run github.com/infinity6-ai/i6sdk/mold apisgen .
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
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiloadresult"
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
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapievent"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapigeturl"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezget"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezlist"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartbatchingest"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartmain"
	"github.com/infinity6-ai/i6sdk/apipipez/scheduleractivate"
	"github.com/infinity6-ai/i6sdk/apipipez/schedulerdeactivate"
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
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/platformresultz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/relevancefashionresultz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
)

type InvokerFunc func(ctx context.Context, api apiz.Api) (int, error)

type SdkApis struct {
	dszapidatasetlist.DszListDatasetsId
	pipezstartbatchingest.PipezStartRelevanceFashionRelevanceFashionBatchIngestId
	pipezstartmain.PipezStartRelevanceFashionRelevanceFashionMainId
	tfzapimodelmetadata.TfzModelMetadataId
	tfzapimodelpredict.TfzModelPredictId

	Invoker InvokerFunc
}

func (a *SdkApis) DszCreateCustomerPortal(ctx context.Context, api *dszapicustomerportalcreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszCreateIngestionToken(ctx context.Context, api *dszapiingestiontokencreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszCreateRecsysToken(ctx context.Context, api *dszapirecsystokencreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszCreateStore(ctx context.Context, api *dszapistorecreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszDatasetSolutions(ctx context.Context, api *dszapidatasetsolutions.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszDeleteToken(ctx context.Context, api *dszapitokendelete.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszDomainGet(ctx context.Context, api *dszapidomainsolutionget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszGetSubscription(ctx context.Context, api *dszapisubscriptionget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

// func (a *SdkApis) DszListDatasets(ctx context.Context, api *dszapidatasetlist.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

func (a *SdkApis) DszListIngestionTokens(ctx context.Context, api *dszapiingestiontokenslist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszListPayments(ctx context.Context, api *dszapipaymentlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszListPlans(ctx context.Context, api *dszapiplanlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszListRecsysTokens(ctx context.Context, api *dszapirecsystokenslist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszListStores(ctx context.Context, api *dszapistorelist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszModelGet(ctx context.Context, api *dszapimodelget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszResultLastBatchSummary(ctx context.Context, api *dszapiloadresult.Api[*platformresultz.BatchSummary]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszResultRelevanceFashionCatalog(ctx context.Context, api *dszapiloadresult.Api[*relevancefashionresultz.CatalogSimple]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszSignedUrl(ctx context.Context, api *dszapisignedurl.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszStoreDisplayCreate(ctx context.Context, api *dszapistoredisplaycreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszStoreDisplayList(ctx context.Context, api *dszapistoredisplaylist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszStripeCheckout(ctx context.Context, api *dszapistripecheckout.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszTableGet(ctx context.Context, api *dszapitableget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszTableList(ctx context.Context, api *dszapitablelist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszUpdatePaymentAutoRenew(ctx context.Context, api *dszapipaymentautorenewupdate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszUpdateSubscriptionPlan(ctx context.Context, api *dszapisubscriptionplanupdate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzEventRelevanceFashionCatalog(ctx context.Context, api *ingestzapievent.Api[*relevancefashiondataz.Catalog]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzEventRelevanceFashionEvent(ctx context.Context, api *ingestzapievent.Api[*relevancefashiondataz.Event]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzGetUrl(ctx context.Context, api *ingestzapigeturl.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzLastBatchSummary(ctx context.Context, api *ingestzapibatchsummary.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) PipezGet(ctx context.Context, api *pipezget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) PipezList(ctx context.Context, api *pipezlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) PipezSchedulerActivateRelevanceFashionRelevanceFashionBatchIngest(ctx context.Context, api *scheduleractivate.Api[*relevancefashionpipez.BatchIngest]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) PipezSchedulerActivateRelevanceFashionRelevanceFashionMain(ctx context.Context, api *scheduleractivate.Api[*relevancefashionpipez.Main]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) PipezSchedulerDeactivateRelevanceFashionRelevanceFashionBatchIngest(ctx context.Context, api *schedulerdeactivate.Api[*relevancefashionpipez.BatchIngest]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) PipezSchedulerDeactivateRelevanceFashionRelevanceFashionMain(ctx context.Context, api *schedulerdeactivate.Api[*relevancefashionpipez.Main]) (int, error) {
	return a.Invoker(ctx, api)
}

// func (a *SdkApis) PipezStartRelevanceFashionRelevanceFashionBatchIngest(ctx context.Context, api *pipezstartbatchingest.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

// func (a *SdkApis) PipezStartRelevanceFashionRelevanceFashionMain(ctx context.Context, api *pipezstartmain.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

func (a *SdkApis) RoutezSampleFraction(ctx context.Context, api *routezsamplefraction.FractionApi) (int, error) {
	return a.Invoker(ctx, api)
}

// func (a *SdkApis) TfzModelMetadata(ctx context.Context, api *tfzapimodelmetadata.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

// func (a *SdkApis) TfzModelPredict(ctx context.Context, api *tfzapimodelpredict.Api) (int, error) {
// 	return a.Invoker(ctx, api)
// }

func (a *SdkApis) UserAuthDeviceCode(ctx context.Context, api *userapiauthdevicecode.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthDeviceToken(ctx context.Context, api *userapiauthdevicetoken.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthExchange(ctx context.Context, api *userapiauthexchange.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthGoogle(ctx context.Context, api *userapiauthgoogle.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthLogout(ctx context.Context, api *userapiauthlogout.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthMe(ctx context.Context, api *userapiauthme.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthPasswordReset(ctx context.Context, api *userapiauthpasswordreset.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthPasswordSignin(ctx context.Context, api *userapiauthpasswordsignin.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthPasswordSignup(ctx context.Context, api *userapiauthpasswordsignup.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserAuthRefresh(ctx context.Context, api *userapiauthrefresh.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserConfSettings(ctx context.Context, api *userapiconfsettings.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserDatasetTokenAccessAttenuate(ctx context.Context, api *userapidatasettokenaccessattenuate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserDatasetTokenAccessCreate(ctx context.Context, api *userapidatasettokenaccesscreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserDatasetTokenRefresh(ctx context.Context, api *userapidatasettokenrefresh.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserDatasetTokenRefreshCreate(ctx context.Context, api *userapidatasettokenrefreshcreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserDatasetTokenRefreshDelete(ctx context.Context, api *userapidatasettokenrefreshdelete.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserDatasetTokenRefreshList(ctx context.Context, api *userapidatasettokenrefreshlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) UserTokenCheck(ctx context.Context, api *userapitokencheck.Api) (int, error) {
	return a.Invoker(ctx, api)
}

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
