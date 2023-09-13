# zigo

`zigo` 是用 golang 写的一个 [zig](https://ziglang.org/) 版本管理工具。

`zigo` 是一个无依赖的二进制文件。

- 使用 [grab](https://github.com/cavaliergopher/grab) 下载文件。
- 使用 [archiver](https://github.com/mholt/archiver) 解压 `.tar.xz` 和 `.zip` 文件。

> 灵感来自 [zigup](https://github.com/marler8997/zigup), [zvm](https://github.com/tristanisham/zvm), [scoop](https://github.com/ScoopInstaller/Scoop)

## 使用

### `zigo <version>`

下载指定版本的 zig 编译器，并设为默认值。

examples: `zigo 0.11.0`, `zigo master`

### `zigo ls`

列出所有已安装的版本。

```sh
❯ zigo ls
* master => 0.12.0-dev.312+cb6201715
  0.11.0
  0.12.0-dev.307+7827265ea
  0.12.0-dev.312+cb6201715
```

### `zigo rm <version>`

删除指定的版本。

### `zigo mv <install-dir>`

移动 zig 安装路径，默认安装路径是 `~/.zig`。

### `zigo -h`

打印帮助信息。

## 编译

需要 golang 1.21.1+

```sh
git clone https://github.com/nopdan/zigo.git
cd zigo
go build
```
