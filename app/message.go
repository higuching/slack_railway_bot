package apps

import (
)

// 表示するテキスト
func GetMessage() string {
    // トラブルが発生している関東の路線を取得
    lineInfos := getTroubleLines()
    if (lineInfos == nil) {
        // トラブル無し
        return ""
    }

    // 表示対象の路線を取得
    targetLines := getTargetLines()

    var message string
    for _, tal := range lineInfos {
        isDispLine := false
        if (len(targetLines) > 0) {
            // 表示対象が限定されている
            for _, trl := range targetLines {
                if (tal.Name == trl) {
                    isDispLine = true
                    break
                }
            }
        } else {
            // 全表示
            isDispLine = true
        }
        if (isDispLine) {
            message = message + tal.Name + " @ " + tal.Outline + "(" + tal.Details + ")" + "\n"
        }
    }
    return message
}
