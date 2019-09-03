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
    client := http.DefaultClient
    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        if "amber-login-url" == req.URL.Host {
            return http.ErrUseLastResponse
        }
        return nil
    }

    ConfigInstance = nil
    request, _ := http.NewRequest("GET", testServer.URL, nil)
    resp, _ := client.Do(request)
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    if "OK" != string(bodyBytes) {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithAppId("1000"),
        WithEncryptKey("0b4c09247ec02edc"),
        WithCookieName("cookie-test"),
        WithAmberLoginUrl("http://amber-login-url"),
        WithLocalUrl(testServer.URL),
    )
    cookieValueBuilt := &CookieValue{
        Username:    "john",
        Random:      ran.String(16),
        ExpiredTime: time.Now().Add(time.Second * 3),
    }
    encrypted := aesEncrypt(gokits.Json(cookieValueBuilt), ConfigInstance.EncryptKey)
    cookie := http.Cookie{Name: ConfigInstance.CookieName,
        Value: encrypted, Path: "/", Expires: cookieValueBuilt.ExpiredTime}

    request.AddCookie(&cookie)
    resp, _ = client.Do(request)
    bodyBytes, _ = ioutil.ReadAll(resp.Body)
    if "OK" != string(bodyBytes) {
        t.Fail()
    }

    time.Sleep(time.Second * 5)
    resp, _ = client.Do(request)
    if http.StatusFound != resp.StatusCode {
        t.Fail()
    }
    redirectUrl := ConfigInstance.AmberLoginUrl +
        "?appId=" + ConfigInstance.AppId +
        "&redirectUrl=" + url.QueryEscape(
        gokits.PathJoin(ConfigInstance.LocalUrl, request.RequestURI))
    if redirectUrl + url.QueryEscape("/") != resp.Header.Get("Location") {
        t.Fail()
    }
}
