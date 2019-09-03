package amber

import (
    "context"
    "github.com/CharLemAznable/gokits"
    "net/http"
    "net/url"
)

const CookieValueContextKey = "AmberCookieValue"

func AuthAmber(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        if nil == ConfigInstance ||
            !ConfigInstance.ForceLogin || "" == ConfigInstance.AppId ||
            "" == ConfigInstance.EncryptKey || "" == ConfigInstance.CookieName ||
            "" == ConfigInstance.AmberLoginUrl || "" == ConfigInstance.LocalUrl {
            handlerFunc(writer, request)
            return
        }

        cookieValue, err := ReadCookieValue(request)
        if err == nil {
            ctx := context.WithValue(request.Context(),
                CookieValueContextKey, cookieValue)
            handlerFunc(writer, request.WithContext(ctx))
            return
        }

        redirectUrl := ConfigInstance.AmberLoginUrl +
            "?appId=" + ConfigInstance.AppId +
            "&redirectUrl=" + url.QueryEscape(
            gokits.PathJoin(ConfigInstance.LocalUrl, request.RequestURI))
        http.Redirect(writer, request, redirectUrl, http.StatusFound)
    }
}
