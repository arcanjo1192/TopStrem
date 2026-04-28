package main

import "topstrem/internal/mobile"

func main() {
    mobile.StartServer()
    select {} // keep the process alive
}