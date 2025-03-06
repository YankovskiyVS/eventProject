package main

func main() {
	server := NewAPISServer(":8080")
	server.Run()
}
