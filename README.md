# line-cli

A command-line interface for the [LINE Messaging API](https://developers.line.biz/en/reference/messaging-api/) and LINE Login.

Designed to be used by humans and AI agents (Claude Code, Codex, etc.) alike.

## Using with AI agents (Claude Code, Codex, etc.)

Paste this into your project's `CLAUDE.md`, `AGENTS.md`, or system prompt to give your AI agent full access to LINE:

````markdown
## LINE CLI

`line` is installed and configured. Use it to interact with the LINE Messaging API.

### Setup (if not already done)
```bash
curl -fsSL https://raw.githubusercontent.com/emailnitram/line-cli/main/install.sh | sh
line config set --token YOUR_ACCESS_TOKEN
```

### Key commands
```bash
line message push --to USER_ID --text "Hello"          # Send a message to a user
line message broadcast --text "Hello everyone"          # Send to all followers
line message push --to USER_ID --file /tmp/msg.json    # Send flex/rich message from file
line profile get --user-id USER_ID                     # Get user profile
line richmenu list                                      # List rich menus
line richmenu upload-image --id ID --file img.png      # Upload image (auto-resizes)
line webhook get                                        # Get current webhook URL
line webhook set --url https://...                     # Set webhook URL
line coupon list                                        # List coupons
line --help                                            # Full command list
line COMMAND --help                                    # Help for a specific command
```

### Output
All commands output JSON (pretty-printed) or a short confirmation string on success.
Errors are written to stderr with a non-zero exit code.

### Sending rich/flex messages
Write the messages array to a temp file, then pass it with --file:
```bash
cat > /tmp/msg.json << 'EOF'
[{"type": "text", "text": "Hello from LINE CLI"}]
EOF
line message push --to USER_ID --file /tmp/msg.json
```

### Notes
- Credentials are stored in ~/.line-cli/config.yaml
- Rich menu images are auto-resized and compressed — any image file works
````

---

## Installation

### curl (macOS, Linux, CI, Claude Code sandboxes)
```bash
curl -fsSL https://raw.githubusercontent.com/emailnitram/line-cli/main/install.sh | sh
```

Detects your OS and architecture automatically. Installs to `/usr/local/bin` (or `~/.local/bin` if permissions require it). No Go or Homebrew needed.

### Homebrew (macOS/Linux)
```bash
brew install emailnitram/tap/line-cli
```

### Go install
```bash
go install github.com/emailnitram/line-cli@latest
```

## Authentication

Get your credentials from the [LINE Developers Console](https://developers.line.biz/console/) under your channel → **Messaging API** tab.

```bash
# If you already have a channel access token
line config set --token YOUR_ACCESS_TOKEN

# If you only have channel credentials
line config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET
line token issue
line config set --token THE_TOKEN_FROM_ABOVE

# Check what's saved
line config show
```

Credentials are stored in `~/.line-cli/config.yaml`.

## Commands

### Token
```bash
line token issue                          # Issue a short-lived channel access token
line token verify                         # Verify the configured token
line token revoke                         # Revoke the configured token
```

### Message
```bash
line message push --to USER_ID --text "Hello"
line message push --to USER_ID --file flex.json       # Send from JSON file
line message broadcast --text "Hello everyone"
line message broadcast --file messages.json
line message reply --reply-token TOKEN --text "Hi"
line message quota                                     # Get monthly quota and usage
line message loading --chat-id USER_ID --seconds 5    # Show loading animation
```

### LIFF
```bash
line liff list
line liff create --url https://example.com/liff --type full
line liff create --file liff.json                           # Full config from file
line liff update --id LIFF_ID --url https://example.com/v2
line liff update --id LIFF_ID --type tall
line liff delete --id LIFF_ID
```

**View types:** `compact` (50% screen height), `tall` (80%), `full` (100%)

### Profile
```bash
line profile get --user-id USER_ID
line profile demographics                 # Get follower age/gender/region breakdown
```

### Rich Menu
```bash
line richmenu list
line richmenu create --file menu.json
line richmenu delete --id RICHMENU_ID
line richmenu set-default --id RICHMENU_ID
line richmenu cancel-default
line richmenu link --id RICHMENU_ID --user-id USER_ID
line richmenu get-user --user-id USER_ID
line richmenu create-alias --alias-id richmenu-a --menu-id RICHMENU_ID

# Upload image — auto-resizes and compresses to fit LINE's 1MB limit
line richmenu upload-image --id RICHMENU_ID --file image.png
line richmenu upload-image --id RICHMENU_ID --file image.png --size 2500x1686
```

**Valid image sizes (px):** `2500x1686` `2500x843` `1200x810` `1200x405` `800x540` `800x270`

### Webhook
```bash
line webhook get
line webhook set --url https://example.com/webhook
```

### Audience
```bash
line audience delete --id AUDIENCE_ID
line audience rename --id AUDIENCE_ID --name "New Name"
line audience get-shared --id AUDIENCE_ID
```

### Coupon
```bash
line coupon create --file coupon.json
line coupon list
line coupon list --status RUNNING --limit 50
line coupon get --id COUPON_ID
line coupon discontinue --id COUPON_ID
line coupon send --id COUPON_ID                       # Broadcast to all followers
line coupon send --id COUPON_ID --to USER_ID          # Send to specific user
```

### Group
```bash
line group members --group-id GROUP_ID
line group member-profile --group-id GROUP_ID --user-id USER_ID
line group leave --group-id GROUP_ID
```

## JSON file format

For `message push --file` and `message broadcast --file`, pass a JSON array of [message objects](https://developers.line.biz/en/reference/messaging-api/#message-objects):

```json
[
  { "type": "text", "text": "Hello!" },
  { "type": "sticker", "packageId": "446", "stickerId": "1988" }
]
```

## Building from source

```bash
git clone https://github.com/emailnitram/line-cli
cd line-cli
make build      # builds ./line
make install    # installs to $GOPATH/bin
make release    # cross-compiles for all platforms into ./dist
```

## License

MIT
