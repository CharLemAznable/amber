package amber

import (
    "testing"
)

func TestConfigInstance(t *testing.T) {
    ConfigInstance = NewConfig(
        WithAppId("1"),
        WithEncryptKey("2"),
        WithCookieName("3"),
        WithAmberLoginUrl("4"),
        WithLocalUrl("5"),
    )
    if !ConfigInstance.ForceLogin {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithAppId("1"),
        WithEncryptKey("2"),
        WithCookieName("3"),
        WithAmberLoginUrl("4"),
        WithLocalUrl("5"),
        WithForceLogin(false),
    )
    if ConfigInstance.ForceLogin {
        t.Fail()
    }
}
