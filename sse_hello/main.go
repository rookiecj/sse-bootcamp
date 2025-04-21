package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// 환경변수에서 prefix를 가져옴 (기본값: "")
	URL_PATH_PREFIX := os.Getenv("URL_PATH_PREFIX")
	
	// 정적 파일 서버 설정
	fs := http.FileServer(http.Dir("public"))
	if URL_PATH_PREFIX != "" {
		// prefix가 있는 경우 StripPrefix 사용
		http.Handle("/", http.StripPrefix(URL_PATH_PREFIX, fs))
	} else {
		// prefix가 없는 경우 루트에서 서비스
		http.Handle("/", fs)
	}

	// SSE 이벤트 핸들러 설정
	ssePath := URL_PATH_PREFIX + "/sse-events"
	fmt.Printf("SSE path: %s\n", ssePath)
	http.HandleFunc(ssePath, sseEventsHandler)

	fmt.Printf("Server starting on :8080 with URL_PATH_PREFIX: %s\n", URL_PATH_PREFIX)
	http.ListenAndServe(":8080", nil)
}

func sseEventsHandler(w http.ResponseWriter, r *http.Request) {
	// 현재 client의 ip 정보를 출력한다.
	log.Printf("Connected: %s\n", r.RemoteAddr)
	//fmt.Println(r.UserAgent())

	// CORS(Cross-Origin Resource Sharing)를 사용하려면 다음 헤더를 설정한다.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Handlers", "Content-Type")

	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	//w.Header().Set("Connection", "keep-alive")

	// 1초마다 현재시간을 내보낸다.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	defer func() {
		flusher.Flush()
		flusher = nil
		fmt.Println("flusher is nil")
	}()

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		id := 0
		for t := range ticker.C {
			// RFC3339 형식으로 시간을 포맷팅 (timezone 포함)
			nowStr := t.Format(time.RFC3339)
			fmt.Println("now", nowStr)
			// SSE format
			data := fmt.Sprintf("id: %d\ndata: %s\n", id, nowStr)
			fmt.Fprintf(w, "%s\n\n", data)
			if flusher != nil {
				flusher.Flush()
			} else {
				break
			}
			id++
		}
		ticker.Stop()
		fmt.Println("ticker closed")
	}()

	closeNotify := w.(http.CloseNotifier).CloseNotify()
	<-closeNotify

	log.Printf("Client Closed: %s\n", r.RemoteAddr)
}

// page를 refresh할때마다 새로운 connection이 연결된다.
//2024/08/01 01:14:49 Connected: [::1]:64727
//2024/08/01 01:14:55 Connected: [::1]:64749
//2024/08/01 01:14:56 Connected: [::1]:64754
//2024/08/01 01:14:57 Connected: [::1]:64759
//2024/08/01 01:14:59 Client Closed: [::1]:64727
//2024/08/01 01:15:05 Client Closed: [::1]:64749
//2024/08/01 01:15:06 Client Closed: [::1]:64754
// 마지막은 connection을 명시적으로 close한 경우이다.
//2024/08/01 01:15:14 Client Closed: [::1]:64759
