package main

func main() {
	cleanup, err := NewReleaseCleanup()
	if err != nil {
		cleanup.Exit(err)
	}
	cleanup.Exit(cleanup.Cleanup())
}
