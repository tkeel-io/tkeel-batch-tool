# Device Mock 

## 创建模板

程序可以通过 API 查询到对应设备的模板信息并创建 MQTT 负载模板
- telemetry   遥测负载模板
- attributes  属性负载模板


使用方式
```
$ ./tkeelBatchTool mock-example
Mock template creat
? Please enter host:  http://preview.tkeel.io:30080/
? Please enter deviceID:  iotd-d33e5f4e-2d99-4aa3-9a1b-45af10bb9c0d
? Please enter Mode:  telemetry
✅  Write mock sample file(telemetry.mock), size 145 byte
```

### 模板类型

模板分为
- telemetry   遥测
- attributes  属性

可以通过键盘箭头选择
```
❯ go run main.go mock-sample
Mock template creat
? Please enter host:  http://preview.tkeel.io:30080/
? Please enter deviceID:  iotd-d33e5f4e-2d99-4aa3-9a1b-45af10bb9c0d
? Please enter Mode:   [Use arrows to move, type to filter]
> telemetry
  attributes
```

### 更新模板

```
$ ./tkeelBatchTool mock-sample
Mock template creat
? Please enter host:  http://preview.tkeel.io:30080/
? Please enter deviceID:  iotd-d33e5f4e-2d99-4aa3-9a1b-45af10bb9c0d
? Please enter Mode:  telemetry
✅  Write mock sample file(telemetry.mock), size 145 byte
```

生成的模板 `telemetry.mock` 如下：
```json
{
	"B0002": "{{randInt 0 10}}.{{randInt 10 99}}",
	"B003": "{{randInt 0 10}}.{{randInt 10 99}}",
	"B1002": "{{randInt 0 10}}.{{randInt 10 99}}"
}
```

### 模板语法

模板语法采用 [Go’s template language](http://golang.org/pkg/text/template/). 并且内置 [sprig](http://masterminds.github.io/sprig/) 模板库。

例如 如下模板：
```json
{
	"B0002": "{{randInt 0 10}}.{{randInt 10 99}}",
	"B003": "{{randInt 10 20}}.{{randInt 10 99}}",
	"B1002": "{{randInt 100 999}}.{{randInt 10 99}}"
}
```

其中， `randInt 10 20` 会渲染为 10 到 20 之间的整数，因此整个模板会渲染为
```json
{
	"B0002": "2.22",
	"B003": "13.45",
	"B1002": "456.87"
}
```

## 模拟设备


```
$ ./tkeelBatchTool mock
Mock device SendData
? Please enter broker:  preview.tkeel.io:31883
? Please enter deviceID:  iotd-d33e5f4e-2d99-4aa3-9a1b-45af10bb9c0d
? Please enter device Token:  Y2M3YTk2NzMtNmZkYS0zMDU0LWEwMDEtNzE3ZjBiNzNhOWQ5
? Please enter Mode:  telemetry
? Please enter template file:  telemetry.mock
Create MQTT client and connection: uri=tcp://iotd-d33e5f4e-2d99-4aa3-9a1b-45af10bb9c0d:Y2M3YTk2NzMtNmZkYS0zMDU0LWEwMDEtNzE3ZjBiNzNhOWQ5@preview.tkeel.io:31883 clientID=IncomingDataPublisher 
[1]Send
[v1/devices/me/telemetry]: 
 {
        "B0002": "7.85",
        "B003": "7.78",
        "B1002": "0.29"
}
[2]Send
[v1/devices/me/telemetry]: 
 {
        "B0002": "2.10",
        "B003": "0.43",
        "B1002": "5.37"
}
```




## 参考

- [1] Go’s template language (http://golang.org/pkg/text/template/)
- [2] sprig 模板库 (http://masterminds.github.io/sprig/) 
