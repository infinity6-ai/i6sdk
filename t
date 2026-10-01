[1mdiff --git a/sdkgo/apiz/dszapi/dszapitablelist/dszapitablelist.go b/sdkgo/apiz/dszapi/dszapitablelist/dszapitablelist.go[m
[1mindex c2b66f6..b437efb 100644[m
[1m--- a/sdkgo/apiz/dszapi/dszapitablelist/dszapitablelist.go[m
[1m+++ b/sdkgo/apiz/dszapi/dszapitablelist/dszapitablelist.go[m
[36m@@ -4,6 +4,7 @@[m [mimport ([m
 	"github.com/infinity6-ai/gox/fsz/fsz"[m
 	"github.com/infinity6-ai/gox/routez/apiz"[m
 	"github.com/infinity6-ai/gox/schemaz/schemaz"[m
[32m+[m	[32m"go.code.infinity6.ai/sdkgo/genjsz/jstypez"[m
 	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"[m
 	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakepathz"[m
 )[m
[36m@@ -24,6 +25,10 @@[m [mtype File struct {[m
 	SignedUrl string         `json:"signed_url"`[m
 }[m
 [m
[32m+[m[32mfunc (*File) JSType() string {[m
[32m+[m	[32mreturn jstypez.TypeObject[m
[32m+[m[32m}[m
[32m+[m
 type Resp struct {[m
 	Files []File `json:"files"`[m
 }[m
[36m@@ -38,8 +43,8 @@[m [mfunc (d *Api) ApiSpec() apiz.ApiSpec {[m
 		Id: "dsz-table-list",[m
 		Desc: func() *schemaz.Desc {[m
 			return &schemaz.Desc{[m
[31m-				Name:    "Get Dataset Signed URLs",[m
[31m-				Summary: "Generates presigned URLs for accessing or mutating files in the specified dataset.",[m
[32m+[m				[32mName:    "List Dataset Tables",[m
[32m+[m				[32mSummary: "Lists files and tables in the specified dataset.",[m
 			}[m
 		},[m
 		Method: "GET",[m
[36m@@ -49,51 +54,59 @@[m [mfunc (d *Api) ApiSpec() apiz.ApiSpec {[m
 [m
 func (d *Api) GetDataRefs() *apiz.DataRefs {[m
 	return &apiz.DataRefs{[m
[31m-		// PathParams: &schemaz.Schema{[m
[31m-		// 	Object: func(read bool) map[string]*schemaz.Schema {[m
[31m-		// 		return map[string]*schemaz.Schema{[m
[31m-		// 			"dataset": {[m
[31m-		// 				Desc: &schemaz.Desc{Summary: "Name of the dataset to generate signed URLs for."},[m
[31m-		// 				Str:  lakedsz.ParseStr(&d.Req.Dataset),[m
[31m-		// 			},[m
[31m-		// 		}[m
[31m-		// 	},[m
[31m-		// },[m
[31m-		// QueryParams: nil,[m
[31m-		// ReqHeaders:  nil,[m
[31m-		// ReqBody: &schemaz.Schema{[m
[31m-		// 	Object: func(read bool) map[string]*schemaz.Schema {[m
[31m-		// 		return map[string]*schemaz.Schema{[m
[31m-		// 			"paths": {[m
[31m-		// 				Desc: &schemaz.Desc{Summary: "List of paths and HTTP methods to generate signed URLs for."},[m
[31m-		// 				Array: schemaz.NewArray(schemaz.ArrayOptions[[]Path, Path]{[m
[31m-		// 					Slice: func() *[]Path { return &d.Req.Paths },[m
[31m-		// 					Get: func(idx int, element *Path) *schemaz.Schema {[m
[31m-		// 						return &schemaz.Schema{[m
[31m-		// 							Desc: &schemaz.Desc{Summary: "Path and HTTP method specification for a signed URL."},[m
[31m-		// 							Object: func(read bool) map[string]*schemaz.Schema {[m
[31m-		// 								return map[string]*schemaz.Schema{[m
[31m-		// 									"method": {[m
[31m-		// 										Desc: &schemaz.Desc{Summary: "HTTP method (e.g. GET, PUT) for the signed URL."},[m
[31m-		// 										Raw:  func() any { return &d.Req.Paths[idx].Method },[m
[31m-		// 									},[m
[31m-		// 									"path": {[m
[31m-		// 										Desc: &schemaz.Desc{Summary: "File path within the dataset."},[m
[31m-		// 										Raw:  func() any { return &d.Req.Paths[idx].Path },[m
[31m-		// 									},[m
[31m-		// 								}[m
[31m-		// 							},[m
[31m-		// 						}[m
[31m-		// 					},[m
[31m-		// 				}),[m
[31m-		// 			},[m
[31m-		// 		}[m
[31m-		// 	},[m
[31m-		// },[m
[31m-		// RespHeaders: nil,[m
[31m-		// RespBody: &schemaz.Schema{[m
[31m-		// 	Desc: &schemaz.Desc{Name: "Get URL Response", Summary: "Presigned URLs generated for the requested paths."},[m
[31m-		// 	Raw:  func() any { return &d.Resp },[m
[31m-		// },[m
[32m+[m		[32mPathParams: &schemaz.Schema{[m
[32m+[m			[32mObject: func(read bool) map[string]*schemaz.Schema {[m
[32m+[m				[32mreturn map[string]*schemaz.Schema{[m
[32m+[m					[32m"dataset": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Name of the dataset to list tables from."},[m
[32m+[m						[32mStr:  lakedsz.ParseStr(&d.Req.Dataset),[m
[32m+[m					[32m},[m
[32m+[m				[32m}[m
[32m+[m			[32m},[m
[32m+[m		[32m},[m
[32m+[m		[32mQueryParams: &schemaz.Schema{[m
[32m+[m			[32mObject: func(read bool) map[string]*schemaz.Schema {[m
[32m+[m				[32mreturn map[string]*schemaz.Schema{[m
[32m+[m					[32m"path": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Path to list."},[m
[32m+[m						[32mRaw:  func() any { return &d.Req.Path },[m
[32m+[m					[32m},[m
[32m+[m					[32m"sign": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Whether to return signed URLs."},[m
[32m+[m						[32mRaw:  func() any { return &d.Req.Sign },[m
[32m+[m					[32m},[m
[32m+[m					[32m"recursive": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Whether to list recursively."},[m
[32m+[m						[32mRaw:  func() any { return &d.Req.Recursive },[m
[32m+[m					[32m},[m
[32m+[m					[32m"stat": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Whether to include file stats."},[m
[32m+[m						[32mRaw:  func() any { return &d.Req.Stat },[m
[32m+[m					[32m},[m
[32m+[m					[32m"cursor": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Pagination cursor."},[m
[32m+[m						[32mStr:  schemaz.ParseStr(&d.Req.Cursor),[m
[32m+[m					[32m},[m
[32m+[m					[32m"limist": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "Limit the number of results."},[m
[32m+[m						[32mStr:  schemaz.ParseStrNumber(&d.Req.Limit),[m
[32m+[m					[32m},[m
[32m+[m				[32m}[m
[32m+[m			[32m},[m
[32m+[m		[32m},[m
[32m+[m		[32mReqHeaders: nil,[m
[32m+[m		[32mReqBody:    nil,[m
[32m+[m		[32mRespHeaders: nil,[m
[32m+[m		[32mRespBody: &schemaz.Schema{[m
[32m+[m			[32mDesc: &schemaz.Desc{Summary: "List of files in the dataset."},[m
[32m+[m			[32mObject: func(read bool) map[string]*schemaz.Schema {[m
[32m+[m				[32mreturn map[string]*schemaz.Schema{[m
[32m+[m					[32m"files": {[m
[32m+[m						[32mDesc: &schemaz.Desc{Summary: "File entries"},[m
[32m+[m						[32mRaw:  func() any { return &d.Resp.Files },[m
[32m+[m					[32m},[m
[32m+[m				[32m}[m
[32m+[m			[32m},[m
[32m+[m		[32m},[m
 	}[m
 }[m
[1mdiff --git a/sdkgo/sdkz/sdkz.go b/sdkgo/sdkz/sdkz.go[m
[1mindex 162e588..a29dc6a 100644[m
[1m--- a/sdkgo/sdkz/sdkz.go[m
[1m+++ b/sdkgo/sdkz/sdkz.go[m
[36m@@ -10,6 +10,7 @@[m [mimport ([m
 	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"[m
 	"go.code.infinity6.ai/datasetz/dszapi"[m
 	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapisignedurl"[m
[32m+[m	[32m"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitablelist"[m
 	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzapigeturl"[m
 	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/pipezapi"[m
 )[m
[36m@@ -33,6 +34,7 @@[m [mfunc init() {[m
 	}[m
 	add(&ingestzapigeturl.GetUrlApi{})[m
 	add(&dszapisignedurl.Api{})[m
[32m+[m	[32madd(&dszapitablelist.Api{})[m
 	add(&pipezapi.PipezStartApi{})[m
 	add(&dszapi.DszSolutionDomainSelectApi{})[m
 }[m
