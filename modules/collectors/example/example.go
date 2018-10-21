package example

import (
	"math/rand"

	"github.com/l2isbad/go.d.plugin/modules"
)

type Example struct {
	modules.Base

	data map[string]int64
}

func (e *Example) Check() bool {
	return true
}

func (Example) GetCharts() *modules.Charts {
	return modules.NewCharts(uCharts...)
}

func (e *Example) GetData() map[string]int64 {
	e.data["random0"] = rand.Int63n(100)
	e.data["random1"] = rand.Int63n(100)

	return e.data
}

func init() {
	modules.Register("example", modules.Creator{
		Create: func() modules.Module {
			return &Example{}
		},
	})
}
