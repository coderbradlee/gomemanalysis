网页显示go程序内存使用情况

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