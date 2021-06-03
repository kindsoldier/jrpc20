/*
 * Copyright: Oleg Borodin <onborodin@gmail.com>
 */


package pmtools

func ArrayIncludes(array []string, target string) bool {
    for i := range array {
        if array[i] == target {
            return true
        }
    } 
    return false
}
//EOF
