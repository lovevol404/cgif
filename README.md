# 下载已经编译后的二进制
在 Release 界面下载对应平台二进制包：
1. cgif_macos 对应 MacOS
2. cgif.exe 对应 Windows
3. cgif_linux 对应 Linux

# 自行编译
* 配置 GO 语言环境（打开 GO MOD模式）
* 以下命令在 Mac下运行： 

go build -o binary/cgif_macos 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o binary/cgif_linux 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o binary/cgif.exe 

# 使用
把可执行文件放到对应的目录下
* 使用样例
```shell script
./cgif_macos  c -f ji.gif -t 8 -b 280

##############################100.000000%(1543) #进度，括号内是完成的帧数

```
* 查看帮助，共用两种功能
```shell script
./cgif_macos -h #查看帮助
NAME:
   Gif转换 - 把Gif进行转换

USAGE:
   cgif_macos [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   create, c   进行黑白转换
   reverse, r  进行颜色反转
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
* 查看对应功能的帮助
```shell script
./cgif_macos c -h # 或者./cgif_macos r -h
NAME:
   cgif_macos create - 进行黑白转换

USAGE:
   cgif_macos create [command options] [arguments...]

OPTIONS:
   --file value, -f value    Gif文件名称
   --black value, -b value   转换为黑色的阈值（源像素点的rgb和超过此值时为黑色，否则为白色） (default: 300)
   --thread value, -t value  线程数 (default: 4)
```
