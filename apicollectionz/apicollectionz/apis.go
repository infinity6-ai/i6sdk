package apicollectionz

import (
	"context"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapibatchsummary"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartbatchingest"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartmain"
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
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapievent"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapigeturl"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezget"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezlist"
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
	"reflect"
	"sync"
)

type InvokerFunc func(ctx context.Context, api apiz.Api) (int, error)

type SdkApis struct {
	Invoker InvokerFunc
}

func (a *SdkApis) DszApiLoadResultBatchSummary(ctx context.Context, api *dszapiloadresult.Api[*platformresultz.BatchSummary]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) DszApiLoadResultCatalogSimple(ctx context.Context, api *dszapiloadresult.Api[*relevancefashionresultz.CatalogSimple]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapicustomerportalcreate(ctx context.Context, api *dszapicustomerportalcreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapidatasetlist(ctx context.Context, api *dszapidatasetlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapidatasetsolutions(ctx context.Context, api *dszapidatasetsolutions.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapidomainsolutionget(ctx context.Context, api *dszapidomainsolutionget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapiingestiontokencreate(ctx context.Context, api *dszapiingestiontokencreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapiingestiontokenslist(ctx context.Context, api *dszapiingestiontokenslist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapimodelget(ctx context.Context, api *dszapimodelget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapipaymentautorenewupdate(ctx context.Context, api *dszapipaymentautorenewupdate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapipaymentlist(ctx context.Context, api *dszapipaymentlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapiplanlist(ctx context.Context, api *dszapiplanlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapirecsystokencreate(ctx context.Context, api *dszapirecsystokencreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapirecsystokenslist(ctx context.Context, api *dszapirecsystokenslist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapisignedurl(ctx context.Context, api *dszapisignedurl.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapistorecreate(ctx context.Context, api *dszapistorecreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapistoredisplaycreate(ctx context.Context, api *dszapistoredisplaycreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapistoredisplaylist(ctx context.Context, api *dszapistoredisplaylist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapistorelist(ctx context.Context, api *dszapistorelist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapistripecheckout(ctx context.Context, api *dszapistripecheckout.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapisubscriptionget(ctx context.Context, api *dszapisubscriptionget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapisubscriptionplanupdate(ctx context.Context, api *dszapisubscriptionplanupdate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapitableget(ctx context.Context, api *dszapitableget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapitablelist(ctx context.Context, api *dszapitablelist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Dszapitokendelete(ctx context.Context, api *dszapitokendelete.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzApiEventCatalog(ctx context.Context, api *ingestzapievent.Api[*relevancefashiondataz.Catalog]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzApiEventEvent(ctx context.Context, api *ingestzapievent.Api[*relevancefashiondataz.Event]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) IngestzApiGetUrl(ctx context.Context, api *ingestzapigeturl.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Ingestzapibatchsummary(ctx context.Context, api *ingestzapibatchsummary.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Ingestzapigeturl(ctx context.Context, api *ingestzapigeturl.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Pipezget(ctx context.Context, api *pipezget.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Pipezlist(ctx context.Context, api *pipezlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Pipezstartbatchingest(ctx context.Context, api *pipezstartbatchingest.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Pipezstartmain(ctx context.Context, api *pipezstartmain.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) RoutezSampleFraction(ctx context.Context, api *routezsamplefraction.FractionApi) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) SchedulerActivateBatchIngest(ctx context.Context, api *scheduleractivate.Api[*relevancefashionpipez.BatchIngest]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) SchedulerActivateMain(ctx context.Context, api *scheduleractivate.Api[*relevancefashionpipez.Main]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) SchedulerDeactivateBatchIngest(ctx context.Context, api *schedulerdeactivate.Api[*relevancefashionpipez.BatchIngest]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) SchedulerDeactivateMain(ctx context.Context, api *schedulerdeactivate.Api[*relevancefashionpipez.Main]) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Tfzapimodelmetadata(ctx context.Context, api *tfzapimodelmetadata.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Tfzapimodelpredict(ctx context.Context, api *tfzapimodelpredict.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthdevicecode(ctx context.Context, api *userapiauthdevicecode.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthdevicetoken(ctx context.Context, api *userapiauthdevicetoken.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthexchange(ctx context.Context, api *userapiauthexchange.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthgoogle(ctx context.Context, api *userapiauthgoogle.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthlogout(ctx context.Context, api *userapiauthlogout.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthme(ctx context.Context, api *userapiauthme.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthpasswordreset(ctx context.Context, api *userapiauthpasswordreset.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthpasswordsignin(ctx context.Context, api *userapiauthpasswordsignin.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthpasswordsignup(ctx context.Context, api *userapiauthpasswordsignup.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiauthrefresh(ctx context.Context, api *userapiauthrefresh.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapiconfsettings(ctx context.Context, api *userapiconfsettings.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapidatasettokenaccessattenuate(ctx context.Context, api *userapidatasettokenaccessattenuate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapidatasettokenaccesscreate(ctx context.Context, api *userapidatasettokenaccesscreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapidatasettokenrefresh(ctx context.Context, api *userapidatasettokenrefresh.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapidatasettokenrefreshcreate(ctx context.Context, api *userapidatasettokenrefreshcreate.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapidatasettokenrefreshdelete(ctx context.Context, api *userapidatasettokenrefreshdelete.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapidatasettokenrefreshlist(ctx context.Context, api *userapidatasettokenrefreshlist.Api) (int, error) {
	return a.Invoker(ctx, api)
}

func (a *SdkApis) Userapitokencheck(ctx context.Context, api *userapitokencheck.Api) (int, error) {
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
