package handler

import (
	"net/http"

	"book/service/search/api/internal/logic"
	"book/service/search/api/internal/svc"
	"book/service/search/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func searchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSearchLogic(r.Context(), svcCtx)
		resp, err := l.Search(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// func TracingHandler(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // **1**
//         carrier, err := trace.Extract(trace.HttpFormat, r.Header)
//         // ErrInvalidCarrier means no trace id was set in http header
//         if err != nil && err != trace.ErrInvalidCarrier {
//             logx.Error(err)
//         }

//         // **2**
//         ctx, span := trace.StartServerSpan(r.Context(), carrier, sysx.Hostname(), r.RequestURI)
//         defer span.Finish()
//         // **5**
//         r = r.WithContext(ctx)

//         next.ServeHTTP(w, r)
//     })
// }

// func StartServerSpan(ctx context.Context, carrier Carrier, serviceName, operationName string) (context.Context, tracespec.Trace) {
//     span := newServerSpan(carrier, serviceName, operationName)
//     // **4**
//     return context.WithValue(ctx, tracespec.TracingKey, span), span
// }

// func newServerSpan(carrier Carrier, serviceName, operationName string) tracespec.Trace {
//     // **3**
//     traceId := stringx.TakeWithPriority(func() string {
//         if carrier != nil {
//             return carrier.Get(traceIdKey)
//         }
//         return ""
//     }, func() string {
//         return stringx.RandId()
//     })
//     spanId := stringx.TakeWithPriority(func() string {
//         if carrier != nil {
//             return carrier.Get(spanIdKey)
//         }
//         return ""
//     }, func() string {
//         return initSpanId
//     })

//     return &Span{
//         ctx: spanContext{
//             traceId: traceId,
//             spanId:  spanId,
//         },
//         serviceName:   serviceName,
//         operationName: operationName,
//         startTime:     timex.Time(),
//         // 标记为server
//         flag:          serverFlag,
//     }
// }
