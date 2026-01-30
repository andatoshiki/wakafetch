# Wakafetch

![Showcase](https://cdn.tosh1ki.de/assets/images/20260129180950.png)

A command-line tool to fetch and display your coding stats from WakaTime or Wakapi in your terminal, without opening or refreshing the web dashboard.

Original author: **sahaj-b**. Current maintainer: [**andatoshiki**](https://toshiki.dev).

## 1: Features

- **Quick stats**: Summary of coding activity for configurable time ranges (`--range` or `--days`).
- **Deep dive**: `--full` shows languages, projects, editors, operating systems, and more.
- **Daily breakdown**: `--daily` shows a day-by-day table.
- **Activity heatmap**: `--heatmap` shows a GitHub-style heatmap; use `--range` (e.g. 7d, 30d, 6m, 1y) for a smaller window, or a year (e.g. 2024) for a full year.
- **WakaTime and Wakapi**: Works with the official [WakaTime](https://wakatime.com) API and [Wakapi](https://github.com/muety/wakapi), including self-hosted instances.
- **Zero-config**: Reads API key from `~/.wakatime.cfg`; override with `--api-key` if needed.

## 2: Installation

### 2.1: One-liner (curl)

Install the latest release (macOS, Linux):

```bash
curl -fsSL https://raw.githubusercontent.com/andatoshiki/wakafetch/master/scripts/install.sh | sh
```

Installs to `/usr/local/bin` if writable, otherwise `~/.local/bin`. Override with `INSTALL_DIR`:

```bash
curl -fsSL https://raw.githubusercontent.com/andatoshiki/wakafetch/master/scripts/install.sh | INSTALL_DIR=~/bin sh
```

### 2.2: From source

```bash
git clone https://github.com/andatoshiki/wakafetch.git
cd wakafetch
go build
./wakafetch --help
```

### 2.3: With Go install

Installs the binary to `$GOPATH/bin` or `$GOBIN`:

```bash
go install github.com/andatoshiki/wakafetch@latest
```

## 3: Configuration

wakafetch reads `~/.wakatime.cfg`. Put your API key and optional API URL there (e.g. from your WakaTime or Wakapi settings):

```ini
[settings]
api_key = your-api-key
api_url = https://wakatime.com/api
```

For Wakapi or a self-hosted instance, set `api_url` to your instance (e.g. `https://wakapi.dev/api` or `https://your-server/api`). The tool normalizes non–WakaTime URLs to the compat API path.

If you use the WakaTime editor extension, this config is usually already present.

## 4: Usage

Default view (last 7 days):

```bash
wakafetch
```

Full list of options:

```bash
wakafetch --help
```

Options:

| Flag | Description |
|------|-------------|
| `-r`, `--range` | Range: today, yesterday, 7d, 30d, 6m, 1y, all, or a year (e.g. 2024). Default: yesterday |
| `-d`, `--days` | Number of days (overrides `--range`) |
| `-f`, `--full` | Full statistics |
| `-D`, `--daily` | Daily breakdown table |
| `-H`, `--heatmap` | Activity heatmap; window set by `--range` (7d, 30d, 6m, 1y, or year) |
| `-k`, `--api-key` | Override API key from config |
| `-t`, `--timeout` | Request timeout in seconds (default: 10) |
| `-j`, `--json` | Output JSON |
| `-h`, `--help` | Help |

> [!WARNING]
> **Historic data and `--range`**: The official [WakaTime](https://wakatime.com) API and hosted [Wakapi](https://wakapi.dev) typically require a **Pro/Premium** plan to return summary or historic data for longer time ranges. Using `--range` (e.g. `1y`, `6m`, or a past year) may result in errors or empty results on free tiers. **Self-hosted Wakapi** has no such limit and returns full historic data.

## 5: Examples

- Yesterday (default): `wakafetch`
- Last 7 days: `wakafetch --range 7d`
- Last 30 days: `wakafetch --range 30d`
- Full stats for the last year: `wakafetch -r 1y -f`
- Last 100 days: `wakafetch --days 100`
- Daily breakdown for 2 weeks: `wakafetch --days 14 --daily`
- Heatmap for last 7 days: `wakafetch -H --range 7d`
- Heatmap for last 30 days: `wakafetch -H --range 30d`
- Heatmap for last 12 months: `wakafetch -H` or `wakafetch -H --range 1y`
- Heatmap for a specific year: `wakafetch -H --range 2024`

## 6: License

MIT. See [LICENSE](LICENSE).