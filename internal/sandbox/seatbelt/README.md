# Seatbelt sandbox for MacOS

## Reason

Obvious way of isolation on macOS (also and linux) is use Docker. 
But for macOS docker brings overhead, like daemon and virtual machine. 
On macOS, we have solution to `sandbox` process using only kernel features - that what `sandbox-exec`

## Implementation

This package provide implementation for sandbox by macOS seatbelt.
It's not well documented by apple, but there doc from [chromium sandbox on macOS](https://source.chromium.org/chromium/chromium/src/+/main:sandbox/mac/seatbelt_sandbox_design.m)

Major part is inspired by [OpenAI Codex](https://github.com/openai/codex) and sbpl configuration files taken from there
and adopted to wide use over classical CLI/TUI applications.

All features from `sandbox` package is supported.

## Useful links

| Name                                                   | Link                                                                           |
|--------------------------------------------------------|--------------------------------------------------------------------------------|
| Doc for rules from Mozilla                             | https://wiki.mozilla.org/Sandbox/OS_X_Rule_Set                                 |
| Caveats from Lucas Wiman on `Sandboxing code on MacOS` | https://lucaswiman.github.io/2023/06/04/macos-sandbox.html                     |
| Apple's Sandbox Guide                                  | https://reverse.put.as/wp-content/uploads/2011/09/Apple-Sandbox-Guide-v1.0.pdf |

## Testing

Tests based on snapshots over generated profile. So, just running test is enough to verify that **only** profile correct.

To update snapshots in tests run
```shell
UPDATE_SNAPSHOTS=1 go test ./internal/sandbox/seatbelt -run Snapshot -v
```