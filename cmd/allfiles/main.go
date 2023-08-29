package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/juanmera/allfiles/pkg/allfiles"
	"github.com/juanmera/allfiles/pkg/downloader"
	"github.com/juanmera/allfiles/pkg/humanize"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
)

type Context struct {
	File string
}
type LsCmd struct {
	ShowTotalSize bool     `name:"totals" short:"t" help:"Show total size"`
	ShowSize      bool     `name:"size" short:"s" help:"Show individual file size"`
	MinSize       string   `default:"1" help:"Min file size (0 includes empty files)"`
	MaxSize       string   `help:"Max file size (0 is unlimited)"`
	IncludeExts   []string `name:"exts" short:"e" help:"Include extesions"`
	ExcludeExts   []string `name:"exclude-exts" short:"x" help:"Exclude extesions"`
	IncludePaths  []string `name:"paths" short:"p" help:"Include paths"`
	ExcludePaths  []string `name:"exclude-paths" short:"q" help:"Exclude paths"`
	IncludeTypes  []string `name:"types" help:"Include file types"`
	ExcludeTypes  []string `name:"exclude-types" help:"Exclude file types"`
}

func (ls *LsCmd) Run(ctx *Context) error {
	ch, err := ls.StartFilter(ctx.File)
	if err != nil {
		return err
	}
	totalSize := 0
	for v := range ch {
		if ls.ShowSize {
			fmt.Printf("%s %s\n", v.HumanSize(), v.FilePath())
		} else {
			fmt.Printf("%s\n", v.FilePath())
		}
		totalSize += v.Size
	}
	if ls.ShowTotalSize {
		fmt.Printf("TOTAL: %s (%d)", humanize.FromBytes(totalSize), totalSize)
	}
	return nil
}

func (ls *LsCmd) StartFilter(file string) (chan *allfiles.File, error) {
	afp := allfiles.NewAllFiles(file)
	afp.SetMinSize(ls.MinSize)
	afp.SetMaxSize(ls.MaxSize)
	afp.ExcludeExts(ls.ExcludeExts...)
	afp.IncludeExts(ls.IncludeExts...)
	afp.ExcludeTypes(ls.ExcludeTypes...)
	afp.IncludeTypes(ls.IncludeTypes...)
	afp.ExcludePaths(ls.ExcludePaths...)
	afp.IncludePaths(ls.IncludePaths...)
	return afp.Start()
}

type DownloadCmd struct {
	LsCmd
	URL              url.URL `short:"u" required:""`
	OutputDir        string  `default:"./data"`
	ProxyURL         url.URL `name:"proxy"`
	Timeout          int
	Threads          int
	WarnAsErrorLimit int `default:"1" short:"w"`
}

func (dl *DownloadCmd) Run(ctx *Context) error {

	ch, err := dl.StartFilter(ctx.File)
	if err != nil {
		return errors.Join(fmt.Errorf("file not found %s", ctx.File), err)
	}
	dp := downloader.New(
		downloader.WithProxy(dl.ProxyURL.String()),
		downloader.WithTimeout(dl.Timeout),
		downloader.WithThreads(dl.Threads),
		downloader.WithWarnAsErrorLimit(dl.WarnAsErrorLimit),
	)
	dp.Start()
	defer dp.Wait()
	for v := range ch {
		fileURL := dl.URL.JoinPath(v.FilePath()).String()
		localPath := filepath.Join(dl.OutputDir, v.FilePath())
		dp.Send(downloader.NewFile(fileURL, localPath))
	}
	return nil
}

var cli struct {
	Debug bool
	File  string      `default:"ALL_FILES" short:"f" help:"Path to file"`
	Ls    LsCmd       `cmd:"" help:"List files"`
	Dl    DownloadCmd `cmd:"" help:"Download files"`
}

func main() {
	ctx := kong.Parse(&cli, kong.Configuration(kong.JSON, "./allfiles.json"))
	level := slog.LevelInfo
	if cli.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
	err := ctx.Run(&Context{File: cli.File})
	ctx.FatalIfErrorf(err)
}
