package main

import (
	"fmt"

	"github.com/alexwoo79/go_learn/gadget"
)

func playList(device gadget.Player, songs []string) {
	for _, song := range songs {
		device.Play(song)
	}
	device.Stop()
}
func TryOut(player gadget.Player) {
	player.Play("song.mp3")
	player.Stop()
	recorder, ok := player.(gadget.TapeRecorder)
	if ok {
		recorder.Record()
	}

}
func main() {
	mixtape := []string{"Jessie's Girl", "Whip It", "9 to 5"}
	var player gadget.Player = gadget.TapePlayer{}
	fmt.Println("TapePlayer Now playing:")
	playList(player, mixtape)
	fmt.Println("TapeRecorder Now playing:")
	player = gadget.TapeRecorder{}
	playList(player, mixtape)

	TryOut(gadget.TapePlayer{})
	TryOut(gadget.TapeRecorder{})

}
