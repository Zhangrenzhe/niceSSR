# 一键科学上网

### 一.开启本地代理
#### 克隆项目
```
git clone git@github.com:va-len-tine/niceSSR.git
```
#### 进入项目编译，生成可执行程序：OpenSocks5.exe
```
cd nicess/
go build OpenSocks5.go
```
#### 打开OpenSocks5.exe，一键开启代理，默认端口为：10808
![blockchain](http://image.youthsweet.com/nicess/nicess-1.png "default")
#### 
### 二.安装浏览器插件SwitchyOmega，切换代理
#### 安装SwitchyOmega插件
插件在chromePlugin目录下，Chrome浏览器插件安装过程，请自行百度
```
$ ls chromePlugin/
SwitchyOmega_Chromium.crx  gfw.txt
```
#### 配置proxy代理
- 代理协议选择：SOCKS5
- 代理服务器：127.0.0.1
- 代理端口：10808
![blockchain](http://image.youthsweet.com/nicess/nicess-2.png "default")
#### 切换至proxy代理，至此就可以愉快的访问GitHub、YouTube了！
![blockchain](http://image.youthsweet.com/nicess/nicess3.png "default")

#### 解决cmd控制台执行暂停的问题
右键控制台->默认值->关闭(快速编辑模式)
![blockchain](http://image.youthsweet.com/nicess/nicess-4.png "default")
