使用`typora`，有时需要插入图片，需要用到图床，该项目就是图床工具

# 工具使用

## server模式

server模式会启动一个tcp服务和一个http服务，tcp服务用于接收hangzhou图片，http服务用于提供图片查看

启动服务的命令如下：

``` shell
pichouse server [-c=/path/to/config/file]
```

配置文件的默认路径为：`~/.pichouse`

配置文件的内容如下：

``` json
{
    "picHouse": "/path/to/image/folder",
    "viewImageUrl": "view image base http address",
    "httpPort": "http service port",
    "tcpPort": "tcp service port"
}
```

## client模式

client模式只负责连接server并通过tcp上传图片，client模式的使用方式如下：

``` shell
pichouse client -uu=server_tcp_address /local/file/path
```

# 图片同步

图片同步有三种方式：

- 使用`git`来同步
- 使用`rsync`命令同步，[rsync官网](https://www.samba.org/ftp/rsync/rsync.html)
- 自己开发同步服务

## rsync使用

同步远程图片文件的命令如下：

``` bash
rsync -r -v --ignore-existing user@host:/path/to/folder ./
```

其中：

- -r   表示在文件夹中进行递归
- -v  显示详细信息
- --ignore-existing  忽略本地已存在的文件

可以配合`expect`工具，免去输入密码的步骤，如下：

``` shell
#!/usr/bin/expect

spawn rsync -r -v --ignore-existing user@host:/path/to/folder ./
expect "*password"
send "password\n"
expect eof
```



# 代办

- 增加图片同步的功能