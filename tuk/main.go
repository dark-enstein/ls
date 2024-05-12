package main

import (
	"context"
	"log"
	"tuk/internal"

	"github.com/spf13/pflag"
)

var path *string

func init() {
	path = pflag.StringP("path", "p", "", "path(s) to begin watching. if passing in multiple paths, separate them with a comma.")
	pflag.Parse()

	if *path == "" {
		log.Fatal("path is required")
	}
}

func main() {
	ctx := context.Background()
	plane, err := internal.NewPlane(ctx, internal.WithPaths(*path))
	if err != nil {
		log.Fatal(err)
	}
	defer plane.Close()

	// Set up listeners
	plane.Listen()

	plane.Log("Watching %s", *path)

	plane.Run()

	<-make(chan struct{})
}
