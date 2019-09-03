package amber

import (
    "testing"
)

func TestAesEncrypt(t *testing.T) {
    encrypted := aesEncrypt("The quick brown fox jumps over the lazy dog", "0b4c09247ec02edc")
    if "3781dU72kqM+ulqyVv7aQlEoowO5jjGkTIjNNPKILa06LZ61DrAl7bhFFR20Ioao" != encrypted {
        t.Fail()
    }
}

func TestAesDecrypt(t *testing.T) {
    decrypted := aesDecrypt("3781dU72kqM+ulqyVv7aQlEoowO5jjGkTIjNNPKILa06LZ61DrAl7bhFFR20Ioao", "0b4c09247ec02edc")
    if "The quick brown fox jumps over the lazy dog" != decrypted {
        t.Fail()
    }
}
