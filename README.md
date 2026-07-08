# ts-proxy-command

`ts-proxy-command` is an SSH `ProxyCommand` helper that uses Tailscale's
`tsnet` library to dial a host inside your tailnet.

It is useful when you want a single SSH command to bring up a small embedded
Tailscale node and proxy the SSH TCP stream through it.

## Build

```sh
nix build
```

The binary is written to:

```sh
./result/bin/ts-proxy-command
```

## Usage

Use it with SSH:

```sh
ssh -o "ProxyCommand=$PWD/result/bin/ts-proxy-command %h %p" yourhost
```

Or add it to `~/.ssh/config`:

```sshconfig
Host yourhost
  ProxyCommand /path/to/ts-proxy-command %h %p
```

## Tailscale State

By default, the embedded Tailscale state is stored in:

```sh
$HOME/.local/state/tsnet-ssh-proxy
```

Set `TSNET_STATE_DIR` to use a different state directory.

On first use, `tsnet` may print a Tailscale login URL. Complete that login once;
subsequent runs should reuse the stored state.

## Notes

This program proxies the SSH TCP connection only. Tools like Mosh use SSH for
setup but then switch to UDP, so they need a separate mosh-aware wrapper or UDP
relay if the target is only reachable through this proxy.

## License

MIT
