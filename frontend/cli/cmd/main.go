package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const (
	chunkSize      = 64 * 1024              // 64KB
	initialBuffer  = 10 * chunkSize         // 640KB
	sampleRate     = 44100                  // 44.1 kHz
	bufferDuration = 100 * time.Millisecond // speaker latency
)

type streamBuffer struct {
	buf  *bytes.Buffer
	done chan struct{}
}

func newStreamBuffer() *streamBuffer {
	return &streamBuffer{
		buf:  bytes.NewBuffer(nil),
		done: make(chan struct{}),
	}
}

func (s *streamBuffer) Read(p []byte) (int, error) {
	for {
		n, err := s.buf.Read(p)
		if err == io.EOF {
			select {
			case <-s.done:
				return n, io.EOF
			default:
				time.Sleep(50 * time.Millisecond)
				continue
			}
		}
		return n, err
	}
}

func (s *streamBuffer) Write(data []byte) {
	s.buf.Write(data)
}

func (s *streamBuffer) Close() error {
	close(s.done)
	return nil
}

func streamMP3ToBuffer(filePath string, sb *streamBuffer, ready chan struct{}) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open mp3 file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, chunkSize)
	totalWritten := 0

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read mp3 file: %v", err)
		}

		sb.Write(buf[:n])
		totalWritten += n

		if totalWritten >= initialBuffer {
			select {
			case ready <- struct{}{}:
			default:
			}
		}

		time.Sleep(100 * time.Millisecond) // имитация сети
	}

	sb.Close()
}

func main() {
	filePath := "song2.mp3" // путь к вашему mp3

	sb := newStreamBuffer()
	ready := make(chan struct{}, 1)

	go streamMP3ToBuffer(filePath, sb, ready)

	// Ждём, пока будет достаточно данных в буфере
	<-ready

	streamer, format, err := mp3.Decode(sb)
	if err != nil {
		log.Fatalf("failed to decode mp3: %v", err)
	}
	defer streamer.Close()

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(bufferDuration)); err != nil {
		log.Fatalf("failed to initialize speaker: %v", err)
	}

	done := make(chan struct{})

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done
}
