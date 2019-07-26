package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yansal/youtube-ar/cmd"
)

func usage() string {
	// TODO: generate automatically from commands in package cmd
	return `usage: youtube-ar [create-url|create-urls-from-youtube-playlist|download-url|list-logs|list-urls|server|worker]`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage())
		return
	}

	cmds := map[string]cmd.Cmd{
		"create-url":                        cmd.CreateURL,
		"create-urls-from-youtube-playlist": cmd.CreateURLsFromYoutubePlaylist,
		"download-url":                      cmd.DownloadURL,
		"list-logs":                         cmd.ListLogs,
		"list-urls":                         cmd.ListURLs,
		"server":                            cmd.Server,
		"worker":                            cmd.Worker,
	}

	cmd, ok := cmds[os.Args[1]]
	if !ok {
		fmt.Printf("error: unknown cmd %s\n", os.Args[1])
		fmt.Println(usage())
		os.Exit(2)
	}

	if err := cmd(context.Background(), os.Args[2:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
