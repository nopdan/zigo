# zigo

GO 语言写的一个 zig 版本管理工具。

灵感来自 [zigup](https://github.com/marler8997/zigup), [zvm](https://github.com/tristanisham/zvm), [scoop](https://github.com/ScoopInstaller/Scoop)

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
