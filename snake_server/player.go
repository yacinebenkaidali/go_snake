package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	streams map[string]*FileSound
}

type FileSound struct {
	file     *os.File
	streamer beep.StreamSeekCloser
}

func NewPlayer(files []string) *Player {
	player := &Player{
		streams: make(map[string]*FileSound),
	}
	for i := 0; i < len(files); i++ {
		f, err := os.Open(files[i])
		if err != nil {
			panic(fmt.Sprintf("unexpected error: couldn't file %s", files[i]))
		}

		streamer, format, err := mp3.Decode(f)
		if err != nil {
			panic(fmt.Sprintf("unexpected error: couldn't decode file %s, expected mp3 file", files[i]))
		}
		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
			panic(err)
		}

		fs := &FileSound{
			file:     f,
			streamer: streamer,
		}
		player.streams[files[i]] = fs

	}
	return player
}

func (p *Player) Close() {
	for _, v := range p.streams {
		v.streamer.Close()
		v.file.Close()
	}
}

func (p *Player) PlaySound(file string, reset bool) {
	if fs, ok := p.streams[file]; ok {
		speaker.Play(beep.Seq(fs.streamer, beep.Callback(func() {
			if reset {
				fs.ResetStream()
			}
		})))
	}
}

func (fs *FileSound) ResetStream() {
	fs.streamer.Seek(0)
}
