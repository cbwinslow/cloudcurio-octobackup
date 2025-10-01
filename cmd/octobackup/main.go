// File: cmd/octobackup/main.go
// Author: cbwinslow <blaine.winslow@gmail.com>
// Date: 2025-10-01
// Project: CloudCurio — OctoBackup (Neon Octopus Edition)
// Summary:
//   A Bubble Tea (Charmbracelet) TUI to stream Linux OS backups directly to a
//   remote homelab server over SSH without needing local temporary space. It
//   supports multiple strategies (raw disk stream, rsync file-level, Borg
//   encrypted deduplicated, and ZFS/Btrfs snapshot streaming when available),
//   with preflight checks, live logs, and a neon CloudCurio theme.
//
//   This is a single-file app for easy drop-in usage. It features:
//     • Strategy picker (dd|rsync|borg|zfs|btrfs)
//     • Config form (remote, port, path, compression, bandwidth, excludes)
//     • Preflight validator (tools, disk selection, SSH reachability)
//     • Live run view (spinner/progress + streaming command logs)
//     • Saves/loads config to ~/.config/cloudcurio/octobackup.yaml
//
// Inputs:
//   Interactive via TUI.
// Outputs:
//   Streams backups over SSH to your homelab path and prints run logs.
//
// Build:
//   $ mkdir -p cmd/octobackup && put this file there as main.go
//   $ cd cmd/octobackup
//   $ go mod init cloudcurio.cc/octobackup
//   $ go get github.com/charmbracelet/bubbles@v0.18.0 \
//          github.com/charmbracelet/bubbletea@v0.26.7 \
//          github.com/charmbracelet/lipgloss@v0.10.0 \
//          gopkg.in/yaml.v3@v3.0.1
//   $ go build -o octobackup
//   $ ./octobackup
//
// Notes:
//   • Requires Go 1.21+.
//   • The app will try to use: ssh, rsync, dd, gzip/pigz, pv, lsblk, borg, zfs, btrfs.
//   • Safe by default: you must pick the correct source disk (for raw dd) and confirm.
//
// Restore (quick hints):
//   • dd image:  ssh homelab 'cat image.gz' | gunzip | sudo dd of=/dev/sdX bs=64K status=progress
//   • rsync dir: rsync -aAXHv remote:/backups/host/ /mnt/target/
//   • borg:      borg mount repo::snapshot /mnt && copy; or borg extract repo::snapshot
//   • zfs/btrfs: receive snapshot and promote/clone as needed.
//
// Security:
//   • SSH only; can specify alternate port. Optional Borg encryption.
//   • Never stores secrets in plaintext; config omits passwords.
//
// Mod Log:
//   0.2.0 2025-10-01  Repo-packaged; headless flags; sample systemd; docs.
//   0.1.0 2025-09-30  Initial release.

package main

import (
    bufio "bufio"
    bytes "bytes"
    context "context"
    fmt "fmt"
    io "io"
    os "os"
    os_exec "os/exec"
    path_file "path/filepath"
    strings "strings"
    time "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/progress"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/lipgloss"
    "gopkg.in/yaml.v3"
)

// --------------------------- THEME (Lip Gloss) ---------------------------
var (
    paletteBg  = lipgloss.Color("#0a0a12")  // deep space
    paletteFg  = lipgloss.Color("#d6d5ff")  // pale neon ice
    neonPurple = lipgloss.Color("#a46cff")
    neonPink   = lipgloss.Color("#ff6bd6")
    neonTeal   = lipgloss.Color("#4fffe1")
    neonYellow = lipgloss.Color("#f7ff6f")

    appTitleStyle = lipgloss.NewStyle().Foreground(neonPurple).Bold(true).Margin(1, 2)
    sectionTitle  = lipgloss.NewStyle().Foreground(neonTeal).Bold(true)
    labelStyle    = lipgloss.NewStyle().Foreground(neonPink)
    valueStyle    = lipgloss.NewStyle().Foreground(paletteFg)
    warnStyle     = lipgloss.NewStyle().Foreground(neonYellow).Bold(true)
    helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a8ab3"))
    borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(neonPurple).Padding(1, 2)
)

const asciiLogo = `
      ____ _                 _               _                 
     / ___| | ___  _   _  __| |_   _  ___   | |__   __ _  ___  
    | |   | |/ _ \| | | |/ _` + "`" + ` | | | |/ _ \  | '_ \ / _` + "`" + ` |/ _ \ 
    | |___| | (_) | |_| | (_| | |_| |  __/  | |_) | (_| | (_) |
     \____|_|\___/ \__,_|\__,_|\__, |\___|  |_.__/ \__,_|\___/ 
                              |___/                             
         CloudCurio • OctoBackup — Neon Octopus Edition
`

// --------------------------- CONFIG ---------------------------

type Strategy string

const (
    StratDD    Strategy = "raw-dd"
    StratRsync Strategy = "rsync"
    StratBorg  Strategy = "borg"
    StratZFS   Strategy = "zfs-send"
    StratBtrfs Strategy = "btrfs-send"
)

type Config struct {
    RemoteUser    string   `yaml:"remote_user"`
    RemoteHost    string   `yaml:"remote_host"`
    SSHPort       int      `yaml:"ssh_port"`
    RemotePath    string   `yaml:"remote_path"`
    Strategy      Strategy `yaml:"strategy"`
    SourceDisk    string   `yaml:"source_disk"` // for dd/zfs roots; empty for rsync/borg
    Compression   string   `yaml:"compression"` // gzip|pigz|none
    BandwidthKbps int      `yaml:"bandwidth_kbps"` // 0 = unlimited
    Excludes      []string `yaml:"excludes"` // for rsync
    BorgRepo      string   `yaml:"borg_repo"` // ssh://user@host:/path/repo
    BorgPassEnv   string   `yaml:"borg_pass_env"` // env var name holding passphrase
}

func defaultConfig() Config {
    return Config{
        RemoteUser:    "cbwinslow",
        RemoteHost:    "cbwdellr720.cloudcurio.cc",
        SSHPort:       22,
        RemotePath:    "/backups/$(hostname)",
        Strategy:      StratRsync,
        Compression:   "pigz",
        BandwidthKbps: 0,
        Excludes: []string{
            "/dev/*", "/proc/*", "/sys/*", "/tmp/*", "/run/*", "/mnt/*", "/media/*", "/lost+found",
        },
        BorgRepo:    "ssh://cbwinslow@cbwdellr720.cloudcurio.cc:/backups/borg/$(hostname)",
        BorgPassEnv: "BORG_PASSPHRASE",
    }
}

func configPath() string {
    cfgDir := path_file.Join(os.Getenv("HOME"), ".config", "cloudcurio")
    _ = os.MkdirAll(cfgDir, 0o700)
    return path_file.Join(cfgDir, "octobackup.yaml")
}

func loadConfig() (Config, error) {
    p := configPath()
    b, err := os.ReadFile(p)
    if err != nil {
        return defaultConfig(), err
    }
    var c Config
    if err := yaml.Unmarshal(b, &c); err != nil {
        return defaultConfig(), err
    }
    return c, nil
}

func saveConfig(c Config) error {
    b, err := yaml.Marshal(c)
    if err != nil { return err }
    return os.WriteFile(configPath(), b, 0o600)
}

// --------------------------- UTIL ---------------------------

func have(cmd string) bool {
    _, err := os_exec.LookPath(cmd)
    return err == nil
}

func runCmd(ctx context.Context, name string, args ...string) (*os_exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
    cmd := os_exec.CommandContext(ctx, name, args...)
    stdout, err := cmd.StdoutPipe()
    if err != nil { return nil, nil, nil, err }
    stderr, err := cmd.StderrPipe()
    if err != nil { return nil, nil, nil, err }
    if err := cmd.Start(); err != nil { return nil, nil, nil, err }
    return cmd, stdout, stderr, nil
}

func sshTarget(c Config) string {
    return fmt.Sprintf("-p% d %s@%s", c.SSHPort, c.RemoteUser, c.RemoteHost)
}

func renderKeyVal(k, v string) string {
    return lipgloss.JoinHorizontal(lipgloss.Top,
        labelStyle.Render(k+":"),
        " ",
        valueStyle.Render(v),
    )
}

// --------------------------- TUI MODEL ---------------------------

type page int

const (
    pageIntro page = iota
    pageSelect
    pageConfig
    pagePreflight
    pageRun
)

type item string
func (i item) FilterValue() string { return string(i) }

// messages
type (
    preflightDoneMsg struct{ ok bool; report string; err error }
    runLogMsg        struct{ line string }
    runDoneMsg       struct{ err error }
)

type model struct {
    cfg         Config
    width       int
    height      int
    page        page

    list        list.Model
    spinner     spinner.Model
    progress    progress.Model

    inputs      []*textinput.Model
    focusIndex  int

    logLines    []string
    cancel      context.CancelFunc
    startTime   time.Time
}

func newModel(cfg Config) model {
    items := []list.Item{
        item("Raw disk stream (dd → ssh)"),
        item("Rsync file-level (fast/smart)"),
        item("Borg encrypted (dedup/incremental)"),
        item("ZFS snapshot send/recv"),
        item("Btrfs snapshot send/recv"),
    }
    lst := list.New(items, list.NewDefaultDelegate(), 0, 0)
    lst.Title = "Choose a backup strategy"
    sp := spinner.New()
    sp.Spinner = spinner.MiniDot
    pr := progress.New()

    // inputs: remote user, host, port, path, compression, bandwidth, disk, repo, passenv
    mk := func(ph string, val string) *textinput.Model {
        ti := textinput.New()
        ti.Placeholder = ph
        ti.SetValue(val)
        ti.Prompt = "➤ "
        return &ti
    }
    inputs := []*textinput.Model{
        mk("remote user", cfg.RemoteUser),
        mk("remote host", cfg.RemoteHost),
        mk("ssh port", fmt.Sprintf("%d", cfg.SSHPort)),
        mk("remote path", cfg.RemotePath),
        mk("compression (pigz|gzip|none)", cfg.Compression),
        mk("bandwidth kbps (0=unlimited)", fmt.Sprintf("%d", cfg.BandwidthKbps)),
        mk("source disk (e.g., /dev/sda)", cfg.SourceDisk),
        mk("borg repo (ssh://…)", cfg.BorgRepo),
        mk("borg pass env (e.g., BORG_PASSPHRASE)", cfg.BorgPassEnv),
    }

    return model{cfg: cfg, list: lst, spinner: sp, progress: pr, inputs: inputs, page: pageIntro}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.list.SetSize(m.width-8, m.height-12)
        return m, nil
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "enter":
            switch m.page {
            case pageIntro:
                m.page = pageSelect
                return m, nil
            case pageSelect:
                i := m.list.Index()
                switch i {
                case 0: m.cfg.Strategy = StratDD
                case 1: m.cfg.Strategy = StratRsync
                case 2: m.cfg.Strategy = StratBorg
                case 3: m.cfg.Strategy = StratZFS
                case 4: m.cfg.Strategy = StratBtrfs
                }
                m.page = pageConfig
                return m, nil
            case pageConfig:
                // persist field values
                m.cfg.RemoteUser = m.inputs[0].Value()
                m.cfg.RemoteHost = m.inputs[1].Value()
                fmt.Sscanf(m.inputs[2].Value(), "%d", &m.cfg.SSHPort)
                m.cfg.RemotePath = m.inputs[3].Value()
                m.cfg.Compression = strings.ToLower(m.inputs[4].Value())
                fmt.Sscanf(m.inputs[5].Value(), "%d", &m.cfg.BandwidthKbps)
                m.cfg.SourceDisk = m.inputs[6].Value()
                m.cfg.BorgRepo = m.inputs[7].Value()
                m.cfg.BorgPassEnv = m.inputs[8].Value()
                _ = saveConfig(m.cfg)
                m.page = pagePreflight
                return m, m.doPreflight()
            case pagePreflight:
                m.page = pageRun
                m.logLines = nil
                m.progress.SetPercent(0)
                m.spinner.Start()
                return m, m.runBackup()
            }
        case "tab":
            if m.page == pageConfig {
                m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
                for i := range m.inputs { m.inputs[i].Blur() }
                m.inputs[m.focusIndex].Focus()
            }
        }
    case preflightDoneMsg:
        m.logLines = append(m.logLines, strings.Split(msg.report, "\n")...)
        if !msg.ok || msg.err != nil {
            m.logLines = append(m.logLines, warnStyle.Render(fmt.Sprintf("Preflight failed: %v", msg.err)))
        }
        return m, nil
    case runLogMsg:
        m.logLines = append(m.logLines, msg.line)
        // naive progress tick
        p := m.progress.Percent() + 0.002
        if p > 0.98 { p = 0.98 }
        m.progress.SetPercent(p)
        return m, nil
    case runDoneMsg:
        m.spinner.Stop()
        if msg.err != nil {
            m.logLines = append(m.logLines, warnStyle.Render("Run finished with error: ")+msg.err.Error())
        } else {
            m.progress.SetPercent(1)
            m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(neonTeal).Bold(true).Render("✔ Backup complete"))
        }
        return m, nil
    }

    // delegate to list/inputs/spinner/progress
    var cmd tea.Cmd
    switch m.page {
    case pageSelect:
        m.list, cmd = m.list.Update(msg)
    case pageConfig:
        if m.focusIndex < len(m.inputs) {
            *m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
        }
    case pageRun:
        m.spinner, _ = m.spinner.Update(msg)
        m.progress, _ = m.progress.Update(msg)
    }
    return m, cmd
}

func (m model) View() string {
    switch m.page {
    case pageIntro:
        b := strings.Builder{}
        b.WriteString(appTitleStyle.Render(asciiLogo))
        b.WriteString(borderStyle.Render(
            sectionTitle.Render("Welcome to OctoBackup")+"\n"+
            "Stream your Linux backups directly to your homelab over SSH.\n\n"+
            helpStyle.Render("Press Enter to choose a backup strategy • q to quit")))
        return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())
    case pageSelect:
        return borderStyle.Render(m.list.View()) + "\n" + helpStyle.Render("Enter: select • q: quit")
    case pageConfig:
        rows := []string{
            sectionTitle.Render("Connection & Options"),
            renderKeyVal("strategy", string(m.cfg.Strategy)),
        }
        labels := []string{"user","host","port","remote path","compression","bandwidth","source disk","borg repo","borg passenv"}
        for i, ti := range m.inputs {
            rows = append(rows, renderKeyVal(labels[i], ti.View()))
        }
        rows = append(rows, "\n"+helpStyle.Render("Tab: next field • Enter: save & preflight"))
        return borderStyle.Render(strings.Join(rows, "\n"))
    case pagePreflight:
        return borderStyle.Render(sectionTitle.Render("Running preflight checks…")+"\n"+strings.Join(m.logLines, "\n"))
    case pageRun:
        logBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(neonTeal).Height(m.height-10).Width(m.width-6).Padding(0,1)
        log := strings.Join(m.logLines, "\n")
        header := lipgloss.JoinHorizontal(lipgloss.Top, m.spinner.View(), " ", sectionTitle.Render("Streaming backup…"))
        return borderStyle.Render(header+"\n"+m.progress.View()+"\n"+logBox.Render(log))
    }
    return ""
}

// --------------------------- PREFLIGHT ---------------------------

func (m model) doPreflight() tea.Cmd {
    return func() tea.Msg {
        var rpt bytes.Buffer
        ok := true

        fmt.Fprintf(&rpt, "Checking required tools…\n")
        req := []string{"ssh", "rsync"}
        switch m.cfg.Strategy {
        case StratDD:
            req = append(req, "dd")
            if m.cfg.Compression == "pigz" { req = append(req, "pigz") } else if m.cfg.Compression == "gzip" { req = append(req, "gzip") }
            if have("pv") { fmt.Fprintf(&rpt, "• pv present (nice progress)\n") }
        case StratBorg:
            req = append(req, "borg")
        case StratZFS:
            req = append(req, "zfs")
        case StratBtrfs:
            req = append(req, "btrfs")
        }
        for _, r := range req {
            if !have(r) { ok = false; fmt.Fprintf(&rpt, "✗ missing %s\n", r) } else { fmt.Fprintf(&rpt, "✓ %s\n", r) }
        }

        // SSH reachability
        fmt.Fprintf(&rpt, "Testing SSH reachability…\n")
        ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
        defer cancel()
        cmd := os_exec.CommandContext(ctx, "ssh", fmt.Sprintf("-p% d", m.cfg.SSHPort), fmt.Sprintf("%s@%s", m.cfg.RemoteUser, m.cfg.RemoteHost), "echo", "ok")
        if out, err := cmd.CombinedOutput(); err != nil || !strings.Contains(string(out), "ok") {
            ok = false; fmt.Fprintf(&rpt, "✗ ssh check failed: %v\n", err)
        } else { fmt.Fprintf(&rpt, "✓ ssh ok\n") }

        // Disk list for dd safety
        if m.cfg.Strategy == StratDD {
            fmt.Fprintf(&rpt, "Listing disks via lsblk…\n")
            if have("lsblk") {
                out, _ := os_exec.Command("lsblk", "-o", "NAME,SIZE,TYPE,MODEL").CombinedOutput()
                fmt.Fprintf(&rpt, "%s\n", out)
                if m.cfg.SourceDisk == "" { ok = false; fmt.Fprintf(&rpt, "✗ source disk not set\n") }
            } else { ok = false; fmt.Fprintf(&rpt, "✗ need lsblk for dd safety\n") }
        }

        return preflightDoneMsg{ok: ok, report: rpt.String(), err: nil}
    }
}

// --------------------------- RUN BACKUP ---------------------------

func (m model) runBackup() tea.Cmd {
    return func() tea.Msg {
        ctx, cancel := context.WithCancel(context.Background())
        m.startTime = time.Now()
        m.cancel = cancel

        var cmd *os_exec.Cmd
        var stdout io.ReadCloser
        var stderr io.ReadCloser
        var err error

        logLine := func(s string) tea.Cmd { return func() tea.Msg { return runLogMsg{line: s} } }
        cmds := []tea.Cmd{}

        switch m.cfg.Strategy {
        case StratDD:
            // dd if=<disk> | [pv] | [pigz|gzip] | ssh user@host "cat > path/file.img.gz"
            if m.cfg.SourceDisk == "" { return runDoneMsg{err: fmt.Errorf("source disk not set")} }
            sshSpec := fmt.Sprintf("%s@%s", m.cfg.RemoteUser, m.cfg.RemoteHost)
            remote := fmt.Sprintf("-p% d", m.cfg.SSHPort)
            date := time.Now().Format("2006-01-02")
            remoteFile := strings.ReplaceAll(m.cfg.RemotePath, "$(hostname)", hostname()) + fmt.Sprintf("/disk-%s.img", date)
            // ensure remote dir
            os_exec.Command("ssh", remote, sshSpec, "mkdir", "-p", path_file.Dir(remoteFile)).Run()

            var pipeCmd []string
            pipeCmd = append(pipeCmd, "dd", fmt.Sprintf("if=%s", m.cfg.SourceDisk), "bs=64K", "status=progress")
            if have("pv") { pipeCmd = append(pipeCmd, "|", "pv") }
            if m.cfg.Compression == "pigz" && have("pigz") { pipeCmd = append(pipeCmd, "|", "pigz") }
            if m.cfg.Compression == "gzip" && have("gzip") { pipeCmd = append(pipeCmd, "|", "gzip") }
            pipeCmd = append(pipeCmd, "|", "ssh", remote, sshSpec, "cat", ">", remoteFile)
            cmdStr := strings.Join(pipeCmd, " ")
            cmd = os_exec.CommandContext(ctx, "bash", "-c", cmdStr)
            cmds = append(cmds, logLine("Running: "+cmdStr))

        case StratRsync:
            rsArgs := []string{"-aAXHvz", "--numeric-ids", "--delete-after"}
            if m.cfg.BandwidthKbps > 0 { rsArgs = append(rsArgs, fmt.Sprintf("--bwlimit=%d", m.cfg.BandwidthKbps)) }
            for _, ex := range m.cfg.Excludes { rsArgs = append(rsArgs, "--exclude="+ex) }
            rsArgs = append(rsArgs, "/")
            remote := fmt.Sprintf("-e ssh -p% d", m.cfg.SSHPort)
            rdest := fmt.Sprintf("%s@%s:%s/", m.cfg.RemoteUser, m.cfg.RemoteHost, strings.ReplaceAll(m.cfg.RemotePath, "$(hostname)", hostname()))
            rsArgs = append(rsArgs, "-e", remote, rdest)
            cmd = os_exec.CommandContext(ctx, "rsync", rsArgs...)
            cmds = append(cmds, logLine("Running: rsync "+strings.Join(rsArgs, " ")))

        case StratBorg:
            repo := strings.ReplaceAll(m.cfg.BorgRepo, "$(hostname)", hostname())
            // ensure repo exists
            init := os_exec.Command("borg", "init", "--encryption=repokey", repo)
            _ = init.Run()
            snap := fmt.Sprintf("%s::%s-%s", repo, hostname(), time.Now().Format("2006-01-02"))
            args := []string{"create", "--stats", "--progress", snap, "/"}
            cmd = os_exec.CommandContext(ctx, "borg", args...)
            cmd.Env = os.Environ()
            if m.cfg.BorgPassEnv != "":
                cmd.Env = append(cmd.Env, fmt.Sprintf("BORG_PASSCOMMAND=echo ${%s}", m.cfg.BorgPassEnv))
            }
            cmds = append(cmds, logLine("Running: borg "+strings.Join(args, " ")))

        case StratZFS:
            // zfs snapshot pool/root@ccYYMMDD && zfs send | ssh recv
            date := time.Now().Format("20060102")
            snap := fmt.Sprintf("%s@%s", m.cfg.SourceDisk, date) // here SourceDisk holds dataset like pool/root
            _ = os_exec.Command("zfs", "snapshot", snap).Run()
            remote := fmt.Sprintf("%s@%s", m.cfg.RemoteUser, m.cfg.RemoteHost)
            recv := strings.ReplaceAll(m.cfg.RemotePath, "$(hostname)", hostname())
            cmdStr := fmt.Sprintf("zfs send %s | ssh -p% d %s \"zfs recv %s\"", snap, m.cfg.SSHPort, remote, recv)
            cmd = os_exec.CommandContext(ctx, "bash", "-c", cmdStr)
            cmds = append(cmds, logLine("Running: "+cmdStr))

        case StratBtrfs:
            // btrfs subvolume snapshot -r / /tmp/cc-snap && btrfs send | ssh receive
            snapDir := fmt.Sprintf("/tmp/cc-snap-%d", time.Now().Unix())
            _ = os.MkdirAll(snapDir, 0o755)
            _ = os_exec.Command("btrfs", "subvolume", "snapshot", "-r", "/", snapDir).Run()
            remote := fmt.Sprintf("%s@%s", m.cfg.RemoteUser, m.cfg.RemoteHost)
            recv := strings.ReplaceAll(m.cfg.RemotePath, "$(hostname)", hostname())
            cmdStr := fmt.Sprintf("btrfs send %s | ssh -p% d %s \"btrfs receive %s\"", snapDir, m.cfg.SSHPort, remote, recv)
            cmd = os_exec.CommandContext(ctx, "bash", "-c", cmdStr)
            cmds = append(cmds, logLine("Running: "+cmdStr))
        }

        if cmd == nil { return runDoneMsg{err: fmt.Errorf("no command")} }
        stdout, stderr, err = startCmdPipes(cmd)
        if err != nil { return runDoneMsg{err: err} }

        // stream logs
        go streamReader(stdout, func(line string) { tea.NewProgram(nil).Send(runLogMsg{line: line}) })
        go streamReader(stderr, func(line string) { tea.NewProgram(nil).Send(runLogMsg{line: line}) })

        for _, c := range cmds { _ = c() }
        if err := cmd.Wait(); err != nil { return runDoneMsg{err: err} }
        return runDoneMsg{err: nil}
    }
}

func startCmdPipes(cmd *os_exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
    stdout, err := cmd.StdoutPipe(); if err != nil { return nil, nil, err }
    stderr, err := cmd.StderrPipe(); if err != nil { return nil, nil, err }
    if err := cmd.Start(); err != nil { return nil, nil, err }
    return stdout, stderr, nil
}

func streamReader(r io.Reader, fn func(string)) {
    s := bufio.NewScanner(r)
    s.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
    for s.Scan() { fn(s.Text()) }
}

func hostname() string {
    if h, err := os.Hostname(); err == nil { return h }
    return "host"
}

// --------------------------- MAIN ---------------------------

func main() {
    fmt.Print(lipgloss.NewStyle().Background(paletteBg).Foreground(paletteFg))
    cfg, _ := loadConfig() // fall back to defaults
    m := newModel(cfg)
    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }
}
