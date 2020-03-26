package amber

import (
    "errors"
    "github.com/CharLemAznable/gokits"
    "net/http"
    "time"
)

type CookieValue struct {
    Username    string
    Random      string
    ExpiredTime time.Time
}

func ReadCookieValue(request *http.Request) (*CookieValue, error) {
    if nil == ConfigInstance {
        return nil, errors.New("未配置Amber.Config")
    }
    cookie, err := request.Cookie(ConfigInstance.CookieName)
    if err != nil {
        return nil, err
    }
    decrypted := gokits.AESDecrypt(cookie.Value, ConfigInstance.EncryptKey)
    if "" == decrypted {
        return nil, errors.New("cookie解密失败")
    }
    cookieValue, ok := gokits.UnJson(decrypted, new(CookieValue)).(*CookieValue)
    if !ok || nil == cookieValue {
        return nil, errors.New("cookie解析失败")
    }
    if cookieValue.ExpiredTime.Before(time.Now()) {
        return nil, errors.New("cookie已过期")
    }
    if "" == cookieValue.Username {
        return nil, errors.New("用户未登录")
    }
    return cookieValue, nil
}
