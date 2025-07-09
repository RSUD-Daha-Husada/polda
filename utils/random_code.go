package utils

import (
    "math/rand"
    "time"
)

func RandomCode(length int) string {
    letters := []rune("0123456789")
    r := rand.New(rand.NewSource(time.Now().UnixNano())) 

    b := make([]rune, length)
    for i := range b {
        b[i] = letters[r.Intn(len(letters))]
    }
    return string(b)
}
