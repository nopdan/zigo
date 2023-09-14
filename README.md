# zigo

![Actions Check](https://badgen.net/github/checks/nopdan/zigo)
[![Stars](https://badgen.net/github/stars/nopdan/zigo)](https://github.com/nopdan/zigo/stargazers)
![Downloads](https://badgen.net/github/assets-dl/nopdan/zigo)
![License](https://badgen.net/github/license/nopdan/zigo)
[![Release](https://badgen.net/github/release/nopdan/zigo)](https://github.com/studyzy/studyzy/releases)

Download and manage Zig compilers.

`zigo` is a dependency free binary program.

> Inspired by [zigup](https://github.com/marler8997/zigup), [zvm](https://github.com/tristanisham/zvm), [scoop](https://github.com/ScoopInstaller/Scoop)

## Install

- Download from [releases](https://github.com/nopdan/zigo/releases/).

```sh
unzip ./zigo-linux-amd64.zip
chmod +x zigo
./zigo
```

- From source: `go install github.com/nopdan/zigo@latest`

Add `~/.zig/current` to the environment variable.

## Usage

### `zigo <version>`

Download the specified version of zig compiler and set it as default.

examples: `zigo 0.11.0`, `zigo master`

```sh
❯ ./zigo master
Downloading https://ziglang.org/builds/zig-windows-x86_64-0.12.0-dev.352+4d29b3967.zip...
100% |███████████████████████████████████████████████████████████████| (74/74 MB, 16 MB/s)
installing master => 0.12.0-dev.352+4d29b3967
successfully installed
❯ ./zig version
0.12.0-dev.352+4d29b3967
```

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

```sh
❯ ./zigo rm 0.10.1
removed 0.10.1
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
removed 0.12.0-dev.312+cb6201715
removed 0.12.0-dev.352+4d29b3967
❯ ./zigo ls
* master => 0.12.0-dev.353+4a44b7993
  0.11.0
  0.12.0-dev.353+4a44b7993
```

### `zigo mv <install-dir>`

Move the zig installation directory. Default installation directory is `~/.zig`.

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
