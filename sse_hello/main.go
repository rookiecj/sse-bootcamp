package main

import (
	"fmt"
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
	fmt.Printf("Connected: %s\n", r.RemoteAddr)
	//fmt.Println(r.UserAgent())

	// CORS(Cross-Origin Resource Sharing)를 사용하려면 다음 헤더를 설정한다.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Handlers", "Content-Type")

	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	//w.Header().Set("Connection", "keep-alive")

	// 30초간 1초마다 현재시간을 내보낸다.
	for i := 0; i < 30; i++ {
		now := time.Now()
		nowStr := now.Format("2006-01-02 15:04:05")
		data := fmt.Sprintf("data: %s\n", nowStr)
		fmt.Fprintf(w, "%s\n\n", data)
		time.Sleep(1 * time.Second)
		w.(http.Flusher).Flush()
	}

	closeNotify := w.(http.CloseNotifier).CloseNotify()
	<-closeNotify

	fmt.Printf("Closed: %s\n", r.RemoteAddr)
}
