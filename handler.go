package amber

import (
    "github.com/CharLemAznable/gokits"
    "net/http"
    "net/url"
    "time"
)

func emptyHandler(_ http.ResponseWriter, _ *http.Request) {}

const CookieValueContextKey = "AmberCookieValue"

func AuthAmber(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        if nil == ConfigInstance {
            emptyHandler(writer, request)
            return
        }
        if !ConfigInstance.ForceLogin {
            handlerFunc(writer, request)
            return
        }
        if "" == ConfigInstance.AppId ||
            "" == ConfigInstance.EncryptKey ||
            "" == ConfigInstance.CookieName ||
            "" == ConfigInstance.AmberLoginUrl ||
            "" == ConfigInstance.LocalUrl {
            emptyHandler(writer, request)
            return
        }

        cookieValue, err := ReadCookieValue(request)
        if err == nil {
            ctx := gokits.ModelContextWithValue(request.Context(),
                CookieValueContextKey, cookieValue)
            handlerFunc(writer, request.WithContext(ctx))
            return
        }

        redirectUrl := ConfigInstance.AmberLoginUrl +
            "?appId=" + ConfigInstance.AppId +
            "&redirectUrl=" + url.QueryEscape(
            ConfigInstance.LocalUrl+request.RequestURI)
        http.Redirect(writer, request, redirectUrl, http.StatusFound)
    }
}

func ServeCocsHandler(writer http.ResponseWriter, request *http.Request) {
    if nil == ConfigInstance {
        gokits.ResponseErrorText(writer, http.StatusInternalServerError,
            http.StatusText(http.StatusInternalServerError))
        return
    }
    if "" == ConfigInstance.AppId ||
        "" == ConfigInstance.EncryptKey ||
        "" == ConfigInstance.CookieName ||
        "" == ConfigInstance.AmberLoginUrl ||
        "" == ConfigInstance.LocalUrl {
        gokits.ResponseErrorText(writer, http.StatusInternalServerError,
            http.StatusText(http.StatusInternalServerError))
        return
    }

    redirect := request.FormValue("redirect")
    if "" == redirect {
        redirect = ConfigInstance.LocalUrl
    }
    expires, err := gokits.Int64FromStr(request.FormValue("e"))
    if nil != err || 0 == expires {
        gokits.ResponseErrorText(writer, http.StatusBadRequest,
            http.StatusText(http.StatusBadRequest))
        return
    }
    cookieValue := request.FormValue(ConfigInstance.CookieName)
    if "" == cookieValue {
        gokits.ResponseErrorText(writer, http.StatusBadRequest,
            http.StatusText(http.StatusBadRequest))
        return
    }
    cookie := http.Cookie{Name: ConfigInstance.CookieName, Value: cookieValue,
        Path: "/", Expires: time.Unix(expires, 0), HttpOnly: true}
    http.SetCookie(writer, &cookie)
    http.Redirect(writer, request, redirect, http.StatusFound)
}
