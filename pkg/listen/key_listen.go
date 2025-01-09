package listen

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func main() {
	// low()
	add()
}

func add() {
	var listen = false
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("退出监听")
		hook.End()
	})
	hook.Register(hook.KeyDown, []string{"u"}, func(e hook.Event) {
		// fmt.Println(listen)
		if !listen {
			fmt.Println("开始录音")
			listen = true
			return
		}
		if listen {
			fmt.Println("结束录音")
			listen = false
		}
	})
	s := hook.Start()
	<-hook.Process(s)
}

// func low() {
// 	evChan := hook.Start()
// 	defer hook.End()

// 	for ev := range evChan {
// 		if ev.Kind != hook.MouseMove {
// 			fmt.Println("hook: ", ev)
// 		}
// 	}
// }
