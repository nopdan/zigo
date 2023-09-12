# zigo

GO 语言写的一个 zig 版本管理工具。

灵感来自 [zigup](https://github.com/marler8997/zigup), [zvm](https://github.com/tristanisham/zvm), [scoop](https://github.com/ScoopInstaller/Scoop)

zigo 使用 [grab](https://github.com/cavaliergopher/grab) 下载文件到内存中

zigo 不依赖 `tar`，而是用 [archiver](https://github.com/mholt/archiver) 解压 `.tar.xz` 和 `.zip` 文件。

TODO:

- [ ] 通过命令行修改 `install-dir`，并移动文件夹内容，重新链接 `current`
- [ ] github action

## 编译

需要 golang 1.21.1+

```sh
git clone https://github.com/nopdan/zigo.git
cd zigo
go build
```

## 使用

```sh
Root Command:
  zigo <version>         download and set the compiler as default

Sub Commands:
  list, ls               list installed compiler versions
  remove, rm <version>   remove compiler
  help, -h
```

## 配置

zigo 的配置文件在 `~/.config/zigo.json`，你可以修改其中的 `install-dir` 改变安装路径，默认安装路径是 `~/.zig`
