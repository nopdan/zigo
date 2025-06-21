# zigo

![Actions Check](https://badgen.net/github/checks/nopdan/zigo)
[![Stars](https://badgen.net/github/stars/nopdan/zigo)](https://github.com/nopdan/zigo/stargazers)
![Downloads](https://badgen.net/github/assets-dl/nopdan/zigo)
![License](https://badgen.net/github/license/nopdan/zigo)
[![Release](https://badgen.net/github/release/nopdan/zigo)](https://github.com/nopdan/zigo/releases)

Download and manage Zig compilers.

`zigo` is a dependency free binary program.

> Inspired by [zigup](https://github.com/marler8997/zigup), [zvm](https://github.com/tristanisham/zvm), [scoop](https://github.com/ScoopInstaller/Scoop)

## Install

- Download from [releases](https://github.com/nopdan/zigo/releases/).

```sh
tar -zxvf ./zigo-linux-amd64.tar.gz
chmod +x ./zigo
./zigo
```

- From source: `go install github.com/nopdan/zigo@latest`

Add `~/.zig/current` to the environment variable.

## Usage

Set environment variable `ZIGO_PATH` to change compiler installation location.

**default:** `~/.zig`

### `zigo <version>`

Download the specified version of zig compiler and set it as default.

examples: `zigo 0.11.0`, `zigo master`

```sh
❯ ./zigo master
Downloading zig-x86_64-windows-0.15.0-dev.864+75d0ec9c0.zip...
url: https://ziglang.org/builds/zig-x86_64-windows-0.15.0-dev.864+75d0ec9c0.zip
save to: C:\Users\Admin\AppData\Local\zigo\zig-x86_64-windows-0.15.0-dev.864+75d0ec9c0.zip
progress: 88.90 MiB / 88.90 MiB | 100.0 % | 2.49 MiB/s         
Done.
Using master => 0.15.0-dev.864+75d0ec9c0
❯ zig version
0.12.0-dev.374+742030a8f
```

### `zigo fetch <version>`

Download the specified version of zig compiler.

### `zigo use <version>`

Set the specific installed version as the default.

### `zigo ls`

List installed compiler versions.

```sh
❯ ./zigo ls
* master => 0.12.0-dev.312+cb6201715
  0.11.0
  0.12.0-dev.307+7827265ea
  0.12.0-dev.312+cb6201715
```

### `zigo rm <version>`

Remove the specified compiler.

Append `-f` or `--force` to force deletion.

```sh
❯ ./zigo rm 0.12.1
Removing 0.12.1... 
Done.

❯ ./zigo rm 0.11.0
Cannot remove the version you are using.
❯ ./zigo rm 0.11.0 -f
Removing 0.11.0... 
Done.
```

### `zigo clean`

Clean up unused dev version compilers.

```sh
❯ ./zigo ls
* master => 0.12.0-dev.353+4a44b7993
  0.11.0
  0.12.0-dev.312+cb6201715
  0.12.0-dev.352+4d29b3967
  0.12.0-dev.353+4a44b7993
❯ ./zigo clean
Removing 0.12.0-dev.312+cb6201715...
Removing 0.12.0-dev.352+4d29b3967...
Done.
❯ ./zigo ls
* master => 0.12.0-dev.353+4a44b7993
  0.11.0
  0.12.0-dev.353+4a44b7993
```

### `zigo -h`

Print help message.

## Build

Install `Go 1.21.0+`

```sh
git clone https://github.com/nopdan/zigo.git
cd zigo
go mod tidy
go build
```
