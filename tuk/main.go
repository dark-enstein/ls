package main

import (
	"context"
	"log"
	"strings"
	"tuk/internal"
	"tuk/internal/config"

	"github.com/spf13/pflag"
)

var path *string
var recursive *bool

type Tuk struct {
	plane  *internal.Plane
	config *config.Config
}

func init() {
	path = pflag.StringP("path", "p", "", "dir path(s) to begin watching. if passing in multiple paths, separate them with a comma.")
	recursive = pflag.BoolP("recursive", "r", false, "watch directories recursively")
	pflag.Parse()
}

func main() {
	// Load the configuration
	var err error
	var ctx context.Context

	// Init Tuk
	tuk := &Tuk{}
	tuk.config, err = config.LoadConfig(ctx, &config.Args{
		Path:      *path,
		Recursive: *recursive,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Init the plane
	tuk.plane, err = internal.NewPlane(ctx, internal.WithConfig(tuk.config))
	if err != nil {
		log.Fatal(err)
	}
	defer tuk.plane.Close()

	// Set up listeners
	tuk.plane.Listen()

	tuk.plane.Log("Watching %s", strings.Split(tuk.config.AllPaths(), ",")...)

	tuk.plane.Run()

	<-make(chan struct{})
}
