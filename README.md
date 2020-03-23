采集Go进程的各项内存指标（包含Go runtime内存管理相关的，以及RSS等），并按时间维度绘制成折线图，可以通过网页实时查看。配合原有Go pprof工具，可以快速监控和分析Go进程的内存使用情况。

```golang
import "github.com/lzxm160/gomemanalysis/core"

core.Start()
```

```shell
./webui
```
`http://<yourhostname>:8081` //for go pprof
`http://<yourhostname>`

todo：
1、加入火焰图
2、加入函数内存分析