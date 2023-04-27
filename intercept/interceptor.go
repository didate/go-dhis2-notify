package intercept

import (
	"net/http"
	"os"
)

type Interceptor struct {
	Core http.RoundTripper
}

func (Interceptor) addBasicAuth(r *http.Request) *http.Request {
	r.Header.Add("Authorization", os.Getenv("DHIS2_AUTH"))
	return r
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	r = i.addBasicAuth(r)
	return i.Core.RoundTrip(r)
}
