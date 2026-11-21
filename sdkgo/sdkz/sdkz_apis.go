package sdkz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
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
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstart"
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
)

var apis map[string]apiz.Api

func init() {
	add := func(as ...apiz.Api) {
		for _, a := range as {
			id := a.ApiSpec().Id
			old, ok := apis[id]
			if ok {
				panic(fmt.Sprintf("api id collision %s: old %T, new %T", id, old, a))
			}
			apis[id] = a
		}
	}

	samples := routezsamplefraction.Services()
	apis = make(map[string]apiz.Api, len(samples))
	for _, svc := range routezsamplefraction.Services() {
		add(svc.Api())
	}
	add(&ingestzapigeturl.Api{})

	// disabled in favor of ingestzapievent
	// add(ingestzapistream.AllApis()...)

	add(ingestzapievent.AllApis()...)
	// add(&ingestzapibatchsummary.Api{})
	add(&dszapisignedurl.Api{})
	add(&dszapitableget.Api{})
	add(&dszapitablelist.Api{})
	add(dszapiloadresult.AllApis()...)
	add(pipezstart.AllApis()...)
	add(scheduleractivate.AllApis()...)
	add(schedulerdeactivate.AllApis()...)
	add(&pipezlist.Api{})
	add(&pipezget.Api{})
	add(&dszapistoredisplaycreate.Api{})
	add(&dszapistoredisplaylist.Api{})
	add(&dszapistorecreate.Api{})
	add(&dszapistorelist.Api{})
	add(&dszapidatasetlist.Api{})
	add(&dszapidatasetsolutions.Api{})
	add(&dszapidomainsolutionget.Api{})
	add(&dszapiingestiontokenslist.Api{})
	add(&dszapiingestiontokencreate.Api{})
	add(&dszapirecsystokenslist.Api{})
	add(&dszapirecsystokencreate.Api{})
	add(&dszapitokendelete.Api{})
	add(&dszapimodelget.Api{})
	add(&tfzapimodelmetadata.Api{})
	add(&tfzapimodelpredict.Api{})
	add(&dszapistripecheckout.Api{})
	add(&dszapicustomerportalcreate.Api{})
	add(&dszapisubscriptionget.Api{})
	add(&dszapipaymentautorenewupdate.Api{})
	add(&dszapipaymentlist.Api{})
	add(&dszapiplanlist.Api{})
	add(&dszapisubscriptionplanupdate.Api{})
	add(&userapiauthdevicecode.Api{})
	add(&userapiauthdevicetoken.Api{})
	add(&userapiauthexchange.Api{})
	add(&userapiauthgoogle.Api{})
	add(&userapiauthlogout.Api{})
	add(&userapiauthme.Api{})
	add(&userapiauthpasswordreset.Api{})
	add(&userapiauthpasswordsignin.Api{})
	add(&userapiauthpasswordsignup.Api{})
	add(&userapiauthrefresh.Api{})
	add(&userapiconfsettings.Api{})
	add(&userapidatasettokenaccessattenuate.Api{})
	add(&userapidatasettokenaccesscreate.Api{})
	add(&userapidatasettokenrefresh.Api{})
	add(&userapidatasettokenrefreshcreate.Api{})
	add(&userapidatasettokenrefreshdelete.Api{})
	add(&userapidatasettokenrefreshlist.Api{})
	add(&userapitokencheck.Api{})
}

func Apis() []apiz.Api {
	return slices.Collect(maps.Values(apis))
}

func Api(id string) apiz.Api {
	ret, ok := apis[id]
	checker.True(ok, "api not found: %s", id)
	return ret
}
