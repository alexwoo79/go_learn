package gadget

import (
	"fmt"
)

// interface 定义接口
type Player interface {
	Play(string)
	Stop()
}

// TapePlayer
type TapePlayer struct {
	Batteries string
}

// TapePlayer methods
func (t TapePlayer) Play(song string) {
	fmt.Println("Playing", song)
}

func (t TapePlayer) Stop() {
	fmt.Println("Stopped")

}

// TapeRecorder
type TapeRecorder struct {
	Microphones int
}

// TapeRecorder methods
func (t TapeRecorder) Play(song string) {
	fmt.Println("Playing", song)
}
func (t TapeRecorder) Stop() {
	fmt.Println("Stopped")
}
func (t TapeRecorder) Record() {
	fmt.Println("Recording")
}
