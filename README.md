# APX - CLI wrapper for platform-specific sandboxes

## Usage

Run [codex](https://github.com/openai/codex) from OpenAI in sandbox

```shell
# Just run
apx codex
# Pass args
apx -- codex -m "gpt-5"
```

## Supported sandboxes

✅ - Support implemented

❌ - Not implemented, planned 

| Name             | OS          | Status | Description                                                                                                     |
|------------------|-------------|--------|-----------------------------------------------------------------------------------------------------------------|
| Seatbelt         | MacOS       | ✅      | Native macos sandbox via sandbox-exec, supported by kernel natively                                             |
| Landlock+seccomp | Linux       | ❌      | Supported by kernel (5.13+) restrictions for process over files access (via Landlock) and network (via seccomp) |
| Docker           | Linux+MacOS | ❌      | Isolation via docker containers                                                                                 |

## Debugging

When some app or cli fails with permission denied without specific details, os-specific tools can help with debug.


### MacOS

Via `log`:

For `Seatbelt`
```shell
sudo log stream --style compact --info --predicate 'subsystem == "com.apple.sandbox" OR process == "sandboxd" OR eventMessage CONTAINS[c] "deny"'
```

