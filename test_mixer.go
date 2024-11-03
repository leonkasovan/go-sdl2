/*
TODO: Implement low level LoadSnd that return a Chunk
*/
package main

/*
#include <stdlib.h>
#if defined(__WIN32)
	#include <SDL2/SDL_mixer.h>
#else
	#include <SDL_mixer.h>
#endif
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/wav"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

type Error string

func (e Error) Error() string { return string(e) }

// type Chunk struct {
// 	allocated int32  // a boolean indicating whether to free abuf when the chunk is freed
// 	buf       *uint8 // pointer to the sample data, which is in the output format and sample rate
// 	len_      uint32 // length of abuf in bytes
// 	volume    uint8  // 0 = silent, 128 = max volume. This takes effect when mixing
// }

type Snd2 struct {
	table     map[[2]int32]*mix.Chunk
	ver, ver2 uint16
}

type Sound struct {
	wavData []byte
	format  beep.Format
	length  int
}

type Snd struct {
	table     map[[2]int32]*Sound
	ver, ver2 uint16
}

func main() {
	// Initialize SDL_mixer with a specific audio format
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 2048); err != nil {
		log.Fatalf("Could not initialize SDL_mixer: %v", err)
	}
	defer mix.CloseAudio()

	// Load the background music as streaming music
	backgroundMusic, err := mix.LoadMUS("background.mp3")
	if err != nil {
		log.Fatalf("Failed to load background music: %v", err)
	}
	defer backgroundMusic.Free()

	// Play the background music on a loop (-1 means loop indefinitely)
	if err := backgroundMusic.Play(-1); err != nil {
		log.Fatalf("Could not play background music: %v", err)
	}
	log.Println("Playing background music...")

	// Load a sound effect into memory (typically a short sound)
	soundEffect, err := mix.LoadWAV("sound_effect1.wav")
	if err != nil {
		log.Fatalf("Failed to load sound effect: %v", err)
	}
	defer soundEffect.Free()

	// Load a sound effect into memory (typically a short sound)
	soundEffect2, err := mix.LoadWAV("sound_effect2.wav")
	if err != nil {
		log.Fatalf("Failed to load sound effect: %v", err)
	}
	defer soundEffect.Free()

	// Query audio specifications (this might change based on the version)
	// frequency, format, channels, _, _ := mix.QuerySpec()
	// log.Printf("Audio Frequency: %d Hz, Format: %v, Channels: %d", frequency, format, channels)
	log.Printf("mix.AllocateChannels(0): %d", mix.AllocateChannels(-1))

	sndFileName := "test.snd"
	charSound, err := LoadSnd(sndFileName)
	if err != nil {
		log.Printf("Can't load %v: %v", sndFileName, err.Error())
	}

	charSound.IterateChunks()

	// Wait 2 seconds and play the sound effect on an available channel
	time.Sleep(2 * time.Second)
	channel, err := soundEffect.Play(-1, 0)
	if err != nil {
		log.Fatalf("Could not play sound effect: %v", err)
	}
	log.Printf("Played sound effect on channel %d", channel)

	// Wait 2 seconds and play the sound effect on an available channel
	time.Sleep(2 * time.Second)
	channel, err = soundEffect2.Play(-1, 0)
	if err != nil {
		log.Fatalf("Could not play sound effect: %v", err)
	}
	log.Printf("Played sound effect on channel %d", channel)

	// Keep the program running to let the audio play
	time.Sleep(3 * time.Second)
}

// ------------------------------------------------------------------
// Sound

func readSound(f *os.File, size uint32) (*Sound, error) {
	if size < 128 {
		return nil, fmt.Errorf("wav size is too small")
	}
	wavData := make([]byte, size)
	if _, err := f.Read(wavData); err != nil {
		return nil, err
	}
	// Decode the sound at least once, so that we know the format is OK
	s, fmt, err := wav.Decode(bytes.NewReader(wavData))
	if err != nil {
		return nil, err
	}
	// Check if the file can be fully played
	var samples [512][2]float64
	for {
		sn, _ := s.Stream(samples[:])
		if sn == 0 {
			// If sound wasn't able to be fully played, we disable it to avoid engine freezing
			if s.Position() < s.Len() {
				return nil, nil
			}
			break
		}
	}
	return &Sound{wavData, fmt, s.Len()}, nil
}

// LoadWAV loads file for use as a sample. This is actually mix.LoadWAVRW(sdl.RWFromFile(file, "rb"), 1). This can load WAVE, AIFF, RIFF, OGG, and VOC files. Note: You must call SDL_OpenAudio before this. It must know the output characteristics so it can convert the sample for playback, it does this conversion at load time. Returns: a pointer to the sample as a mix.Chunk.
// (https://www.libsdl.org/projects/SDL_mixer/docs/SDL_mixer_19.html)
func LoadWAV(wavData []byte, size uint32) (chunk *mix.Chunk, err error) {
	// why doesn't this call Mix_LoadWAV ?
	chunk = (*mix.Chunk)(unsafe.Pointer(C.Mix_LoadWAV_RW(C.SDL_RWFromMem(unsafe.Pointer(&wavData[0]), (C.int)(size)), 1)))
	if chunk == nil {
		err = sdl.GetError()
	}
	return
}

func readSound2(f *os.File, size uint32) (*mix.Chunk, error) {
	if size < 128 {
		return nil, fmt.Errorf("wav size is too small")
	}
	wavData := make([]byte, size)
	if _, err := f.Read(wavData); err != nil {
		return nil, err
	}

	return LoadWAV(wavData, size)
}

func (s *Sound) GetStreamer() beep.StreamSeeker {
	streamer, _, _ := wav.Decode(bytes.NewReader(s.wavData))
	return streamer
}

func newSnd() *Snd {
	return &Snd{table: make(map[[2]int32]*Sound)}
}

func newSnd2() *Snd2 {
	return &Snd2{table: make(map[[2]int32]*mix.Chunk)}
}

func LoadSnd(filename string) (*Snd2, error) {
	return LoadSndFiltered(filename, func(gn [2]int32) bool { return gn[0] >= 0 && gn[1] >= 0 }, 0)
}

// Parse a .snd file and return an Snd structure with its contents
// The "keepItem" function allows to filter out unwanted waves.
// If max > 0, the function returns immediately when a matching entry is found. It also gives up after "max" non-matching entries.
func LoadSndFiltered(filename string, keepItem func([2]int32) bool, max uint32) (*Snd2, error) {
	s := newSnd2()
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { f.Close() }()
	buf := make([]byte, 12)
	var n int
	if n, err = f.Read(buf); err != nil {
		return nil, err
	}
	if string(buf[:n]) != "ElecbyteSnd\x00" {
		return nil, Error("Unrecognized SND file, invalid header")
	}
	read := func(x interface{}) error {
		return binary.Read(f, binary.LittleEndian, x)
	}
	if err := read(&s.ver); err != nil {
		return nil, err
	}
	if err := read(&s.ver2); err != nil {
		return nil, err
	}
	var numberOfSounds uint32
	if err := read(&numberOfSounds); err != nil {
		return nil, err
	}
	var subHeaderOffset uint32
	if err := read(&subHeaderOffset); err != nil {
		return nil, err
	}
	loops := numberOfSounds
	if max > 0 && max < numberOfSounds {
		loops = max
	}
	log.Printf("%v numberOfSounds=%v", filename, numberOfSounds)
	for i := uint32(0); i < loops; i++ {
		f.Seek(int64(subHeaderOffset), 0)
		var nextSubHeaderOffset uint32
		if err := read(&nextSubHeaderOffset); err != nil {
			return nil, err
		}
		var subFileLength uint32
		if err := read(&subFileLength); err != nil {
			return nil, err
		}
		var num [2]int32
		if err := read(&num); err != nil {
			return nil, err
		}
		if keepItem(num) {
			_, ok := s.table[num]
			if !ok {
				tmp, err := readSound2(f, subFileLength)
				if err != nil {
					log.Printf("%v sound %v,%v can't be read: %v\n", filename, num[0], num[1], err)
					if max > 0 {
						return nil, err
					}
				} else {
					// Sound is corrupted and can't be played, so we export a warning message to the console
					if tmp == nil {
						log.Printf("WARNING: %v sound %v,%v is corrupted and can't be played, so it was disabled", filename, num[0], num[1])
					}
					s.table[num] = tmp
					if max > 0 {
						break
					}
				}
			}
		}
		subHeaderOffset = nextSubHeaderOffset
	}
	return s, nil
}

func (s *Snd2) IterateChunks() {
	for key, chunk := range s.table {
		// key is of type [2]int32
		// chunk is of type *mix.Chunk
		fmt.Printf("Key: %v, Chunk: %v\n", key, chunk)

		// Do something with chunk if needed
		if chunk != nil {
			// Example: access fields or methods of chunk
			chunk.Play(-1, 0)
			time.Sleep(1 * time.Second)
		}
	}
}
