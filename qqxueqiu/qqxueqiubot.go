package main

import (
	//	"bytes"
	// "code.google.com/p/go.net/publicsuffix"
	"encoding/json"
	"errors"
	. "fmt"
	//"github.com/PuerkitoBio/goquery"
	//	htmler "html"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	//	"sort"
	"crypto/md5"
	"encoding/base64"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	//	"unicode/utf8"
)

var logfile *os.File = nil

type qquser struct {
	qq   int
	nick string
}
type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	defer jar.lk.Unlock()
	//	Println("SetCookies: ", u.Host, cookies)
	s := jar.cookies[u.Host]
	for i := 0; i != len(cookies); i++ {
		find := false
		for j := 0; j < len(s); j++ {
			if cookies[i].Name == s[j].Name {
				find = true
				s[j] = cookies[i]
				break
			}
		}
		if find == false {
			s = append(s, cookies[i])
		}
	}

	jar.cookies[u.Host] = s

}

//var glock *sync.Mutex

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) (cookies []*http.Cookie) {
	jar.lk.Lock()
	defer jar.lk.Unlock()
	//	Println("GetCookies: ", u.Host)
	for _, v := range jar.cookies {
		for k := 0; k != len(v); k++ {
			bu := false
			for i := 0; i != len(cookies); i++ {
				if cookies[i].Name == v[k].Name {
					if cookies[i].Value == "" {
						cookies[i] = v[k]
						bu = true
						break
					}
				}
			}
			if bu == false {
				cookies = append(cookies, v[k])
			}
		}

		//	cookies = append(cookies, v...)
	}

	return
	//	return jar.cookies[u.Host]
}

var gheader map[string]string

//var gCurCookies []*http.Cookie
var gcookieJar *cookiejar.Jar //*Jar
var ghttpClient *http.Client
var ClientID int

func redirectaddheader(req *http.Request, via []*http.Request) error {
	if len(via) > 10 {
		return Errorf("%d consecutive requests(redirects)", len(via))
	}
	if len(via) == 0 {

		return nil
	}

	// mutate the subsequent redirect requests with the first Header
	for key, val := range via[0].Header {
		req.Header[key] = val
	}
	return nil
}

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}
func init() {
	//	glock = new(sync.Mutex)
	gheader = map[string]string{
		"Accept":          "application/javascript, */*;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.101 Safari/537.36",
		"Accept-Encoding": "identity",
		"Connection":      "close",
	}
	//	gcookieJar = NewJar()
	// options := cookiejar.Options{
	// 	PublicSuffixList: publicsuffix.List,
	// }
	transport := http.Transport{
		Dial: dialTimeout,
	}
	_ = transport
	gcookieJar, _ = cookiejar.New(nil)
	ghttpClient = &http.Client{
		Jar: gcookieJar,
		//	CheckRedirect: redirectaddheader,
		//		Transport: &transport,
		//		Timeout:   time.Duration(5 * time.Second),
	}
	rand.Seed(time.Now().UnixNano())
	ClientID = rand.Intn(888888-111111) + 111111
}
func addheads(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}

func Get(url, refer string) ([]byte, error) {
	//	glock.Lock()
	//	defer glock.Unlock()

	if len(url) == 0 {
		return nil, errors.New("url is empty")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	addheads(req, gheader)
	if len(refer) > 0 {
		req.Header.Add("Referer", refer)
	}

	rsp, err := ghttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//defer rsp.Body.Close()

	b, _ := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	return b, nil
}

func Post(url, data, refer string) ([]byte, error) {
	//	glock.Lock()
	//	defer glock.Unlock()

	log.Println("Post url: ", url)

	if len(url) == 0 {
		return nil, errors.New("url is empty")
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	addheads(req, gheader)
	if len(refer) > 0 {
		req.Header.Add("Referer", refer)
	}

	rsp, err := ghttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//	defer rsp.Body.Close()

	b, _ := ioutil.ReadAll(rsp.Body)
	log.Println("Post rst: ", string(b))
	rsp.Body.Close()
	return b, nil
}
func down(url, path string) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	req, _ := http.NewRequest("GET", url, nil)
	addheads(req, gheader)
	rsp, err := ghttpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	b, _ := ioutil.ReadAll(rsp.Body)

	_, err = f.Write(b)

	return err

}
func getcookie(key string) string {
	u, err := url.Parse("http://qq.com")
	if err != nil {
		return ""
	}
	cooks := gcookieJar.Cookies(u)
	for _, v := range cooks {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}
func printcookie(website string) {
	u, _ := url.Parse(website)

	cooks := gcookieJar.Cookies(u)
	for _, v := range cooks {
		Println("cookie: ", v)
	}

}

var PTWebQQ string = ""
var APPID string = ""
var msgId int = 0

var FriendList map[int]qquser = make(map[int]qquser)
var GroupList map[int]int = make(map[int]int)

var ThreadList []string = make([]string, 0)
var GroupThreadList []string = make([]string, 0)
var GroupWatchList []string = make([]string, 0)

var PSessionID string = ""

var Referer string = "http://d.web2.qq.com/proxy.html?v=20130916001&callback=1&id=2"
var SmartQQUrl string = "http://w.qq.com/login.html"

var VFWebQQ string = ""
var AdminQQ string = "0"

// 聊天授权 qq号
type chatqq map[int]int

var gauthchatqq chatqq = chatqq(make(map[int]int))

func (g chatqq) check(qq int) bool {
	if _, ok := g[qq]; ok {
		return true
	} else {
		return false
	}
}
func (g chatqq) add(qq int) {
	g[qq] = qq
}

// http://tool.chinaz.com/Tools/MD5.aspx?q=xhbot2015032816&md5type=1
// xhbot2015032816
func check_chatmd5(checkmd5 string) bool {
	checkmd5 = strings.TrimSpace(checkmd5)
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	timestr := Sprintf("%d%02d%02d%02d", year, month, day, hour)
	plain := Sprintf("xhbot%s", timestr)

	tarmd5 := Sprintf("%x", md5.Sum([]byte(plain)))
	Println(tarmd5, checkmd5)
	tarmd5 = tarmd5[:6]
	Println(tarmd5, checkmd5)
	if tarmd5 == checkmd5 {
		return true
	} else {
		return false
	}

}

type info struct {
	msgid    int
	msgtxt   string
	msgtitle string
}

type groupmsg struct {
	tuin     int
	msgtxt   string
	msgtitle string
}

var begtime int64 = time.Now().Unix()

func starttime() int64 {
	begtime = time.Now().Unix()
	return begtime
}
func passtime() int64 {
	cost := time.Now().Unix() - begtime
	begtime = time.Now().Unix()
	return cost
}

func getReValue(html []byte, rex string) (string, error) {
	re, err := regexp.Compile(rex)
	if err != nil {
		return "", err
	}
	matchs := re.FindSubmatch(html)
	if len(matchs) > 1 {
		return string(matchs[1]), nil
	}
	return "", errors.New("FindSubmatch did not match")
}

func getqqname(tuin int) (int, string, error) {
	ret, exist := FriendList[tuin]
	if exist {
		if ret.qq == 0 {
			qq, err := getqq(tuin)
			if err != nil {
				return 0, "", err
			} else {
				ret.qq = qq
				FriendList[tuin] = ret
			}
		}
		if ret.nick == "" {
			nick, err := getname(tuin)
			if err != nil {
				return 0, "", err
			} else {
				ret.nick = nick
				FriendList[tuin] = ret
			}
		}
		return FriendList[tuin].qq, FriendList[tuin].nick, nil
	}

	qq, err := getqq(tuin)
	if err != nil {
		return 0, "", errors.New("getqq err")
	}
	nick, err := getname(tuin)
	if err != nil {
		return 0, "", errors.New("getname err")
	}

	FriendList[tuin] = qquser{
		qq:   qq,
		nick: nick,
	}
	return qq, nick, nil
}
func getname(tuin int) (name string, err error) {

	urlpath := Sprintf("http://s.web2.qq.com/api/get_friend_info2?tuin=%d&vfwebqq=%s&clientid=%d&psessionid=%s", tuin, VFWebQQ, ClientID, PSessionID)
	html, err := Get(urlpath, Referer)
	if err != nil {
		return "", err
	}
	//  {"retcode":0,"result":{"face":285,"birthday":{"month":0,"year":0,"day":0},"occupation":"","phone":"","allow":1,"college":"","uin":1959520491,"constel":0,"blood":0,"homepage":"","stat":20,"vip_info":0,"country":"","city":"","personal":"","nick":"Ten.Ten","shengxiao":0,"email":"","province":"","gender":"unknown","mobile":"-"}}

	info := make(map[string]interface{})
	err = json.Unmarshal(html, &info)
	if err != nil {
		return "", err
	}
	if retcode, exist := info["retcode"]; exist {
		retcodeint := int(info["retcode"].(float64))
		if retcodeint != 0 {
			return "", errors.New(Sprintf("retcode != 0. %d", retcode))
		}

		if result, exist := info["result"]; exist {
			if userinfo, ok := result.(map[string]interface{}); ok {
				if nickif, exist := userinfo["nick"]; exist {
					if nick, ok := nickif.(string); ok {

						return nick, nil
					}
				}
			}
		}
	}
	return "", errors.New("json err")
}

func getqq(tuin int) (int, error) {

	urlpath := Sprintf("http://s.web2.qq.com/api/get_friend_uin2?tuin=%d&type=1&vfwebqq=%s", tuin, VFWebQQ)
	html, err := Get(urlpath, Referer)
	if err != nil {
		return 0, err
	}
	// {"retcode":0,"result":{"uiuin":"","account":1982141,"uin":2859374232}}

	info := make(map[string]interface{})
	err = json.Unmarshal(html, &info)
	if err != nil {
		return 0, err
	}
	if retcode, exist := info["retcode"]; exist {
		retcodeint := int(info["retcode"].(float64))
		if retcodeint != 0 {
			return 0, errors.New(Sprintf("retcode != 0. %d", retcode))
		}

		if result, exist := info["result"]; exist {
			if userinfo, ok := result.(map[string]interface{}); ok {
				if account, exist := userinfo["account"]; exist {
					if accountnum, ok := account.(float64); ok {

						return int(accountnum), nil
					}
				}
			}
		}
	}
	return 0, errors.New("json err")
}

func getaireplay(qq int, msg string) (string, error) {
	// http://www.tuling123.com/openapi/api?key=ec930a04c591f56d26a2c4fde771221c&info=helloworld&userid=1111
	urlpath := Sprintf("http://www.tuling123.com/openapi/api?key=%s&info=%s&userid=%d", "ec930a04c591f56d26a2c4fde771221c",
		url.QueryEscape(msg), qq)
	req, err := http.NewRequest("GET", urlpath, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept-Encoding", "text/json; charset=utf-8")
	transport := http.Transport{
		Dial: dialTimeout,
	}

	cli := &http.Client{
		Transport: &transport,
	}

	rsp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	html, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return "", nil
	}
	// {"code":100000,"text":"请不要放弃治疗"}
	info := make(map[string]interface{})
	err = json.Unmarshal(html, &info)
	if err != nil {
		return "", err
	}
	code := 0
	text := ""

	if codeif, exist := info["code"]; exist {
		code = int(codeif.(float64))
	}

	if textif, exist := info["text"]; exist {
		text = textif.(string)
	}
	if code == 100000 {
		{

			return text, nil
		}
	} else if code == 305000 {
		return Sprintf("%s.额,解析好麻烦,你打12360呢.", text), nil
	} else if code == 306000 {
		return Sprintf("%s.额,解析好麻烦,http://www.ctrip.com 上看看呢.", text), nil
	}

	return "", errors.New(Sprintf("code: %d,text: %s", code, text))
}

var gLoginName string = "loginqqname"

func login() error {

	html, err := Get(SmartQQUrl, "")
	if err != nil {
		return err
	}

	//	Println(string(html))
	initurl, err := getReValue(html, `\.src = "(.+?)"`)
	if err != nil {
		return err
	}

	html, err = Get(initurl+"0", "")
	if err != nil {
		return err
	}
	APPID, err = getReValue(html, `var g_appid =encodeURIComponent\("(\d+)"\);`)
	sign, err := getReValue(html, `var g_login_sig=encodeURIComponent\("(.+?)"\);`)
	JsVer, err := getReValue(html, `var g_pt_version=encodeURIComponent\("(\d+)"\);`)
	MiBaoCss, err := getReValue(html, `var g_mibao_css=encodeURIComponent\("(.+?)"\);`)

	path := "./td.jpg"

	err = down(Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid=%s&e=0&l=L&s=8&d=72&v=4", APPID), path)
	if err != nil {
		return err
	}

	Println("登陆二维码下载成功，请扫描")
	log.Println("登陆二维码下载成功，请扫描")

	logincode := -1
	var loginret []string
	maxtrytime := 150 //5分钟
	for t := 0; t < maxtrytime; t++ {
		html, err = Get(Sprintf("https://ssl.ptlogin2.qq.com/ptqrlogin?webqq_type=10&remember_uin=1&login2qq=1&aid=%s&u1=http%%3A%%2F%%2Fw.qq.com%%2Fproxy.html%%3Flogin2qq%%3D1%%26webqq_type%%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-%d&mibao_css=%s&t=undefined&g=1&js_type=0&js_ver=%s&login_sig=%s",
			APPID, t*1000*2, MiBaoCss, JsVer, sign), initurl)
		strthml := string(html)
		loginret = strings.Split(strthml, "'")

		if len(loginret) < 12 {
			return errors.New(Sprintf("%s len < 2", strthml))
		}
		if loginret[1] == "65" || loginret[1] == "0" { // 65: QRCode 失效, 0: 验证成功, 66: 未失效, 67: 验证中
			logincode, _ = strconv.Atoi(loginret[1])
			break
		}
		time.Sleep(2 * time.Second)
		Println("等待扫描...")
		log.Println("等待扫描...")
	}

	if logincode != 0 {
		return errors.New(Sprintf("login ret: %d", logincode))
	}

	os.Remove(path)

	gLoginName = loginret[11]

	Println("已扫描,正在登录")
	log.Println("已扫描,正在登录")

	html, err = Get(loginret[5], "")
	if err != nil {
		log.Println("Get Fail loginret[5]: ", loginret[5])
		return err
	}

	strurl, _ := getReValue(html, `src="(.+?)"`)
	if len(strurl) != 0 {
		html, err = Get(strings.Replace(strurl, "&amp;", "&", -1), "")
		if err != nil {
			return err
		}
		strurl, err = getReValue(html, `location\.href="(.+?)"`)
		html, err = Get(strurl, "")
		if err != nil {
			return err
		}
	}

	PTWebQQ = getcookie("ptwebqq")

	strr := Sprintf("{\"ptwebqq\":\"%s\",\"clientid\":%d,\"psessionid\":\"%s\",\"status\":\"online\"}",
		PTWebQQ, ClientID, PSessionID)
	data := url.Values{"r": {strr}}

	html, err = Post("http://d.web2.qq.com/channel/login2", data.Encode(), Referer)
	if err != nil {
		return err
	}

	info := make(map[string]interface{})
	err = json.Unmarshal(html, &info)
	if err != nil {
		return err
	}
	if retcode, exist := info["retcode"]; exist {
		retcodeint := int(info["retcode"].(float64))
		if retcodeint != 0 {
			return errors.New(Sprintf("retcode != 0. %d", retcode))
		}

		if result, exist := info["result"]; exist {
			if userinfo, ok := result.(map[string]interface{}); ok {
				if vfwebqq, exist := userinfo["vfwebqq"]; exist {
					VFWebQQ = vfwebqq.(string)
				}
				if psessionid, exist := userinfo["psessionid"]; exist {
					PSessionID = psessionid.(string)
				}
				uin, exist := userinfo["uin"]
				uinid := int(uin.(float64))
				if exist {
					Println(Sprintf("QQ号：%d 登陆成功, 用户名：%s", uinid, gLoginName))
					log.Println(Sprintf("QQ号：%d 登陆成功, 用户名：%s", uinid, gLoginName))
					msgId = rand.Intn(50000000+20000000) - 20000000
					return nil
				}
			}

		}
	}
	return errors.New("shit")
}

func sendpersonmsg(tuin int, msg string) error {
	if len(msg) == 0 {
		return errors.New("msg too short")
	}
	msgId++
	log.Println("sendpersonmsg: ", tuin, msg)
	strfmt := `{"to":%d,"content":"[\"%s\",[\"font\",{\"name\":\"Arial\",\"size\":10,\"style\":[0,0,0],\"color\":\"000000\"}]]","face":6,"clientid":%d,"msg_id":%d,"psessionid":"%s"}`
	strr := Sprintf(strfmt, tuin, msg, ClientID, msgId, PSessionID)
	data := url.Values{"r": {strr}}
	html, err := Post("http://d.web2.qq.com/channel/send_buddy_msg2", data.Encode(), Referer)
	if err != nil {
		log.Println(string(html))
	}
	return err
}
func sendgroupmsg(tuin int, msg, color string) error {
	if len(msg) == 0 {
		return errors.New("msg too short")
	}
	msgId++
	if len(color) == 0 {
		color = COLOR[COR_BLACK]
	}
	log.Println("sendgroupmsg: ", tuin, msg)
	strfmt := `{"group_uin":%d,"content":"[\"%s\",[\"font\",{\"name\":\"Arial\",\"size\":10,\"style\":[0,0,0],\"color\":\"%s\"}]]","face":6,"clientid":%d,"msg_id":%d,"psessionid":"%s"} `
	strr := Sprintf(strfmt, tuin, msg, color, ClientID, msgId, PSessionID)
	data := url.Values{"r": {strr}}
	html, err := Post("http://d.web2.qq.com/channel/send_qun_msg2", data.Encode(), Referer)

	if err != nil {
		log.Println(string(html))
	}
	return err
}
func handletxt(txtif []interface{}) string {
	txt := ""
	for itemidx := 1; itemidx < len(txtif); itemidx++ {
		if stritem, ok := txtif[itemidx].(string); ok {
			txt += stritem
		} else if picitemif, ok := txtif[itemidx].([]interface{}); ok {
			if len(picitemif) > 1 {
				if strpicitem, ok := picitemif[0].(string); ok {
					txt += Sprintf("<%s>", strpicitem)
				}
			}

		}
	}
	return txt
}

func handlemsg(msgif interface{}) error {
	if arr, ok := msgif.([]interface{}); ok {
		for _, v := range arr {
			if msg, ok := v.(map[string]interface{}); ok {
				if msgtypeif, exist := msg["poll_type"]; exist {
					if msgtype, ok := msgtypeif.(string); ok {
						if msgtype == "message" || msgtype == "sess_message" {

							valif := msg["value"].(map[string]interface{})
							txtif, _ := valif["content"].([]interface{})
							if len(txtif) < 2 {
								continue
							}
							txt := handletxt(txtif)

							fromid := int(valif["from_uin"].(float64))
							fromqq, fromname, _ := getqqname(fromid)

							Println("person message: ", fromid, fromname, fromqq, txt)
							log.Println("person message: ", fromid, fromname, fromqq, txt)
							if gauthchatqq.check(fromid) == false {
								if strings.Contains(txt, "@"+gLoginName) {
									wakeupmd5, err := getReValue([]byte(txt), `\s*@xhbot\s*wakeup\s*(\w{6})\s*`)
									if err != nil {
										sendpersonmsg(fromid, "姿势不对")
										continue
									} else {
										if check_chatmd5(wakeupmd5) {
											gauthchatqq.add(fromid)
											sendpersonmsg(fromid, "好吧,我听见你心中动人的天籁,登上天外云霄的舞台,嘿!")
										}
									}

								} else {
									continue
								}
								continue
							}

							aireply, err := getaireplay(fromid, txt)
							Println("getaimsg", aireply, err)
							log.Println("getaimsg", aireply, err)

							if err == nil {
								sendpersonmsg(fromid, aireply)
							} else {
								sendpersonmsg(fromid, "我大脑短路咯 ⊙﹏⊙!")
							}

						} else if msgtype == "group_message" {
							valif := msg["value"].(map[string]interface{})
							txtif, _ := valif["content"].([]interface{})
							if len(txtif) < 2 {
								continue
							}
							txt := handletxt(txtif)
							from_uin := int(valif["from_uin"].(float64)) // 群session id
							info_seq := int(valif["info_seq"].(float64)) // 群号
							send_uin := int(valif["send_uin"].(float64)) // 发送人session id

							GroupList[from_uin] = info_seq

							if xindagroupv, bexist := gqqfroup[info_seq]; bexist {
								if from_uin != xindagroupv {
									guserxinda.add_groupuins_byfile(info_seq, from_uin)
									gqqfroup[info_seq] = from_uin
								}
							}

							send_qq, send_name, _ := getqqname(send_uin)

							Println("group message: ", from_uin, info_seq, send_uin, send_qq, send_name, txt)
							log.Println("group message: ", from_uin, info_seq, send_uin, send_qq, send_name, txt)

							strxueqiuid, err := getReValue([]byte(txt), `\s*@xhbot\s*recv\s*(\d+)\s*`)
							if err == nil {
								log.Println(strxueqiuid)
								if xueqiuid, err := strconv.Atoi(strxueqiuid); err == nil {
									straddret := guserxinda.add_groupuins(xueqiuid, info_seq, from_uin)
									sendgroupmsg(from_uin, straddret, COLOR[COR_87843b])
								}
							} else {
								if gauthchatqq.check(send_uin) == false {
									if strings.Contains(txt, "@"+gLoginName) {
										wakeupmd5, err := getReValue([]byte(txt), `\s*@xhbot\s*wakeup\s*(\w{6})\s*`)
										if err != nil {
											sendgroupmsg(from_uin, "姿势不对!", COLOR[COR_YELLOW])
											continue
										} else {
											if check_chatmd5(wakeupmd5) {
												gauthchatqq.add(send_uin)
												sendgroupmsg(from_uin, "好吧,我听见你心中动人的天籁,登上天外云霄的舞台,嘿!", COLOR[COR_87843b])
											}
										}
									} else {
										continue
									}
								} else {
									if strings.Contains(txt, "@"+gLoginName) {
										txt = strings.Replace(txt, "@"+gLoginName, "", -1)
										aireply, err := getaireplay(send_qq, txt)
										if err == nil {
											sendgroupmsg(from_uin, aireply, COLOR[COR_BLACK])
										} else {
											sendgroupmsg(from_uin, "谁叫我 (⊙０⊙)!", COLOR[COR_YELLOW])
										}
									}
								}

							}

						} else {
							Println("un handle message type. ", msgtype)
							log.Println("un handle message type. ", msgtype, msg)
							continue
						}
					} else {
						continue
					}
				} else {
					continue
				}
			} else {
				continue
			}

		}
	} else {
		return errors.New("msg is not array")
	}
	return nil
}
func checkmsg() error {

	strr := Sprintf("{\"ptwebqq\":\"%s\",\"clientid\":%d,\"psessionid\":\"%s\",\"key\":\"\"}",
		PTWebQQ, ClientID, PSessionID)
	data := url.Values{"r": {strr}}
	html, err := Post("http://d.web2.qq.com/channel/poll2", data.Encode(), Referer)
	if err != nil {
		return err
	}

	info := make(map[string]interface{})
	err = json.Unmarshal(html, &info)
	if err != nil {
		return err
	}
	if retcodeif, exist := info["retcode"]; exist {
		retcode := int(retcodeif.(float64))
		switch retcode {
		case 100006:
			return errors.New(Sprintf("retcode: %d. %s", retcode, string(html)))
		case 102:
			return nil
		case 116:
			if pif, exist := info["p"]; exist {
				PTWebQQ = pif.(string)
				log.Println("update PTWebQQ ", PTWebQQ)
				return nil
			} else {
				return errors.New(Sprintf("p dont exist. %s", string(html)))
			}
		case 0:
			if rstif, exist := info["result"]; exist {
				return handlemsg(rstif)
			} else {
				return errors.New(Sprintf("result dont exist. %s", string(html)))
			}
			//	case 121:
			//		return errors.New("账户被登出")
		default:
			log.Println("checkmsg unhandle code", string(html))
			return errors.New(Sprintf("unhandle code. %s", string(html)))
		}
	}
	return errors.New(Sprintf("retcode dont exist. %s", string(html)))
}

type Timelinejson struct {
	Ttotal    int                      `json:"total"`
	TmaxPage  int                      `json:"maxPage"`
	Tcount    int                      `json:"count"`
	Tstatuses []map[string]interface{} `json:"statuses"`
	Tpage     int                      `json:"page"`
}
type timelinemsg struct {
	id                int
	title             string
	created_at        int
	commentId         int
	retweet_status_id int
	description       string
	edited_at         int
	text              string //
	source            string
}

type xueqiuportfoliostockjson struct {
	Scomment   string  `json:"comment"`
	SsellPrice float64 `json:"sellPrice"`
	SbuyPrice  float64 `json:"buyPrice"`
	ScreateAt  int     `json:"createAt"`
	SstockName string  `json:"stockName"`
}

func (s xueqiuportfoliostockjson) String() string {
	return Sprintf("股票:%s,买入价:%.2f,卖出价:%.2f,添加时间:%s,备注:%s", s.SstockName, s.SbuyPrice, s.SsellPrice,
		time.Unix(int64(s.ScreateAt/1000), 0).String(), s.Scomment)
}
func (s xueqiuportfoliostockjson) FormatDiff(n xueqiuportfoliostockjson) string {
	buyprice := Sprintf("%0.2f", s.SbuyPrice)
	sellprice := Sprintf("%0.2f", s.SsellPrice)
	comment := s.Scomment

	if s.SbuyPrice != n.SbuyPrice {
		buyprice = Sprintf("%0.2f(修改前:%0.2f)", n.SbuyPrice, s.SbuyPrice)
	}
	if s.SsellPrice != n.SsellPrice {
		sellprice = Sprintf("%0.2f(修改前:%0.2f)", n.SsellPrice, s.SsellPrice)
	}

	if s.Scomment != n.Scomment {
		comment = Sprintf("%s(修改前:%s)", n.Scomment, s.Scomment)
	}

	return Sprintf("股票:%s,买入价:%s,卖出价:%s,添加时间:%s,备注:%s", s.SstockName, buyprice, sellprice,
		time.Unix(int64(s.ScreateAt/1000), 0).String(), comment)
}

type xueqiuportfoliojson struct {
	SisPublic bool                       `json:"isPublic"`
	Scount    int                        `json:"count"`
	Sstocks   []xueqiuportfoliostockjson `json:"stocks"`
}

const (
	COR_BLACK int = iota
	COR_BLUE
	COR_YELLOW
	COR_RED
	COR_f15b6c
	COR_6b473c
	COR_87843b
)

var COLOR map[int]string = map[int]string{
	COR_BLACK:  "293047",
	COR_RED:    "f15a22",
	COR_BLUE:   "76becc",
	COR_YELLOW: "e0861a",
	COR_f15b6c: "f15b6c",
	COR_6b473c: "6b473c",
	COR_87843b: "87843b",
}

const (
	MSG_ADD int = iota
	MSG_EDIT
	MSG_DEL
	MSG_STOCK_ADD
	MSG_STOCK_EDIT
	MSG_STOCK_DEL
	MSG_COMMENT
)

var MSG_HEAD map[int]string = map[int]string{
	MSG_ADD:        "新发表:",
	MSG_EDIT:       "编辑了:",
	MSG_DEL:        "删除了:",
	MSG_STOCK_ADD:  "增加自选:",
	MSG_STOCK_EDIT: "编辑自选:",
	MSG_STOCK_DEL:  "删除自选:",
	MSG_COMMENT:    "评论了：",
}

type qqiduin struct {
	id  int
	uin int
}
type notifymsg struct {
	itype     int // 0 增加,1编辑,2删除  ,当为MSG_COMMENT 特殊处理
	msg       timelinemsg
	groupuins []qqiduin
	name      string
	uid       int
}

type tlmslice []timelinemsg

func (p tlmslice) Len() int           { return len(p) }
func (p tlmslice) Less(i, j int) bool { return p[i].id > p[j].id }
func (p tlmslice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type xueqiuusr struct {
	lock      *sync.Mutex
	uid       int
	name      string
	actives   map[int]timelinemsg
	groupid   int
	stocks    map[string]xueqiuportfoliostockjson
	groupuins []qqiduin // 通过 @xhbot  recv  2821861040  regexp:^@xhbot\s*recv\s*(\d+)$
}

var guserxinda *xueqiuusr

var gxuueqiuusrsendqueue chan (notifymsg) = make(chan notifymsg, 100)

var gxindaemail []string = make([]string, 0)
var gqqfroup map[int]int = make(map[int]int)
var EmailUsr = "xxxx@163.com"
var EmailPwd = "xxxx"
var EmailSrv = "smtp.163.com:25"

var gxueqiucli *http.Client
var gxueqiucookie *cookiejar.Jar

func init() {
	guserxinda = &xueqiuusr{
		new(sync.Mutex),
		2821861040,
		//3893368287, // 测试id
		"炒的是心",
		make(map[int]timelinemsg),
		0,
		make(map[string]xueqiuportfoliostockjson),
		make([]qqiduin, 0),
	}
}

func read_email_qq() error {
	f, err := os.Open("./xinemail.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	emails, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	gxindaemail = strings.Split(string(emails), ",")

	fq, err := os.Open("./qqgroup.txt")
	if err != nil {
		return err
	}
	defer fq.Close()
	qs, err := ioutil.ReadAll(fq)
	if err != nil {
		return err
	}

	qqarr := strings.Split(string(qs), ",")

	for _, v := range qqarr {
		group, _ := strconv.Atoi(v)
		gqqfroup[group] = 0
	}
	return nil
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}
func SendMail(from, to, subject, body string) error {

	hp := strings.Split(EmailSrv, ":")
	auth := smtp.PlainAuth("", EmailUsr, EmailPwd, hp[0])

	subject = erasebractket(subject)
	tmpsubject := subject

	subidx := 0
	for len(tmpsubject) > 0 {
		_, size := utf8.DecodeRuneInString(tmpsubject)
		if subidx+size > 70 {
			break
		}
		subidx += size
		tmpsubject = tmpsubject[size:]
	}

	subject = subject[:subidx]

	log.Println("SendMail", from, subject, to)

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	msg := []byte("To: " + to + "\r\nFrom: " + encodeRFC2047(from) + "<" + EmailUsr + ">\r\nSubject: " + encodeRFC2047(subject) + "\r\n" + "MIME-Version: 1.0\r\n" + "Content-Type: text/html; charset=UTF-8 \r\nContent-Transfer-Encoding: base64\r\n\r\n" + b64.EncodeToString([]byte(body)))
	send_to := strings.Split(to, "|")
	err := smtp.SendMail(EmailSrv, auth, EmailUsr, send_to, msg)
	return err
}

func SendToEvernote(from, to, subject, body string) error {

	hp := strings.Split(EmailSrv, ":")
	auth := smtp.PlainAuth("", EmailUsr, EmailPwd, hp[0])

	subject = erasebractket(subject)
	tmpsubject := subject

	subidx := 0
	for len(tmpsubject) > 0 {
		_, size := utf8.DecodeRuneInString(tmpsubject)
		if subidx+size > 70 {
			break
		}
		subidx += size
		tmpsubject = tmpsubject[size:]
	}

	subject = subject[:subidx]

	subject = strings.Replace(subject, "@", "", -1)
	subject += " @股市"

	log.Println("SendToEvernote", from, subject, to)

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	msg := []byte("To: " + to + "\r\nFrom: " + "xhbot" + "<" + EmailUsr + ">\r\nSubject: " + encodeRFC2047(subject) + "\r\n" + "MIME-Version: 1.0\r\n" + "Content-Type: text/html; charset=UTF-8 \r\nContent-Transfer-Encoding: base64\r\n\r\n" + b64.EncodeToString([]byte(body)))
	send_to := strings.Split(to, "|")
	err := smtp.SendMail(EmailSrv, auth, EmailUsr, send_to, msg)
	return err
}

func sendmailoreverote(msg notifymsg) error {

	subject := ""
	if len(msg.msg.title) > 0 {
		subject = Sprintf("(动态id=%d,被引用动态id=%d)%s %s %s", msg.msg.id, msg.msg.retweet_status_id, msg.name, MSG_HEAD[msg.itype], msg.msg.title)
	} else {
		subject = Sprintf("(动态id=%d,被引用动态id=%d)%s %s %s", msg.msg.id, msg.msg.retweet_status_id, msg.name, MSG_HEAD[msg.itype], msg.msg.description)
	}

	tomailarr := make([]string, 0)
	toevernotearr := make([]string, 0)

	for _, v := range gxindaemail {
		if strings.Contains(v, "yinxiang.com") {
			toevernotearr = append(toevernotearr, v)
		} else {
			tomailarr = append(tomailarr, v)
		}
	}
	var err error = nil
	if len(tomailarr) > 0 {
		err = SendMail(msg.name, strings.Join(tomailarr, "|"), subject, msg.msg.text)
		log.Println("SendMail:", strings.Join(tomailarr, "|"), subject, err)
	}

	if len(toevernotearr) > 0 {
		err = SendToEvernote(msg.name, strings.Join(toevernotearr, "|"), subject, msg.msg.text)
		log.Println("SendToEvernote:", strings.Join(toevernotearr, "|"), subject, err)
	}

	return err
}
func xueqiugettimelineroution() {

	//	rand.Seed(time.Now().Unix())
	err := getxueqiucookie()
	if err != nil {
		Println("getxueqiucookie: ", err)
		log.Println("getxueqiucookie: ", err)
	}

	for {
		urlstr := Sprintf("http://xueqiu.com/v4/statuses/user_timeline.json?user_id=%d", guserxinda.uid)
		mp, err := getxueqiutimeline(urlstr)
		Println("getxueqiutimeline: ", len(mp), err)
		log.Println("getxueqiutimeline: ", len(mp), err)

		if err == nil {
			guserxinda.handle_newmsg(mp)

			urlstr := Sprintf("http://xueqiu.com/stock/portfolio/stocks.json?size=1000&tuid=%d", guserxinda.uid)
			ms, err := getxueqiuportfolio(urlstr)
			if err == nil {
				guserxinda.handle_stocks(ms)
			}
		}

		if err != nil {
			time.Sleep(time.Minute * 2)
			getxueqiucookie()
		} else {
			time.Sleep(time.Second * 45)
		}
	}
}
func xueqiuusrsendgrouproution() {
	for {
		select {
		case msg := <-gxuueqiuusrsendqueue:
			log.Println("xueqiuusrsendgrouproution: ", msg.name, msg.itype, msg.groupuins)

			var color int = msg.itype

			switch msg.itype {
			case MSG_ADD:
				color = COR_RED
			case MSG_EDIT:
				color = COR_BLUE
			case MSG_DEL:
				color = COR_YELLOW
			}
			for _, v := range msg.groupuins {
				time.Sleep(100 * time.Millisecond)
				if len(msg.msg.title) > 0 {
					sendgroupmsg(v.uin, Sprintf(`%s %s %s( http://xueqiu.com/%d/%d ,被引用id=%d)`, msg.name, MSG_HEAD[msg.itype], trimforqqmsg(erasebractket(msg.msg.title)), msg.uid, msg.msg.id, msg.msg.retweet_status_id), COLOR[color])
				} else {
					sendgroupmsg(v.uin, Sprintf(`%s %s %s( http://xueqiu.com/%d/%d ,被引用id=%d)`, msg.name, MSG_HEAD[msg.itype], trimforqqmsg(erasebractket(msg.msg.description)), msg.uid, msg.msg.id, msg.msg.retweet_status_id), COLOR[color])
				}
			}

			sendmailoreverote(msg)
		case <-time.After(time.Millisecond * 100):

		}
		time.Sleep(time.Second)
	}
}
func (u *xueqiuusr) add_groupuins_byfile(gid, guin int) {
	u.lock.Lock()
	defer u.lock.Unlock()

	bfind := false
	for i := 0; i < len(u.groupuins); i++ {
		if u.groupuins[i].id == gid {
			u.groupuins[i].uin = guin
			bfind = true
			break
		}
	}
	if bfind == false {
		u.groupuins = append(u.groupuins, qqiduin{gid, guin})
	}
	Println("xinda qq atuo update: ", gid)
	log.Println("xinda qq atuo update: ", gid)
}
func (u *xueqiuusr) add_groupuins(xueqiuid, gid, guin int) string {
	log.Println(u.name, "add_groupuins", xueqiuid, gid, guin)
	u.lock.Lock()
	defer u.lock.Unlock()

	if u.uid == xueqiuid {

		bfind := false
		for i := 0; i < len(u.groupuins); i++ {
			if u.groupuins[i].id == gid {
				if u.groupuins[i].uin != guin {
					u.groupuins[i].uin = guin
					bfind = true
					return Sprintf("啫,已更新收听 %s", u.name)
				} else {
					return Sprintf("已收听 %s", u.name)
				}
			}
		}
		if bfind == false {
			u.groupuins = append(u.groupuins, qqiduin{gid, guin})
		}
		return Sprintf("OK,开始收听 %s", u.name)
	}
	return "目前不能收听该人"

}

func (u *xueqiuusr) handle_stocks(socksmap map[string]xueqiuportfoliostockjson) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if len(u.stocks) == 0 {
		u.stocks = socksmap
		return
	}
	if len(socksmap) == 0 {
		return
	}
	tmp := timelinemsg{
		id:                0,
		title:             "",
		created_at:        0,
		commentId:         0,
		retweet_status_id: 0,
		description:       "",
		edited_at:         0,
		text:              "",
		source:            "雪球",
	}
	for k, vold := range u.stocks {
		if vnew, ok := socksmap[k]; ok {
			if vold != vnew {
				strmsg := vold.FormatDiff(vnew)
				Println("stock edit:", strmsg)
				log.Println("stock edit:", strmsg)

				tmp.created_at = vnew.ScreateAt
				tmp.description = strmsg
				tmp.edited_at = vnew.ScreateAt
				tmp.text = strmsg
				tmp.title = strmsg

				ntfmsg := notifymsg{MSG_STOCK_EDIT, tmp, append(make([]qqiduin, 0), u.groupuins...), u.name, u.uid}
				gxuueqiuusrsendqueue <- ntfmsg
				// 被修改的
			}
		} else {
			// 被删
			Println("stock del:", vold)
			log.Println("stock del:", vold)

			tmp.created_at = vold.ScreateAt
			tmp.description = vold.String()
			tmp.edited_at = vold.ScreateAt
			tmp.text = vold.String()
			tmp.title = vold.String()

			ntfmsg := notifymsg{MSG_STOCK_DEL, tmp, append(make([]qqiduin, 0), u.groupuins...), u.name, u.uid}
			gxuueqiuusrsendqueue <- ntfmsg
		}
	}

	for k, vnew := range socksmap {
		if _, ok := u.stocks[k]; ok == false {
			Println("stock add:", vnew)
			log.Println("stock add:", vnew)

			tmp.created_at = vnew.ScreateAt
			tmp.description = vnew.String()
			tmp.edited_at = vnew.ScreateAt
			tmp.text = vnew.String()
			tmp.title = vnew.String()

			ntfmsg := notifymsg{MSG_STOCK_ADD, tmp, append(make([]qqiduin, 0), u.groupuins...), u.name, u.uid}
			gxuueqiuusrsendqueue <- ntfmsg
		}
	}

	u.stocks = socksmap
}
func (u *xueqiuusr) handle_newmsg(msgmap map[int]timelinemsg) {
	u.lock.Lock()
	defer u.lock.Unlock()

	if len(u.actives) == 0 {
		u.actives = msgmap
		return
	}

	for k, _ := range u.actives {
		if v, ok := msgmap[k]; ok {
			if v.edited_at > u.actives[k].edited_at {
				Println(k, v.edited_at, u.actives[k].edited_at)
				ntfmsg := notifymsg{MSG_EDIT, v, append(make([]qqiduin, 0), u.groupuins...), u.name, u.uid}
				gxuueqiuusrsendqueue <- ntfmsg
				// 被修改的
			}
		} else {
			for tmpmixk, _ := range msgmap {
				if tmpmixk < k { // 如果在最新的检索中发现比待检的id更小,说明被删了.
					ntfmsg := notifymsg{MSG_DEL, u.actives[k], append(make([]qqiduin, 0), u.groupuins...), u.name, u.uid}
					gxuueqiuusrsendqueue <- ntfmsg
					break
				}
			}
		}
	}

	for k, _ := range msgmap {
		if _, ok := u.actives[k]; ok == false {
			bbig := true
			for tmpmaxk, _ := range u.actives {
				if tmpmaxk > k { // 如果在老的集合中没有找到比这个id大的说明就是新增的
					bbig = false
					break
				}
			}
			if bbig {
				ntfmsg := notifymsg{MSG_ADD, msgmap[k], append(make([]qqiduin, 0), u.groupuins...), u.name, u.uid}
				gxuueqiuusrsendqueue <- ntfmsg
			}
		}
	}
	u.actives = msgmap
}

func getxueqiucookie() error {
	gxueqiucookie, _ := cookiejar.New(nil)
	gxueqiucli = &http.Client{

		Jar: gxueqiucookie,
	}

	req, err := http.NewRequest("GET", "http://xueqiu.com/", nil)
	if err != nil {
		return err
	}

	//	req.Header.Add("X-Forwarded-For", Sprintf("%d.%d.%d.%d", rand.Intn(60)+192, rand.Intn(250), rand.Intn(250), rand.Intn(250)))
	//	log.Println(req.Header.Get("X-Forwarded-For"))
	rsp, err := gxueqiucli.Do(req)
	if err != nil {
		return err
	}
	rsp.Body.Close()

	time.Sleep(time.Second)
	req, err = http.NewRequest("GET", "http://xueqiu.com/2821861040", nil)
	if err != nil {
		return err
	}

	rsp, err = gxueqiucli.Do(req)
	if err != nil {
		return err
	}
	rsp.Body.Close()

	return nil
}

// 获取关注stock状态
// http://xueqiu.com/stock/portfolio/stocks.json?size=1000&pid=-1&tuid=2821861040&cuid=3893368287&_=1427343673358
// http://xueqiu.com/stock/portfolio/stocks.json?tuid=2821861040

func getxueqiuportfolio(urlstr string) (map[string]xueqiuportfoliostockjson, error) {
	Println("getxueqiuportfolio:", urlstr)
	log.Println("getxueqiuportfolio:", urlstr)
	socks := make(map[string]xueqiuportfoliostockjson)

	req, _ := http.NewRequest("GET", urlstr, nil)

	rsp, err := gxueqiucli.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	jsontxt, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	portfolio := &xueqiuportfoliojson{}
	err = json.Unmarshal(jsontxt, portfolio)
	if err != nil {
		log.Println("json.Unmarshal err: ", string(jsontxt))
		return nil, err
	}
	for _, v := range portfolio.Sstocks {

		socks[v.SstockName] = v
	}
	Println("getxueqiuportfolio:", len(socks))
	log.Println("getxueqiuportfolio:", len(socks))
	return socks, nil
}

func getxueqiutimelinecomment(activeid int) error {
	return nil
}
func getxueqiutimeline(urlstr string) (map[int]timelinemsg, error) {

	msgmap := make(map[int]timelinemsg)

	ipage := 0
	imaxpage := 2

	blatest := true

	for blatest && ipage < imaxpage {
		time.Sleep(time.Second * 5)
		tmpurl := ""
		ipage++
		tmpurl = Sprintf("%s&page=%d", urlstr, ipage)
		Println(tmpurl)
		log.Println(tmpurl)
		req, _ := http.NewRequest("GET", tmpurl, nil)

		rsp, err := gxueqiucli.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		jsontxt, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		timeline := &Timelinejson{}
		err = json.Unmarshal(jsontxt, timeline)
		if err != nil {
			log.Println("json.Unmarshal err: ", string(jsontxt))
			return nil, err
		}
		if imaxpage > timeline.TmaxPage {
			imaxpage = timeline.TmaxPage
		}
		for _, v := range timeline.Tstatuses {
			id := 0
			if idif, ok := v["id"]; ok {
				if idif == nil {
					continue
				} else {
					id = int(idif.(float64))
				}
			} else {
				continue
			}

			title := ""
			if titleif, ok := v["title"]; ok {
				if titleif != nil {
					title = titleif.(string)
				}
			}

			created_at := 0
			if created_atif, ok := v["created_at"]; ok {
				if created_atif != nil {
					created_at = int(created_atif.(float64))
				}
			}

			timecreate := time.Unix(int64(created_at)/1000, 0)

			if time.Now().Sub(timecreate) > time.Hour*24*30 {
				blatest = false
				//	Println("第 ", ipage, "页的 ", title, " 以后的文章太旧了.")
				break
			}

			commentId := 0
			if commentIdif, ok := v["commentId"]; ok {
				if commentIdif == nil {
					continue
				} else {
					commentId = int(commentIdif.(float64))
				}
			} else {
				continue
			}
			retweet_status_id := 0
			if retweet_status_idif, ok := v["retweet_status_id"]; ok {
				if retweet_status_idif != nil {
					retweet_status_id = int(retweet_status_idif.(float64))
				}
			}

			description := ""
			if descriptionif, ok := v["description"]; ok {
				if descriptionif != nil {
					description = descriptionif.(string)
				}
			}

			edited_at := 0
			if edited_atif, ok := v["edited_at"]; ok {
				if edited_atif != nil {
					edited_at = int(edited_atif.(float64))
				}
			}

			text := ""
			if textif, ok := v["text"]; ok {
				if textif != nil {
					text = textif.(string)
				}
			}

			source := "雪球"
			if sourceif, ok := v["source"]; ok {
				if sourceif != nil {
					source = sourceif.(string)
				}
			}

			msg := timelinemsg{
				id,
				title,
				created_at,
				commentId,
				retweet_status_id,
				description,
				edited_at,
				text,
				source,
			}

			msgmap[msg.id] = msg
		}
	}
	return msgmap, nil
}
func trimforqqmsg(s string) string {

	tmpsubject := s

	subidx := 0
	for len(tmpsubject) > 0 {
		_, size := utf8.DecodeRuneInString(tmpsubject)
		if subidx+size > 300 {
			break
		}
		subidx += size
		tmpsubject = tmpsubject[size:]
	}

	s = s[:subidx]
	return s
}
func erasebractket(s string) string {
	idxl := strings.Index(s, "<")
	idxr := strings.Index(s, ">")

	ret := s
	for idxr != -1 && idxr != -1 {

		face := s[idxl:idxr]

		facetxt, err := getReValue([]byte(face), `title="(\S+)"`)

		ret = s[0:idxl]
		if err == nil {
			ret += facetxt
		}
		ret += s[idxr+1:]

		s = ret

		idxl = strings.Index(s, "<")
		idxr = strings.Index(s, ">")

	}

	nbsp := "&nbsp;"
	lennbsp := len(nbsp)

	idxnbsp := strings.Index(s, nbsp)

	for idxnbsp != -1 {

		ret = s[0:idxnbsp]
		ret += s[idxnbsp+lennbsp:]

		s = ret

		idxnbsp = strings.Index(s, nbsp)

	}
	return ret
}

func main() {
	logfile, _ = os.Create("log.txt")
	defer logfile.Close()
	log.SetOutput(logfile)

	err := read_email_qq()
	log.Println("mails: ", gxindaemail, err)
	log.Println("qq group: ", gqqfroup)
	if err != nil {
		return
	}
	err = login()
	if err == nil {
		go xueqiugettimelineroution()
		go xueqiuusrsendgrouproution()
		for {
			log.Println("checkmsg err:", checkmsg())
			time.Sleep(time.Second * 1)
		}
	}
	log.Println("login err: ", err)

}
