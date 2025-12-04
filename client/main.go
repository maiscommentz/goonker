package main

import (
	"Goonker/client/ui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.Println("Start client")
	game := &Game{}

	log.Println(ui.WindowWidth)
	log.Println(ui.WindowHeight)

	ebiten.SetWindowSize(ui.WindowWidth, ui.WindowHeight)
	ebiten.SetWindowTitle(ui.GameTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
