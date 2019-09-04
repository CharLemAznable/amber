package amber

import (
    "fmt"
    "github.com/CharLemAznable/gokits"
    "github.com/bingoohuang/gou/ran"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestReadCookieValue(t *testing.T) {
    ConfigInstance = nil
    request := httptest.NewRequest("", "http://a.cn", nil)
    cookieValue, err := ReadCookieValue(request)
    if nil != cookieValue || nil == err {
        t.Fail()
    }
    if "未配置Amber.Config" != err.Error() {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithAppID("1000"),
        WithEncryptKey("0b4c09247ec02edc"),
        WithCookieName("cookie-test"),
        WithAmberLoginURL("amber-login-url"),
        WithLocalURL("local-url"),
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
    cookieValue, err = ReadCookieValue(request)
    if nil != err {
        t.Fail()
    }
    if cookieValueBuilt.Username != cookieValue.Username {
        t.Fail()
    }
    if cookieValueBuilt.Random != cookieValue.Random {
        t.Fail()
    }
    if cookieValueBuilt.ExpiredTime.Format("2006-01-02 15:04:05") !=
        cookieValue.ExpiredTime.Format("2006-01-02 15:04:05") {
        t.Fail()
    }

    time.Sleep(time.Second * 5)
    cookieValue, err = ReadCookieValue(request)
    if nil == err {
        fmt.Println("7")
        t.Fail()
    }
    if "cookie已过期" != err.Error() {
        fmt.Println("8")
        t.Fail()
    }
}
