/*
 * Copyright: Oleg Borodin <onborodin@gmail.com>
 */


package pmtools

import (
    "time"
)

func GetIsoTimestamp() string {
    return time.Now().Format(time.RFC3339)
}


//EOF

