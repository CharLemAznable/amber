package amber

import (
    "context"
    "github.com/CharLemAznable/gokits"
    "net/http"
    "net/url"
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
        if "" == ConfigInstance.AppID ||
            "" == ConfigInstance.EncryptKey ||
            "" == ConfigInstance.CookieName ||
            "" == ConfigInstance.AmberLoginURL ||
            "" == ConfigInstance.LocalURL {
            emptyHandler(writer, request)
            return
        }

        cookieValue, err := ReadCookieValue(request)
        if err == nil {
            ctx := context.WithValue(request.Context(),
                CookieValueContextKey, cookieValue)
            handlerFunc(writer, request.WithContext(ctx))
            return
        }

        redirectUrl := ConfigInstance.AmberLoginURL +
            "?appID=" + ConfigInstance.AppID +
            "&redirectUrl=" + url.QueryEscape(
            gokits.PathJoin(ConfigInstance.LocalURL, request.RequestURI))
        http.Redirect(writer, request, redirectUrl, http.StatusFound)
    }
}
