package amber

type Config struct {
    AppId         string
    EncryptKey    string
    CookieName    string
    AmberLoginUrl string
    LocalUrl      string
    ForceLogin    bool
}

type ConfigOptions struct {
    AppId         string
    EncryptKey    string
    CookieName    string
    AmberLoginUrl string
    LocalUrl      string
    ForceLogin    bool
}

var defaultConfigOptions = ConfigOptions{
    AppId:         "",
    EncryptKey:    "",
    CookieName:    "amber-login-auth",
    AmberLoginUrl: "",
    LocalUrl:      "",
    ForceLogin:    true,
}

type ConfigOption func(*ConfigOptions)

func WithAppId(appId string) ConfigOption {
    return func(o *ConfigOptions) { o.AppId = appId }
}

func WithEncryptKey(encryptKey string) ConfigOption {
    return func(o *ConfigOptions) { o.EncryptKey = encryptKey }
}

func WithCookieName(cookieName string) ConfigOption {
    return func(o *ConfigOptions) { o.CookieName = cookieName }
}

func WithAmberLoginUrl(amberLoginUrl string) ConfigOption {
    return func(o *ConfigOptions) { o.AmberLoginUrl = amberLoginUrl }
}

func WithLocalUrl(localUrl string) ConfigOption {
    return func(o *ConfigOptions) { o.LocalUrl = localUrl }
}

func WithForceLogin(forceLogin bool) ConfigOption {
    return func(o *ConfigOptions) { o.ForceLogin = forceLogin }
}

func NewConfig(opts ...ConfigOption) *Config {
    options := defaultConfigOptions
    for _, o := range opts {
        o(&options)
    }
    return &Config{
        AppId:         options.AppId,
        EncryptKey:    options.EncryptKey,
        CookieName:    options.CookieName,
        AmberLoginUrl: options.AmberLoginUrl,
        LocalUrl:      options.LocalUrl,
        ForceLogin:    options.ForceLogin,
    }
}

var ConfigInstance *Config = nil
