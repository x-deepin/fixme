# 功能需求

## 问题编写要求
1. 问题描述
2. 问题检测
3. 问题修复

## 客户端
1. 拉取问题数据库
2. 显示问题列表
3. 执行修复
4. 结果汇报

## 服务器
0. 问题描述页面
1. 结果展示页面
2. 结果汇报API
3. 元素下载


# 用法
```
NAME:
   fixme - Fix urgent bugs in deepin and eventually fix itself.

USAGE:
   fixme [global options] command [command options] [arguments...]
   
VERSION:
   0.0.1
   
COMMANDS:
   show		[pid ...]
   fix		pid1 [pid2 ...]
   update	list all knowned problems
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --server, -s "https://fixme.deepin.com"	server url for updating and reporting
   --db, -d "db.json"				database path
   --help, -h					show help
   --version, -v				print the version
```
