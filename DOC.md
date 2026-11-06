# TASK: add missing descriptions to API and SDK function definitions in the Go repo at /home/hex/proj/i6/i6sdk.

Work ONLY on definitions that are missing a description. Anything that already has one is out of scope: do not edit, reword, reformat or "improve" it.

## SCOPE
1. APIs: every ApiSpec() under sdkgo/apiz/**. A missing description is Desc: nil (the other APIs use Desc: func() *schemaz.Desc { return &schemaz.Desc{Name: ..., Summary: ...} }).
2. SDK funcs: every SdkFuncSpec() under sdkgo/productz/** (a missing description is Desc: nil in sdkfuncz.Spec). Before editing, read sdkgo/sdkfuncz/sdkfuncz.go to see the exact type of the Desc field, and use that type. If it is not the same func-returning-*schemaz.Desc shape as the APIs, follow whatever the type requires.
3. Schema fields: every entry in a map[string]*schemaz.Schema (PathParams, QueryParams, headers, ReqBody, RespBody, nested Object and Array elements) that has Field/Object/Array/UrlValue/etc. but NO Desc. This includes shared types the APIs and funcs use: sdkgo/dsz/dsztoken/dsztoken.go, sdkgo/pipez/pipez.go, sdkgo/resultz/**, sdkgo/tabledataz/**, sdkgo/productz/relevance/pipez/**.
   Do not touch genjsz/ or sdkweblegacy/ (not in scope).

### STEP 0 - FIND WHAT IS MISSING (do not trust any old list; the repo may have changed)
API level: grep -rn "Desc: *nil" sdkgo
Field level: write a throwaway Go AST scanner OUTSIDE the repo (use the session scratchpad directory, never inside the repo and never in /tmp) that parses each non-test .go file importing schemaz and prints file:line for every map entry whose value is an untyped composite literal containing one of the keys Field, Object, Array, UrlValue, Str, Int, Bool, Map but NOT the key Desc. Run it over: sdkgo (excluding _test.go files).
Print the full list of findings to the user first, grouped by file.

### STEP 1 - WRITING RULES (follow the existing descriptions as the style guide)
Read 3 or 4 already-described files first for style, for example sdkgo/apiz/dszapi/dszapistorelist/dszapistorelist.go, dszapipaymentlist.go, dszapisubscriptionget.go and tfzapimodelpredict.go.


API and SDK func level (Name AND Summary are both mandatory):
Name: Title Case, 2 to 5 words, verb first, e.g. "List Stores", "Get Subscription", "Create Ingestion Token".
Summary: ONE sentence, ends with a period, starts with a third-person verb ("Lists ..." says what it does and for what scope, e.g. "Lists all stores configured under a specificdataset.". Mention a non-obvious behavior only if it matters (e.g. "The list is empty for datasets without a Stripe customer.").
Add Labels only if the file already has a convention for it (e.g. tfzapimodelpredict  do not invent labels.
Do NOT use Markdown. Use it only if a description genuinely cannot be expressed without it (almost never). If you do use it, say so in your report.
Field level (Summary only, never Name):
Form: Desc: &schemaz.Desc{Summary: "..."}, placed as the first key of the field literds.
The SHORTEST text that is still understandable. Target 2 to 8 words, hard limit about 12. Examples of the target style: "Dataset identifier.", "List of store configurations.", "Dataset to list stores from.".
End with a period. Start with a capital letter. No trailing explanations, no examples, ot obvious (then give it in 1 word, e.g. "Expiration time, UTC.").
Describe the meaning, not the type or the JSON key (bad: "The id string field"; good: "Pipeline run id.").
Same field name in the same domain gets the same wording everywhere (e.g. "dataset" is or the exact wording the nearest already-described file uses for it; reuse it).
Do not guess business meaning you cannot verify. Read the surrounding struct, the handler or the model types it maps to, and the other described files before writing. If a meaning is truly unclear, write the
most neutral true statement and add the field to a "needs human review" list in your fina


### STEP 2 - ASCII ONLY (strict)
Every character you add must be plain 7-bit ASCII. Forbidden: em dash, en dash, curly/smart quotes, ellipsis character, arrows, bullets, non-breaking spaces, accents, emojis, any unicode. Use "-" , straight
quotes only where Go requires them, and "..." as three dots if ever needed.
Inside Go strings, do not use backticks or double quotes; rephrase instead. No escape sequences.
After editing, run: grep -nP '[^\x00-\x7F]' <every file you edited>
  It must print nothing for the lines you added. If a pre-existing non-ASCII character is found on a line you did not add, leave it alone and mention it in the report.


### STEP 3 - EDITING RULES
Minimal diffs. Change only the Desc lines (and the Desc: nil line) of the definitiofactor, rename, move code, extract helpers, or change any function, tag, logic or import.
Do NOT run gofmt/goimports on whole files (some files have pre-existing misaligned lines such as Id:     "..." that must stay unchanged). Match the indentation (tabs) of the neighbouring lines by hand.
If you replace Desc: nil for an API, use exactly the same multi-line shape as the des` import only if the file does not already import"github.com/infinity6-ai/gox/schemaz/schemaz".
Do not touch the Go SDK vs web SDK naming; this task changes no names, only description
Do not edit test files unless a test asserts on a description you changed. Existing tests that assert Markdown: "" must keep passing.
No git commit, no push. Leave changes in the working tree.

### STEP 4 - VERIFY
go build ./... and go vet on the touched packages.
Run unit tests only: go test ./... -run 'TestUnit' (do not run Remote or Manual tests). output; if a failure was pre-existing, confirm that by checking on a clean stash/checkoutbefore claiming it.
Re-run the Step 0 scanner and grep. Expected result: zero Desc: nil for in-scope APIslds without Desc, zero non-ASCII characters added.
git diff --stat and confirm only description lines changed (git diff should contain no changes other than added/changed Desc lines and possibly one import).


### WORKFLOW AND REPORTING
Final report: counts (APIs fixed, funcs fixed, fields fixed), the "needs human review" list, any anomalies (non-ASCII found, pre-existing test failures), and nothing else. Do not invent work beyond the scope above