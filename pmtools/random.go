/*
 * Copyright: Oleg Borodin <onborodin@gmail.com>
 */


package pmtools

import (
    "math/rand"
    "time"
)

func GetRandomPercent() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(100)
}

func GetRandomInt(min int, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min + 1) + min
}

func GetRandomBool() bool {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(2) == 1
}

//EOF

