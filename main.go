package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/starfederation/datastar-go/datastar"
)

var count atomic.Int64

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/increment", increment)
	http.HandleFunc("/datastar.js", datastarJS)

	fmt.Println("listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Datastar + Go</title>
  <script type="module" src="/datastar.js"></script>
</head>
<body>
  <main style="font-family: sans-serif; max-width: 32rem; margin: 3rem auto;">
    <h1>Datastar counter</h1>
    <p>A minimal Go + Datastar example.</p>
    <div data-signals='{"count": 0}'>
      <p>Count: <strong data-text="$count">0</strong></p>
      <button data-on:click="@get('/increment')">Increment</button>
    </div>
  </main>
</body>
</html>`)
}

func increment(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	next := count.Add(1)
	_ = sse.PatchSignals([]byte(fmt.Sprintf(`{"count": %d}`, next)))
}

func datastarJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "datastar.js")
}
