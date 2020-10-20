package ip

import (
	"context"
	"log"
	"net"
	"net/http"
)

//IPCtxKey is
var IPCtxKey = &contextKey{"client"}

type contextKey struct {
	name string
}

//DeviceDetails is
type DeviceDetails struct {
	IP   string
	Port string
}

//Middleware is
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, port, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				log.Println(err)
			}

			deviceDetails := &DeviceDetails{
				IP:   ip,
				Port: port,
			}

			//UserAgent: r.UserAgent(), NOT WORKING ON MOBILE APPS

			ctx := context.WithValue(r.Context(), IPCtxKey, deviceDetails)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
