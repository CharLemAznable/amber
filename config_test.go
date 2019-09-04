package amber

import (
    "testing"
)

func TestConfigInstance(t *testing.T) {
    ConfigInstance = NewConfig(
        WithAppID("1"),
        WithEncryptKey("2"),
        WithCookieName("3"),
        WithAmberLoginURL("4"),
        WithLocalURL("5"),
    )
    if !ConfigInstance.ForceLogin {
        t.Fail()
    }

    ConfigInstance = NewConfig(
        WithAppID("1"),
        WithEncryptKey("2"),
        WithCookieName("3"),
        WithAmberLoginURL("4"),
        WithLocalURL("5"),
        WithForceLogin(false),
    )
    if ConfigInstance.ForceLogin {
        t.Fail()
    }
}
