/*
 * Copyright: Oleg Borodin <onborodin@gmail.com>
 */

package pmtools

import (
    "github.com/satori/go.uuid"
)

func GetNewUUID() string {
    id := uuid.NewV4()
    return id.String()
}

//EOF

