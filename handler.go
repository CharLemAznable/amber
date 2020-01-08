package amber

import (
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
