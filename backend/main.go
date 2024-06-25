package main

import (
    "fmt"
    "net/http"
)

func main() {
    port := "8080"

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Server online")
    })

    fmt.Printf("Executando na porta %s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        fmt.Printf("Erro ao iniciar o Server: %s\n", err)
        return
    }
}
