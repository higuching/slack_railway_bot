package apps

import (
    "io/ioutil"
    yaml "gopkg.in/yaml.v2"
)

type TargetLines struct {
    Filter  bool        `yaml:"filter"`
    Lines   []string    `yaml:"lines"`
}

// 定義した情報の欲しい路線情報を取得
func getTargetLines() []string {
    buf, err := ioutil.ReadFile("configs/railways.yml")
    if err != nil {
        panic(err)
    }

    targetLines := TargetLines{}
    err = yaml.Unmarshal(buf, &targetLines)
    if err != nil {
      panic(err)
    }
    return targetLines.Lines
}
