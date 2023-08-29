package allfiles

import (
	"bufio"
	"github.com/juanmera/allfiles/pkg/extensions"
	"github.com/juanmera/allfiles/pkg/humanize"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type File struct {
	Path string
	Name string
	Size int
}

func (aff *File) FilePath() string {
	return filepath.Join(aff.Path, aff.Name)
}

func (aff *File) Ext() string {
	return strings.ToLower(filepath.Ext(aff.Name))
}

func (aff *File) HumanSize() string {
	return humanize.FromBytes(aff.Size)
}

type Parser struct {
	FilePath     string
	includeExts  []string
	excludeExts  []string
	includePaths []string
	excludePaths []string
	MinSize      int
	MaxSize      int
}

func NewAllFiles(filePath string) *Parser {
	return &Parser{
		FilePath:     filePath,
		excludePaths: make([]string, 0),
		excludeExts:  make([]string, 0),
		includePaths: make([]string, 0),
		includeExts:  make([]string, 0),
	}
}

func (afp *Parser) IncludeExts(exts ...string) {
	for _, v := range exts {
		if !strings.HasPrefix(v, ".") {
			v = "." + v
		}
		afp.includeExts = append(afp.includeExts, v)
	}
}

func (afp *Parser) IncludeTypes(types ...string) {
	afp.IncludeExts(extensions.Get(types...)...)
}

func (afp *Parser) ExcludeExts(exts ...string) {
	for _, v := range exts {
		if !strings.HasPrefix(v, ".") {
			v = "." + v
		}
		afp.excludeExts = append(afp.excludeExts, v)
	}
}

func (afp *Parser) ExcludeTypes(types ...string) {
	afp.ExcludeExts(extensions.Get(types...)...)
}

func (afp *Parser) IncludePaths(paths ...string) {
	afp.includePaths = append(afp.includePaths, paths...)
}

func (afp *Parser) ExcludePaths(paths ...string) {
	afp.excludePaths = append(afp.excludePaths, paths...)
}

func (afp *Parser) Start() (chan *File, error) {
	f, err := os.Open(afp.FilePath)
	if err != nil {
		return nil, err
	}
	c := make(chan *File)
	go func() {
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		relativePath := ""
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ":") {
				relativePath = strings.TrimSuffix(line, ":")
			} else if relativePath == "" {
				slog.Error("Empty relative path")
				break
			} else {
				if len(line) == 0 || line[0] == 'd' {
					continue
				}
				if !afp.validatePath(relativePath) {
					continue
				}
				file := afp.parseLine(relativePath, line)
				if file.Size < afp.MinSize {
					continue
				}
				if afp.MaxSize != 0 && file.Size > afp.MaxSize {
					continue
				}
				if !afp.validateExt(file.Ext()) {
					continue
				}
				c <- file
			}
		}
		close(c)
	}()
	return c, nil
}

func (afp *Parser) validatePath(path string) bool {
	for _, v := range afp.excludePaths {
		if strings.HasPrefix(path, v) {
			return false
		}
	}
	if len(afp.includePaths) == 0 {
		return true
	}
	for _, v := range afp.includePaths {
		if strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}

func (afp *Parser) validateExt(ext string) bool {
	var exts []string
	included := len(afp.includeExts) > 0
	if included {
		exts = afp.includeExts
	} else {
		exts = afp.excludeExts
	}
	if slices.Contains(exts, ext) {
		return included
	}
	return !included
}

func (afp *Parser) parseLine(path, line string) *File {
	parts := strings.SplitN(line, " ", 3)
	return &File{
		Path: path,
		Name: parts[2],
		Size: humanize.ToBytes(parts[1]),
	}
}

func (afp *Parser) SetMinSize(size string) {
	afp.MinSize = humanize.ToBytes(size)
}

func (afp *Parser) SetMaxSize(size string) {
	afp.MaxSize = humanize.ToBytes(size)
}
