package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func main() {
	tickerM := time.NewTicker(time.Minute * 20)
	tickerS := time.NewTicker(time.Second * 20)
	running := make(chan struct{})
	started := false

	go func() {
		for {
			select {
			case <-tickerM.C:
				started = true

				if err := playStart(); err != nil {
					log.Fatal(err)
				}
			case <-tickerS.C:
				if !started {
					break
				}

				started = false
				if err := playStop(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	log.Println("partime started, [ctrl + c] for stop")
	<-running
}

// playStart plays start sound
func playStart() (err error) {
	file, err := os.Open("start.mp3")
	if err != nil {
		return err
	}

	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}

	ctx, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}

	defer ctx.Close()

	player := ctx.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, decoder); err != nil {
		return err
	}

	return nil
}

// playStop plays stop sound
func playStop() (err error) {
	file, err := os.Open("stop.mp3")
	if err != nil {
		return err
	}

	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}

	ctx, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}

	defer ctx.Close()

	player := ctx.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, decoder); err != nil {
		return err
	}

	return nil
}
