package amber

type Config struct {
    AppID         string
    EncryptKey    string
    CookieName    string
    AmberLoginURL string
    LocalURL      string
    ForceLogin    bool
}

type ConfigOptions struct {
    AppID         string
    EncryptKey    string
    CookieName    string
    AmberLoginURL string
    LocalURL      string
    ForceLogin    bool
}

var defaultConfigOptions = ConfigOptions{
    AppID:         "",
    EncryptKey:    "",
    CookieName:    "amber-login-auth",
    AmberLoginURL: "",
    LocalURL:      "",
    ForceLogin:    true,
}

type ConfigOption func(*ConfigOptions)

func WithAppID(appID string) ConfigOption {
    return func(o *ConfigOptions) { o.AppID = appID }
}

func WithEncryptKey(encryptKey string) ConfigOption {
    return func(o *ConfigOptions) { o.EncryptKey = encryptKey }
}

func WithCookieName(cookieName string) ConfigOption {
    return func(o *ConfigOptions) { o.CookieName = cookieName }
}

func WithAmberLoginURL(amberLoginURL string) ConfigOption {
    return func(o *ConfigOptions) { o.AmberLoginURL = amberLoginURL }
}

func WithLocalURL(localURL string) ConfigOption {
    return func(o *ConfigOptions) { o.LocalURL = localURL }
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
        AppID:         options.AppID,
        EncryptKey:    options.EncryptKey,
        CookieName:    options.CookieName,
        AmberLoginURL: options.AmberLoginURL,
        LocalURL:      options.LocalURL,
        ForceLogin:    options.ForceLogin,
    }
}

var ConfigInstance *Config = nil
