/*
 * Copyright: Oleg Borodin <onborodin@gmail.com>
 */

package pmlog

import (
    "log"
)

func LogDebug(message ...interface{}) {
    log.Println("debug:", message)
    return
}

func LogError(message ...interface{}) {
    log.Println("error:", message)
    return
}


func LogWarning(message ...interface{}) {
    log.Println("warning:", message)
    return
}

func LogInfo(message ...interface{}) {
    log.Println("info:", message)
    return
}

//EOF


