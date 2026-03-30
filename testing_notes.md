# CLI Testing Notes: `line`

## Observations
- The CLI has several top-level commands: `audience`, `config`, `coupon`, `group`, `message`, `order`, `profile`, `richmenu`, `token`, `webhook`.
- `line config show` reveals four expected configuration keys:
    - `access_token`
    - `client_id`
    - `client_secret`
    - `api_key`
- Currently, all configuration values are "(not set)".
- `line profile get` requires a `--user-id` flag.
- `line config show` masks the `access_token` (e.g., `dumm...oken`).
- The error message for missing token correctly suggests `line config set --token <token>`.

## Summary of Investigation
The `line` CLI provides a broad interface for many LINE APIs, including messaging, rich menus, coupons, and orders. It has a robust configuration system but suffers from some inconsistencies in flag naming, subcommand availability (especially `list` and `unlink`), and output formatting. The error handling is mostly direct pass-through of HTTP responses, which provides technical detail but could be more user-friendly.

### Key Improvement Areas
1. **Consistency**: Standardize flag names across similar commands (e.g., `--id` vs `--group-id` vs `--order-no`).
2. **Missing Features**: Implement `list` subcommands for `audience` and `order`, and `unlink` for `richmenu`.
3. **Machine-readable Output**: Add `--json` support for commands like `config show`.
4. **Validation**: Add local validation for flags like `loading --seconds` to save unnecessary API calls.
5. **Masking**: Standardize sensitive data masking in `config show`.
6. **Error Handling**: Parse and format API error responses for better readability.
