package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/trace"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/golang/base"
	"github.com/rusq/slackdump/v3/internal/cache"
)

var CmdWspList = &base.Command{
	UsageLine: baseCommand + " list [flags]",
	Short:     "list saved authentication information",
	Long: `
# Auth List Command

**List** allows to list Slack Workspaces, that you have previously authenticated in.
`,
	FlagMask:   flagmask,
	PrintFlags: true,
}

const timeLayout = "2006-01-02 15:04:05"

var (
	bare = CmdWspList.Flag.Bool("b", false, "bare output format (just names)")
	all  = CmdWspList.Flag.Bool("a", false, "all information, including user")
)

func init() {
	CmdWspList.Run = runList
}

func runList(ctx context.Context, cmd *base.Command, args []string) error {
	m, err := cache.NewManager(cfg.CacheDir())
	if err != nil {
		base.SetExitStatus(base.SCacheError)
		return err
	}

	formatter := printFull
	if *bare {
		formatter = printBare
	} else if *all {
		formatter = printAll
	}

	entries, err := m.List()
	if err != nil {
		if errors.Is(err, cache.ErrNoWorkspaces) {
			base.SetExitStatus(base.SUserError)
			return errors.New("no authenticated workspaces, please run \"slackdump " + baseCommand + " new\"")
		}
		base.SetExitStatus(base.SCacheError)
		return err
	}
	current, err := m.Current()
	if err != nil {
		if !errors.Is(err, cache.ErrNoDefault) {
			base.SetExitStatus(base.SWorkspaceError)
			return fmt.Errorf("error getting the current workspace: %s", err)
		}
		current = entries[0]
		if err := m.Select(current); err != nil {
			base.SetExitStatus(base.SWorkspaceError)
			return fmt.Errorf("error setting the current workspace: %s", err)
		}

	}

	formatter(m, current, entries)
	return nil
}

const defMark = "=>"

var hdrItems = []hdrItem{
	{"C", 1},
	{"name", 7},
	{"filename", 12},
	{"modified", 19},
	{"team", 9},
	{"user", 8},
	{"error", 5},
}

func printAll(m manager, current string, wsps []string) {
	ctx, task := trace.NewTask(context.Background(), "printAll")
	defer task.End()

	tw := tabwriter.NewWriter(os.Stdout, 2, 8, 1, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, printHeader(hdrItems...))

	rows := wspInfo(ctx, m, current, wsps)
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
}

func wspInfo(ctx context.Context, m manager, current string, wsps []string) [][]string {
	var rows = [][]string{}

	var (
		wg   sync.WaitGroup
		rowC = make(chan []string)
		pool = make(chan struct{}, 8)
	)
	for _, name := range wsps {
		wg.Add(1)
		go func() {
			pool <- struct{}{}
			defer func() {
				<-pool
				wg.Done()
			}()
			rowC <- wspRow(ctx, m, current, name)
		}()
	}
	go func() {
		wg.Wait()
		close(rowC)
	}()
	for row := range rowC {
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][1] < rows[j][1]
	})
	return rows
}

func wspRow(ctx context.Context, m manager, current, name string) []string {
	curr := ""
	if current == name {
		curr = "*"
	}
	fi, err := m.FileInfo(name)
	if err != nil {
		return []string{curr, name, "", "", "", "", err.Error()}
	}
	info, err := userInfo(ctx, m, name)
	if err != nil {
		return []string{curr, name, fi.Name(), fi.ModTime().Format(timeLayout), "", "", err.Error()}
	}
	return []string{curr, name, fi.Name(), fi.ModTime().Format(timeLayout), info.Team, info.User, "OK"}
}

type hdrItem struct {
	name string
	size int
}

func (h *hdrItem) String() string {
	return h.name
}

func (h *hdrItem) Size() int {
	if h.size == 0 {
		h.size = len(h.String())
	}
	return h.size
}

func (h *hdrItem) Underline(char ...string) string {
	if len(char) == 0 {
		char = []string{"-"}
	}
	return strings.Repeat(char[0], h.Size())
}

func printHeader(hi ...hdrItem) string {
	var sb strings.Builder
	for i, h := range hi {
		if i > 0 {
			sb.WriteByte('\t')
		}
		sb.WriteString(h.String())
	}
	sb.WriteByte('\n')
	for i, h := range hi {
		if i > 0 {
			sb.WriteByte('\t')
		}
		sb.WriteString(h.Underline())
	}
	return sb.String()
}

func userInfo(ctx context.Context, m manager, name string) (*slack.AuthTestResponse, error) {
	prov, err := m.LoadProvider(name)
	if err != nil {
		return nil, err
	}
	return prov.Test(ctx)
}

func printFull(m manager, current string, wsps []string) {
	fmt.Printf("Workspaces in %q:\n\n", cfg.CacheDir())
	for _, row := range simpleList(context.Background(), m, current, wsps) {
		fmt.Printf("%s (file: %s, last modified: %s)", row[0], row[1], row[2])
	}
	fmt.Printf("\nCurrent workspace is marked with ' %s '.\n", defMark)
}

func simpleList(_ context.Context, m manager, current string, wsps []string) [][]string {
	var rows = make([][]string, 0, len(wsps))
	for _, name := range wsps {
		timestamp := "unknown"
		filename := "-"
		if fi, err := m.FileInfo(name); err == nil {
			timestamp = fi.ModTime().Format(timeLayout)
			filename = fi.Name()
		}
		if name == current {
			name = defMark + " " + name
		} else {
			name = "   " + name
		}
		rows = append(rows, []string{name, filename, timestamp})
	}
	return rows
}

func printBare(_ manager, current string, workspaces []string) {
	for _, name := range workspaces {
		if current == name {
			fmt.Print("*")
		}
		fmt.Println(name)
	}
}
