package amber

import (
    "github.com/CharLemAznable/gokits"
    "github.com/bingoohuang/gou/ran"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "net/url"
    "testing"
    "time"
)

func okHandler(writer http.ResponseWriter, _ *http.Request) {
    gokits.ResponseText(writer, "OK")
}

func TestAuthAmber(t *testing.T) {
    testServer := httptest.NewServer(AuthAmber(okHandler))
    request, _ := http.NewRequest("GET", testServer.URL, nil)
    client := http.DefaultClient
    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        if "amber-login-url" == req.URL.Host {
            return http.ErrUseLastResponse
        }
        return nil
    }

    ConfigInstance = nil
    resp, _ := client.Do(request)
    if http.StatusOK != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    if "" != string(bodyBytes) {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithForceLogin(false),
    )
    resp, _ = client.Do(request)
    if http.StatusOK != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if "OK" != string(bodyBytes) {
        t.Fail()
    }

    ConfigInstance = NewConfig()
    resp, _ = client.Do(request)
    if http.StatusOK != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if "" != string(bodyBytes) {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithAppId("1000"),
        WithEncryptKey("0b4c09247ec02edc"),
        WithCookieName("cookie-test"),
        WithAmberLoginUrl("http://amber-login-url"),
        WithLocalUrl(testServer.URL),
    )
    redirectUrl := ConfigInstance.AmberLoginUrl +
        "?appId=" + ConfigInstance.AppId +
        "&redirectUrl=" + url.QueryEscape(
        ConfigInstance.LocalUrl+request.RequestURI) +
        url.QueryEscape("/")
    resp, _ = client.Do(request)
    if http.StatusFound != resp.StatusCode {
        t.Fail()
    }
    if redirectUrl != resp.Header.Get("Location") {
        t.Fail()
    }

    cookieValue := &CookieValue{
        Username:    "john",
        Random:      ran.String(16),
        ExpiredTime: time.Now().Add(time.Second * 3),
    }
    encrypted := gokits.AESEncrypt(gokits.Json(cookieValue), ConfigInstance.EncryptKey)
    cookie := http.Cookie{Name: ConfigInstance.CookieName,
        Value: encrypted, Path: "/", Expires: cookieValue.ExpiredTime}
    request.AddCookie(&cookie)
    resp, _ = client.Do(request)
    if http.StatusOK != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if "OK" != string(bodyBytes) {
        t.Fail()
    }

    time.Sleep(time.Second * 5)
    resp, _ = client.Do(request)
    if http.StatusFound != resp.StatusCode {
        t.Fail()
    }
    if redirectUrl != resp.Header.Get("Location") {
        t.Fail()
    }
}

func TestServeCocsHandler(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(ServeCocsHandler))
    client := http.DefaultClient
    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        if "amber-cocs-url" == req.URL.Host {
            return http.ErrUseLastResponse
        }
        return nil
    }

    ConfigInstance = nil
    request, _ := http.NewRequest("GET", testServer.URL, nil)
    resp, _ := client.Do(request)
    if http.StatusInternalServerError != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    if http.StatusText(http.StatusInternalServerError) != string(bodyBytes) {
        t.Fail()
    }

    ConfigInstance = NewConfig()
    resp, _ = client.Do(request)
    if http.StatusInternalServerError != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if http.StatusText(http.StatusInternalServerError) != string(bodyBytes) {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithAppId("1000"),
        WithEncryptKey("0b4c09247ec02edc"),
        WithCookieName("cookie-test"),
        WithAmberLoginUrl("http://amber-login-url"),
        WithLocalUrl("http://amber-cocs-url"),
    )
    resp, _ = client.Do(request)
    if http.StatusBadRequest != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if http.StatusText(http.StatusBadRequest) != string(bodyBytes) {
        t.Fail()
    }

    e := gokits.StrFromInt64(time.Now().Add(time.Second * 3).Unix())
    request, _ = http.NewRequest("GET", testServer.URL + "?e=" + e, nil)
    resp, _ = client.Do(request)
    if http.StatusBadRequest != resp.StatusCode {
        t.Fail()
    }
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if http.StatusText(http.StatusBadRequest) != string(bodyBytes) {
        t.Fail()
    }

    request, _ = http.NewRequest("GET", testServer.URL + "?e=" + e + "&cookie-test=cookieValue", nil)
    resp, _ = client.Do(request)
    if http.StatusFound != resp.StatusCode {
        t.Fail()
    }
    cookies := resp.Cookies()
    for _, cookie := range cookies {
        if "cookie-test" == cookie.Name {
            if "cookieValue" != cookie.Value {
                t.Fail()
            }
            if e != gokits.StrFromInt64(cookie.Expires.Unix()) {
                t.Fail()
            }
        }
    }
}
