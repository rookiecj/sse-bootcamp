package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./public")))

	// '/events'는 event source
	http.HandleFunc("/sse-events", sseEventsHandler)
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

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		id := 0
		for t := range ticker.C {
			nowStr := t.Format("2006-01-02 15:04:05")
			// SSE format
			data := fmt.Sprintf("id: %d\ndata: %s\n", id, nowStr)
			fmt.Fprintf(w, "%s\n\n", data)
			flusher.Flush()
			id++
		}
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
