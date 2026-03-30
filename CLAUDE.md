# line-cli

This is the source repo for `line-cli`. When working here, use `make build` to
compile and `./line` to test commands locally.

## LINE CLI usage (for agents using this tool externally)

`line` is a CLI for the LINE Messaging API. Install and configure it once:

```bash
go install github.com/emailnitram/line-cli@latest
line config set --token YOUR_ACCESS_TOKEN
```

### Sending messages
```bash
line message push --to USER_ID --text "Hello"
line message broadcast --text "Hello everyone"

# Rich/flex messages — write JSON to a temp file first
cat > /tmp/msg.json << 'EOF'
[{"type": "text", "text": "Hello"}]
EOF
line message push --to USER_ID --file /tmp/msg.json
```

### Common operations
```bash
line profile get --user-id USER_ID
line richmenu list
line richmenu upload-image --id RICHMENU_ID --file image.png   # auto-resizes
line webhook get
line webhook set --url https://...
line coupon list
line order get --order-no ORDER_NO
```

### Discovering commands
```bash
line --help
line message --help
line richmenu --help
```

### Output format
All commands output pretty-printed JSON or a short confirmation on success.
Errors go to stderr with a non-zero exit code — check exit code to detect failures.
