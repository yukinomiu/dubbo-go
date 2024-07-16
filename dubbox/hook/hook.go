package hook

import "sync"

type (
	Hook func()
)

var (
	afterRCLoadOnce  sync.Once
	afterRCLoadHooks = make([]Hook, 0, 16)
)

func RegisterAfterRCLoadHook(hook Hook) {
	if hook == nil {
		return
	}

	afterRCLoadHooks = append(afterRCLoadHooks, hook)
}

func ExecuteAfterRCLoadHook() {
	afterRCLoadOnce.Do(func() {
		for _, hook := range afterRCLoadHooks {
			hook()
		}
	})
}
