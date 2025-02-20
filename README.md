# Alert Notify
这是一个用于通知的程序，当有告警时，会通过webhook的方式通知到其他系统。

## 编译

```bash
go build -o alert-notify main.go
```

## k8s 部署

```bash
skaffold run --namespace=monitoring
```


