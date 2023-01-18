package hello

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"gopkg.in/yaml.v2"
)

type HelloMsg struct {
	Msg   string   `yaml:"msg"`
	Image []string `yaml:"image"`
}

type Config struct {
	BlackList []int64  `yaml:"blacklist"`
	Group     HelloMsg `yaml:"group"`
	Friend    HelloMsg `yaml:"friend"`
	Delay     struct {
		Enable bool `yaml:"enable"`
		MaxMin int  `yaml:"max"`
	} `yaml:"delay"`
}

var instance *hello
var logger = utils.GetModuleLogger("com.aimerneige.hello")
var helloConfig Config

type hello struct {
}

func init() {
	instance = &hello{}
	bot.RegisterModule(instance)
}

func (h *hello) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "com.aimerneige.hello",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (h *hello) Init() {
	path := config.GlobalConfig.GetString("aimerneige.hello.path")
	if path == "" {
		path = "./hello.yaml"
	}
	bytes := utils.ReadFile(path)
	if err := yaml.Unmarshal(bytes, &helloConfig); err != nil {
		logger.WithError(err).Errorf("Unable to read config file in %s", path)
	}
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (h *hello) PostInit() {
}

// Serve 注册服务函数部分
func (h *hello) Serve(b *bot.Bot) {
	b.GroupInvitedEvent.Subscribe(func(client *client.QQClient, event *client.GroupInvitedRequest) {
		logger.Infof("「%s」(%d) 邀请机器人加入群「%s」(%d)。", event.InvitorNick, event.InvitorUin, event.GroupName, event.GroupCode)
		if inBlacklist(event.InvitorUin) {
			randomDelay() // 随机等待一会再处理请求
			event.Reject(true, "您已被列入黑名单")
			return
		}
		randomDelay() // 随机等待一会再处理请求
		event.Accept()
		sayHello(client, message.Source{
			SourceType: message.SourceGroup,
			PrimaryID:  event.GroupCode,
		})
	})
	b.NewFriendRequestEvent.Subscribe(func(client *client.QQClient, event *client.NewFriendRequest) {
		logger.Infof("「%s」(%d) 请求添加为好友，验证消息为“%s”。", event.RequesterNick, event.RequesterUin, event.Message)
		if inBlacklist(event.RequesterUin) {
			randomDelay() // 随机等待一会再处理请求
			event.Reject()
			return
		}
		randomDelay() // 随机等待一会再处理请求
		event.Accept()
		sayHello(client, message.Source{
			SourceType: message.SourcePrivate,
			PrimaryID:  event.RequesterUin,
		})
	})
}

// Start 此函数会新开携程进行调用
// ```go
//
//	go exampleModule.Start()
//
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (h *hello) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (h *hello) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}

func sayHello(c *client.QQClient, target message.Source) {
	helloMsg := getHelloMsg(c, target)
	switch target.SourceType {
	case message.SourceGroup:
		c.SendGroupMessage(target.PrimaryID, helloMsg)
	case message.SourcePrivate:
		c.SendPrivateMessage(target.PrimaryID, helloMsg)
	}
}

func getHelloMsg(c *client.QQClient, target message.Source) *message.SendingMessage {
	var helloMsg HelloMsg
	switch target.SourceType {
	case message.SourceGroup:
		helloMsg = helloConfig.Group
	case message.SourcePrivate:
		helloMsg = helloConfig.Friend
	}
	msg := message.NewSendingMessage()
	msg.Append(message.NewText(helloMsg.Msg))
	for _, img := range helloMsg.Image {
		if img == "" {
			continue
		} else {
			imgBytes, err := ioutil.ReadFile(img)
			if err != nil {
				logger.WithError(err).Errorf("Fail to read image %s", img)
				msg.Append(message.NewText("\n\n[图片 - 读取失败]"))
				continue
			}
			imgMsgElement, err := c.UploadImage(target, bytes.NewReader(imgBytes))
			if err != nil {
				logger.WithError(err).Errorf("Fail to upload image %s", img)
				msg.Append(message.NewText("\n\n[图片 - 上传失败]"))
				continue
			}
			msg.Append(imgMsgElement)
		}
	}
	return msg
}

func randomDelay() {
	if helloConfig.Delay.Enable {
		maxMin := helloConfig.Delay.MaxMin
		// 如果配置文件中未指定，则默认为 30
		if helloConfig.Delay.MaxMin <= 0 {
			maxMin = 30
		}
		randomMinutes := rand.Intn(maxMin)
		randomSeconds := rand.Intn(60)
		time.Sleep(time.Minute * time.Duration(randomMinutes))
		time.Sleep(time.Second * time.Duration(randomSeconds))
	}
}

func inBlacklist(userID int64) bool {
	for _, v := range helloConfig.BlackList {
		if userID == v {
			return true
		}
	}
	return false
}
