1、获取安装

```
git clone https://github.com/tkeel-io/tkeel-batch-tool.git
go build
```

2、预览

```
tanli@ubuntu:~/Desktop/workspace/project/tkeel-batch-tool$ ./tkeelBatchTool -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  tkeelBatchTool [command]

Available Commands:
  dev         Creat device from excel
  help        Help about any command
  mapper      Creat mapper from excel
  spaceTree   Creat spaceTree from excel
  template    Creat template from excel

Flags:
  -c, --conf string   The iot api config
  -f, --file string   The data excel
  -h, --help          help for tkeelBatchTool
  -o, --op string     add or del

Use "tkeelBatchTool [command] --help" for more information about a command.
```

3、运行前置条件

登录平台

3.1 直接输入
```
> tkeelBatchTool login http://tkeel.io:30080/ --tenant <your tenant name> --username <your username> --password <your password>

✅  You are Login as admin in tenant an!
✅  AccessToken is [eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0a2VlbCIsImV4cCI6MTY1NTk5MDQ4NSwic3ViIjoidXNyLTc3MzUwNWZkZDM1MzY1MTI3NDEyNmJlNmQyMDIifQ.HF9cUKDOmhvHooqwPSH1e3xKnpPNqyfhDlSLQ0aP45nSDptHigu06PbKkjcEEEvmDCMVlZ6wY_P55QZVZnb6Lg]
✅  RefreshToken is [YZG5NME4YJMTZTEXYY01NJJHLWJIMZKTMMY0YTZJOTZHNJMW]
✅  Login Token save in ./config.json!
✅  You are Login Success!
```

3.2 交互式输入
```
>  tkeelBatchTool login http://tkeel.io:30080/

? Please enter your tenant:  an
? Please enter your username:  admin
? Please enter your password:  ******
✅  You are Login as admin in tenant an!
✅  AccessToken is [eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0a2VlbCIsImV4cCI6MTY1NTk5MDU5Miwic3ViIjoidXNyLTc3MzUwNWZkZDM1MzY1MTI3NDEyNmJlNmQyMDIifQ.lsha1wiI-RDB9ZPAbFWarQAyJqRqYsDEkofqy5I9aXVvbPnPKute1hbNcJ5qwQ4AZ85A_BDSvPM7L41Yfi12Sg]
✅  RefreshToken is [ODY4OGM0NTKTODU0ZC01NJC2LWFKNDMTNJE5MDA1MWEYMZEY]
✅  Login Token save in ./config.json!
✅  You are Login Success!
```


4、命令概览：以下命令是一个完整的流程 注意顺序

```
批量新增：
./tkeelBatchTool template -o add -f excel_file/template.xlsx //批量新增模板
./tkeelBatchTool spaceTree -o add -f excel_file/spaceTree.xlsx  //批量新增空间节点（设备组）
./tkeelBatchTool dev -o add -f excel_file/devices.xlsx  //批量新增设备
./tkeelBatchTool mapper -o add -f excel_file/mapper.xlsx        //批量新增设备数据映射关系

批量删除：
./tkeelBatchTool mapper -o del -f excel_file/mapper.xlsx        //批量删除设备数据映射关系 
./tkeelBatchTool dev -o del -f excel_file/devices.xlsx  //批量删除设备
./tkeelBatchTool template -o del -f excel_file/template.xlsx  //批量删除模板
./tkeelBatchTool spaceTree -o del -f excel_file/spaceTree.xlsx  //批量删除空间节点（设备组）

配置文件默认读取./config.json   指定路径 -c xxx/path
```



5、具体细节及excel  格式详细说明见如下链接

https://docs.tkeel.io/developer_cookbook/device/batch/batch_1
