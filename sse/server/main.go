package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// サーバー側の設定
const (
	chunkSize = 3 // 1回に送信する文字数
)

var (
	sampleText string
	chunks     []string
)

func init() {
	// テキストファイルの読み込み
	data, err := os.ReadFile("./sample.txt")
	if err != nil {
		log.Fatalf("Error reading sample.txt: %v", err)
	}
	sampleText = string(data)

	// テキストを事前にチャンク化
	words := strings.Split(sampleText, "")
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], ""))
	}
}

// getChunk returns a Message containing the accumulated text up to the specified chunk index
func getChunk(chunkIndex int) Message {
	var content string
	for i := 0; i <= chunkIndex; i++ {
		content += chunks[i]
	}
	return Message{
		ID:      chunkIndex + 1,
		Content: content,
		Done:    chunkIndex == len(chunks)-1,
	}
}

func main() {
	http.HandleFunc("/stream", handleStream)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	// SSEヘッダーの設定
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// クライアントが切断したことを検知するためのコンテキスト
	ctx := r.Context()

	// Last-Event-IDヘッダーの取得
	lastEventId := 0
	if idStr := r.Header.Get("Last-Event-ID"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			lastEventId = id
			log.Printf("Client requested to resume from ID: %d", lastEventId)
		}
	}

	// メッセージを送信するためのチャネル
	messageChan := make(chan Message)
	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	// メッセージを生成するゴルーチン
	go func() {
		for i := lastEventId; i < len(chunks); i++ {
			select {
			case <-ctx.Done():
				return
			default:
				messageChan <- getChunk(i)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// メッセージをクライアントに送信
	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		case <-pingTicker.C:
			fmt.Fprintf(w, ": ping\n\n")
			w.(http.Flusher).Flush()
		case message := <-messageChan:
			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling JSON: %v", err)
				continue
			}
			fmt.Fprintf(w, "id: %d\ndata: %s\n\n", message.ID, jsonData)
			w.(http.Flusher).Flush()

			if message.Done {
				return
			}
		}
	}
}
