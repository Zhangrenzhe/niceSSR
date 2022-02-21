package utils

import (
	"bufio"
	"encoding/base64"
	"errors"
	"github.com/va-len-tine/niceSSR/config"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Shadowsocks struct {
}
var DfShadowsocks = &Shadowsocks{}

var TestPort = "1080"
var TestTimeOut = 5
var TestClient = DfShadowsocks.NewHttpProxyClient(TestPort, TestTimeOut)
var TestUrl = []string{"https://www.youtube.com/", "https://github.com/avelino/awesome-go"}

// NewHttpProxyClient 新建sock5代理client,设置超时时间
func (ssk *Shadowsocks)NewHttpProxyClient(port string, timeout int) *http.Client{
	dialSocksProxy, err := proxy.SOCKS5("tcp", "127.0.0.1:"+port, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}
	return &http.Client{
		Transport: &http.Transport{Dial: dialSocksProxy.Dial},
		Timeout: time.Second * time.Duration(timeout),
	}
}

// NewSock5Proxy 新建一个cmd
func (ssk *Shadowsocks)NewSock5Proxy(ss string, port string) *exec.Cmd {
	return exec.Command(config.ShadowPath,"-c", ss, "-verbose", "-socks", ":"+port, "-u")
}

// OpenSock5Proxy 在本地1080端口开启一个socks5代理
func (ssk *Shadowsocks)OpenSock5Proxy(ss string, port string)(*exec.Cmd, error){
	cmd := exec.Command(config.ShadowPath,"-c", ss, "-verbose", "-socks", ":"+port, "-u")
	err := cmd.Start()
	if err != nil {
		return nil,err
	}
	time.Sleep(time.Second*1)
	return  cmd,nil
	// 调用后要释放资源，防止资源泄露
	//defer cmd.Process.Kill()
}

// TestSS 测速ss
func (ssk *Shadowsocks)TestSS(ssurl string) (float64,error){
	cmd,err := ssk.OpenSock5Proxy(ssurl, TestPort)
	if err != nil {
		return 100,err
	}
	defer cmd.Process.Kill()

	var t float64
	for _,v := range TestUrl {
		t += ssk.TestUrl(v)
	}
	return t,nil
}

// TestUrl 测试单个URL访问时间
func (ssk *Shadowsocks)TestUrl(url string) float64  {
	t1 := time.Now()

	response, err := TestClient.Get(url)
	if err != nil {
		return float64(100)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return float64(100)
	}
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return float64(100)
	}

	return time.Since(t1).Seconds()
}

// GetFastSS 测试全部SS，返回时间最短的一个
func (ssk *Shadowsocks)GetFastSS(url string, path string)(string, error)  {
	result := ""
	//MinTime := float64(TestTimeOut * len(TestUrl))
	var MinTime float64 = 200

	ss,err := ssk.GetAllSS(url,path)
	if err != nil {
		return "",err
	}
	for _,v := range ss {
		t,_ := ssk.TestSS(v)
		log.Printf("%s %.2fs\n", v, t)
		if  t < MinTime {
			MinTime = t
			result = v
		}
	}
	if result == ""{
		return "",errors.New("未找到可用的ss")
	}
	return result,nil
}

// GetAvailSS 找到一个可用的即返回
func (ssk *Shadowsocks)GetAvailSS(url string, path string)(string, error)  {
	result := ""
	MinTime := float64(TestTimeOut * len(TestUrl))

	ss,err := ssk.GetAllSS(url,path)
	if err != nil {
		return "",err
	}
	for _,v := range ss {
		t,_ := ssk.TestSS(v)
		log.Printf("%s %.2fs\n", v, t)
		if  t < MinTime {
			MinTime = t
			result = v
			break
		}
	}
	if result == ""{
		return "",errors.New("未找到可用的ss")
	}
	return result,nil
}

// GetAllSS 获取SS
func (ssk *Shadowsocks)GetAllSS(url string, path string) ([]string,error) {
	s2,err := ssk.GetSSFromUrl(url)
	if err != nil {
		return nil,err
	}
	s1,err := ssk.GetSSFromTxt(path)
	if err != nil {
		return nil,err
	}
	s1 = append(s1, s2...)
	return s1,nil
}

// GetSSFromTxt 本地文本读取ss链接
func (ssk *Shadowsocks)GetSSFromTxt(path string) ([]string, error) {
	var ss []string
	file,err:=os.Open(path)
	if err != nil {
		return nil,err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineStr:=strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(lineStr, "ss://"){
			v1 := strings.Split(lineStr, "#")[0]
			v2,_ := base64.URLEncoding.DecodeString(v1[5:])
			v3 := "ss://" + string(v2)
			ss = append(ss, v3)
		}
	}
	return ss,nil
}

// GetSSFromUrl 通过订阅链接获取ss
func (ssk *Shadowsocks)GetSSFromUrl(subsUrl string) ([]string, error) {
	r,err := http.Get(subsUrl)
	if err != nil {
		return nil,err
	}
	defer r.Body.Close()

	b,err := io.ReadAll(r.Body)
	if err != nil {
		return nil,err
	}

	bb,err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil,err
	}

	var ss []string
	for _,v := range strings.Split(string(bb), "\n"){
		if strings.HasPrefix(v, "ss://"){
			v1 := strings.Split(v, "#")[0]
			v2 := strings.Split(v1, "@")
			v3,_ := base64.RawURLEncoding.DecodeString(v2[0][5:])
			v4 := "ss://" + string(v3) + "@" + v2[1]
			ss = append(ss, v4)
		}
	}
	return ss,nil
}
