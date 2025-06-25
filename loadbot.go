package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	lksdk "github.com/livekit/server-sdk-go/v2"
)

type logEntry struct {
	Time  time.Time `json:"time"`
	Bot   string    `json:"bot"`
	Room  string    `json:"room"`
	Event string    `json:"event"`
	Error string    `json:"error,omitempty"`
}

func main() {
	roomPrefix := flag.String("room-prefix", "demo", "Prefix for room names")
	rooms := flag.Int("rooms", 1, "Number of rooms")
	bots := flag.Int("bots", 1, "Bots per room")
	duration := flag.Duration("d", time.Minute, "Connection duration")
	videoFile := flag.String("video", "sample720p.ivf", "Video file path (.ivf)")
	audioFile := flag.String("audio", "sample.ogg", "Audio file path (.ogg)")
	url := flag.String("url", "wss://meet.cst.ro", "LiveKit WS URL")
	tokenURL := flag.String("token-url", "https://meet.cst.ro/token", "Token endpoint")
	debug := flag.Bool("debug", false, "Enable debug output")
	logJSON := flag.Bool("log", false, "Write ndjson logs to last_run.json")

	flag.Parse()

	if _, err := os.Stat(*videoFile); err != nil {
		log.Printf("video file error: %v", err)
		os.Exit(1)
	}
	if _, err := os.Stat(*audioFile); err != nil {
		log.Printf("audio file error: %v", err)
		os.Exit(1)
	}

	totalBots := (*rooms) * (*bots)
	fmt.Printf("Spawning %d bots across %d room(s)...\n", totalBots, *rooms)

	var (
		wg    sync.WaitGroup
		logCh chan logEntry
	)

	if *logJSON {
		logCh = make(chan logEntry, 100)
		f, err := os.Create("last_run.json")
		if err != nil {
			log.Fatalf("log file: %v", err)
		}
		go func() {
			enc := json.NewEncoder(f)
			for e := range logCh {
				enc.Encode(e)
			}
			f.Close()
		}()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < *rooms; i++ {
		roomName := fmt.Sprintf("%s-%02d", *roomPrefix, i)
		for j := 0; j < *bots; j++ {
			identity := fmt.Sprintf("bot-%03d-%d", i*(*bots)+j, time.Now().UnixMilli())
			wg.Add(1)
			go func(room, id string) {
				defer wg.Done()
				logEvt := func(event string, err error) {
					if logCh != nil {
						e := logEntry{Time: time.Now(), Bot: id, Room: room, Event: event}
						if err != nil {
							e.Error = err.Error()
						}
						logCh <- e
					}
				}

				token, err := fetchToken(*tokenURL, room, id)
				if err != nil {
					log.Printf("token fetch error: %v", err)
					logEvt("token_error", err)
					return
				}
				roomConn, err := lksdk.ConnectToRoomWithToken(*url, token, nil)
				if err != nil {
					log.Printf("connect error: %v", err)
					logEvt("connect_error", err)
					return
				}
				defer roomConn.Disconnect()
				logEvt("join", nil)

				if *debug {
					log.Printf("%s joined %s", id, room)
				}

				if vt, err := lksdk.NewLocalFileTrack(*videoFile); err == nil {
					if _, err := roomConn.LocalParticipant.PublishTrack(vt, nil); err != nil {
						log.Printf("video publish error: %v", err)
						logEvt("video_error", err)
					}
				} else {
					log.Printf("video track error: %v", err)
					logEvt("video_error", err)
				}
				if at, err := lksdk.NewLocalFileTrack(*audioFile); err == nil {
					if _, err := roomConn.LocalParticipant.PublishTrack(at, nil); err != nil {
						log.Printf("audio publish error: %v", err)
						logEvt("audio_error", err)
					}
				} else {
					log.Printf("audio track error: %v", err)
					logEvt("audio_error", err)
				}

				select {
				case <-time.After(*duration):
				case <-stop:
					if *debug {
						log.Printf("%s interrupted", id)
					}
				}
				logEvt("leave", nil)
			}(roomName, identity)
		}
	}
	wg.Wait()
	if logCh != nil {
		close(logCh)
	}
	fmt.Println("Done.")
}

func fetchToken(endpoint, room, identity string) (string, error) {
	u := fmt.Sprintf("%s?room=%s&identity=%s", endpoint, room, identity)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var j map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return "", err
	}
	return j["token"], nil
}
