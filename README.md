采集Go进程的内存指标绘制成图，可通过网页查看。
可以生成火焰图，可以查看函数分配内存情况。

use:
```
import "github.com/lzxm160/gomemanalysis/core"

core.Start()
```
build webui:

```./build.sh```

explore:
```shell
./webui
```
`http://<yourhostname>:8081` //for go pprof

`http://<yourhostname>`

todo：
1、加入火焰图
2、加入函数内存分析