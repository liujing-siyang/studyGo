package painkiller
//go:generate stringer -type=Pill
type Pill int

const (
    Placebo Pill = iota
    Aspirin
    Ibuprofen
    Paracetamol
    Acetaminophen = Paracetamol
)
// "stack trace:
// \nError 1525: Incorrect TIMESTAMP value: '2023-00-01'
// \nGet data from mysql AllocateCommission Staff fail
// \nform_push/internal/data.(*formAllocateCommissionRepo).GetPushFiledStaff
// \n\tD:/jy7188/workproject/form_push/internal/data/yy_allocate_commission.go:91
// \nform_push/internal/biz.(*FormSet).PushDataOfCommission
// \n\tD:/jy7188/workproject/form_push/internal/biz/yy_allocate_commission.go:104
// \nform_push/internal/service.(*FormPushService).AllocateCommission
// \n\tD:/jy7188/workproject/form_push/internal/service/formpush.go:174
// \nform_push/api/form_push/v1._FormPush_AllocateCommission0_HTTP_Handler.func1.1
// \n\tD:/jy7188/workproject/form_push/api/form_push/v1/form_push_http.pb.go:122
// \nform_push/internal/middleware/urlcheck.PingUrl.func1.1
// \n\tD:/jy7188/workproject/form_push/internal/middleware/urlcheck/pingurl.go:28
// \ngithub.com/go-kratos/kratos/v2/middleware/selector.selector.func1.1
// \n\tD:/jy7188/goproject/pkg/mod/github.com/go-kratos/kratos/v2@v2.5.3/middleware/selector/selector.go:125
// \ngithub.com/go-kratos/kratos/v2/middleware/recovery.Recovery.func2.1
// \n\tD:/jy7188/goproject/pkg/mod/github.com/go-kratos/kratos/v2@v2.5.3/middleware/recovery/recovery.go:60
// \nform_push/api/form_push/v1._FormPush_AllocateCommission0_HTTP_Handler.func1
// \n\tD:/jy7188/workproject/form_push/api/form_push/v1/form_push_http.pb.go:124
// \ngithub.com/go-kratos/kratos/v2/transport/http.(*Router).Handle.func1
// \n\tD:/jy7188/goproject/pkg/mod/github.com/go-kratos/kratos/v2@v2.5.3/transport/http/router.go:54
// \nnet/http.HandlerFunc.ServeHTTP
// \n\tD:/jy7188/go_sdk/go/src/net/http/server.go:2109
// \ngithub.com/go-kratos/kratos/v2/transport/http.(*Server).filter.func1.1
// \n\tD:/jy7188/goproject/pkg/mod/github.com/go-kratos/kratos/v2@v2.5.3/transport/http/server.go:278
// \nnet/http.HandlerFunc.ServeHTTP
// \n\tD:/jy7188/go_sdk/go/src/net/http/server.go:2109\ngithub.com/gorilla/mux.(*Router).ServeHTTP
// \n\tD:/jy7188/goproject/pkg/mod/github.com/gorilla/mux@v1.8.0/mux.go:210
// \nnet/http.serverHandler.ServeHTTP\n\tD:/jy7188/go_sdk/go/src/net/http/server.go:2947
// \nnet/http.(*conn).serve
// \n\tD:/jy7188/go_sdk/go/src/net/http/server.go:1991
// \nruntime.goexit
// \n\tD:/jy7188/go_sdk/go/src/runtime/asm_amd64.s:1594\n"