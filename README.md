## 线报推送脚本
* 线报接口：http://new.xianbao.fun/plus/json/push.json?230406
* 微信推送：https://www.pushplus.plus

---

### 使用说明
* 注册pushplus账号，获取token，token添加到pushToken中
* 启动  ```go  run main.go```

---

### 添加规则，筛选掉不想要的线报
````go
    // 例：推送时间
    func (p *XBPush) pushTimeRule() bool {
        begin := 8
        end := 23
        
        if time.Now().Hour() <= end && time.Now().Hour() >= begin {
        return true
        }
        return false
    }
````
