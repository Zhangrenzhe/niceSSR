package main

import (
	"flag"
	"github.com/robfig/cron"
	"github.com/va-len-tine/niceSSR/utils"
	"log"
	"os/exec"
	"time"
)

var ch = make(chan string, 1)
var cmd *exec.Cmd
var ss string
var err error
var SSPath = "ss.txt"
var SSUrl = "https://bulink.me/sub/mruxq/ss"

func main()  {
	var testUrl string
	var testInterval string
	var port string
	flag.StringVar(&testUrl, "u", "", "测速网址")
	flag.StringVar(&testInterval, "t", "15", "自动测速时间间隔/分钟")
	flag.StringVar(&port, "p", "10808", "代理端口")
	flag.Parse()

	var Interval = "0 */" + testInterval +" * * * *"
	if testUrl != "" {
		utils.TestUrl = utils.TestUrl[0:0]
		utils.TestUrl = append(utils.TestUrl, testUrl)
		log.Printf("测速网址：%s 测速间隔：%sMin\n", utils.TestUrl, testInterval)
	}else{
		log.Printf("使用默认测速地址，测速间隔：%sMin\n", testInterval)
	}

	log.Printf("正在获取代理...")
	ss,err = utils.DfShadowsocks.GetAvailSS(SSUrl, SSPath)
	if err != nil {
		log.Fatal(err)
	}

	cmd = utils.DfShadowsocks.NewSock5Proxy(ss, port)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer cmd.Process.Kill()
	log.Printf("代理启动成功！端口:%s PID:%d %s\n\n", port, cmd.Process.Pid, ss)

	// 开启定时任务自动切换代理
	cron2 := cron.New()
	err = cron2.AddFunc(Interval, test)
	if err != nil {
		log.Fatal("添加定时任务失败！")
	}
	cron2.Start()
	defer cron2.Stop()
	for {
		select {
		case ss = <-ch:
			cmd.Process.Kill()
			cmd = utils.DfShadowsocks.NewSock5Proxy(ss ,port)
			cmd.Start()
			time.Sleep(time.Second*1)
			log.Printf("自动切换代理！端口:%s PID:%d %s\n\n", port, cmd.Process.Pid, ss)
		}
	}
}

func test()  {
	log.Println("开启自动测速...")
	ss,err = utils.DfShadowsocks.GetFastSS(SSUrl, SSPath)
	if err != nil {
		log.Fatal(err)
	}
	ch <- ss
}
