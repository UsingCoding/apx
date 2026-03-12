# APX - CLI wrapper for platform-specific sandboxes

## Usage

Run [codex](https://github.com/openai/codex) from OpenAI in sandbox

```shell
# Just run
apx codex
# Pass args
apx -- codex -m "gpt-5"
```

## Install

Via homebrew on Linux/MacOS

```shell
brew install usingcoding/public/apx
```

## Concept and portability

APX is designed for AI agents usage, but suitable for generic every-day cases.
However, there are core concepts how APX works for **all** OS and sandboxes

* Read (like `/etc`, `/usr`) and write (like `/var`, `/tmp`, `/dev`) access for system paths is **always granted**.
Without it most apps will be unable to start or work correctly
* Home: access for specific paths, not for all dir. If set `home` in `sandboxes.policy.fs` APX gives access to home dir.
But only for specific paths. It reads all paths and passes it explicitly to sandbox implementation.
  * Excluded dirs: `.ssh`, `.kube`.
* APX configuration portable for all [supported sandboxes](#supported-sandboxes).
But Linux Landlock there are more limitations (like lack of support `denyPaths`).
Unsupported cases not causes to crash of APX - you can see them in `apx -v` in logs

## apx.toml

Configuration format for app that can be run in APX supported sandboxes

Format:
```toml
# Name of app. Will be runned with sandbox
name = "opencode"

# Sandboxes defined for app. Can be many
[[sandboxes]]
# Sandbox type. Supported types: seatbelt, landlock
type = "seatbelt"

# Policy for sandbox
# Here all rules placed
[sandboxes.policy]

# Env passed to app when runs in sandbox
env = { K = "V" }

[sandboxes.policy.fs]
# Specific paths that should be RO
roPaths = [
  "$HOME/.state",
]
# Specific paths that should be RW
rwPaths = [
  "$HOME/.state",
]
# Specific paths that should be Denied
denyPaths = [
  "$HOME/.state",
]

# Access to $HOME
[sandboxes.policy.fs.home]
# Full $HOME access
allPaths = true
# Do not use default denyList
skipDefaultDenyList = true
# Forbid specific paths under $HOME
denyList = [ "$HOME/.secret" ]
# Access with RW. If rw ommitted - RO access to home
rw = true

[sandboxes.policy.network]
# Revoke network access
deny = true
```

## Supported sandboxes

✅ - Support implemented

❌ - Not implemented, planned 

| Name             | OS          | Status | Description                                                                                                     |
|------------------|-------------|--------|-----------------------------------------------------------------------------------------------------------------|
| Seatbelt         | MacOS       | ✅      | Native macos sandbox via sandbox-exec, supported by kernel natively                                             |
| Landlock+seccomp | Linux       | ✅      | Supported by kernel (5.13+) restrictions for process over files access (via Landlock) and network (via seccomp) |
| Docker           | Linux+MacOS | ❌      | Isolation via docker containers                                                                                 |

## Configuration

### Local Registry and sandboxes

APX has built-in apps in registry defined at [registry](registry)

You can define own collection of `<app>.apx.toml`

* At `$HOME/.config/apx`
* Create `<app>.apx.toml` in format shown above. For example `git.apx.toml` (sandbox for git ¯\_(ツ)_/¯)
* Run `apx -- git`

**Important**

When you have local and built-in `apx.toml` for same app - local **replaces** built-in, **no merge**

## Debugging

When some app or cli fails with permission denied without specific details, os-specific tools can help with debug.


### MacOS

Via `log`:

For `Seatbelt`
```shell
sudo log stream --style compact --info --predicate 'subsystem == "com.apple.sandbox" OR process == "sandboxd" OR eventMessage CONTAINS[c] "deny"'
```

