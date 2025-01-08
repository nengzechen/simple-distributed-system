# simple-distrubuted-system
使用go构建简单的分布式系统
## 项目内容
服务注册、用户门户、日志服务、业务服务
## 服务注册
- 创建web服务
- 创建注册服务
- 注册web服务
- 取消注册web服务
## 服务发现
- 业务服务
- 服务发现
- 依赖服务变化的通知
## 状态检测

---

### 服务注册
1. registry包用于统一的服务注册，服务注册本身也是一个webservice，他在cmd/registryservice中启动
2. 其他服务在cmd下进行启动，启动时会调用service包，这是一个用于服务启动的包
3. service包会调用registry包中函数用于服务注册，同时service包也可服务注销