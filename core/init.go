package core

import (
	"bytes"
	"config"
	"controller"
	"core/library"
	"net/http"
	"net/url"
	"reflect"
	"routes"
	"strings"
)

var RegisterMessage = make(map[string]interface{})
var routeList = routes.Web() //加载路由

// LoadRouteHttp 加载控制器函数
func LoadRouteHttp(w http.ResponseWriter, r *http.Request) {
	// 设置所有用户都能访问
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//解析请求中的数据
	r.ParseMultipartForm(1024)

	//初始化http处理结构体，把http信息压入结构体内
	var HttpInfo = library.HttpInfo{}
	HttpInfo.IsCli = false //非cli模式
	HttpInfo.ResponseWriter = w
	HttpInfo.Request = r
	HttpInfo.Form = r.Form
	HttpInfo.MultipartForm = r.MultipartForm

	//预制body内容raw访问使用
	var buf = new(bytes.Buffer)
	from, err := buf.ReadFrom(r.Body)
	if err != nil {
		library.SetLog(from, "错误输出")
		library.SetLog(err, "错误输出")
		library.OutJson(HttpInfo, map[string]interface{}{"code": "0", "msg": "预制body失败"})
		return
	}
	HttpInfo.Body = buf.String()

	//获得访问路径并去掉get参数
	var route = HttpInfo.GetReUrl()

	//获取当前url的路由设置 map[ac:order_list ct:CtlOrder method:GET route:/order/order_list]
	lr := library.Request{}
	Mount, RInfo := lr.GetRInfo(r, routeList, route)
	HttpInfo.Mount = Mount

	//判断是否存在路由
	if RInfo != nil {
		//循环控制器列表
		for k1, v1 := range RegisterMessage {
			//fmt.Print(k1, RInfo["ct"])
			//fmt.Print("\n")
			if k1 == RInfo["ct"] { //找到控制器
				//预创建控制器对象
				var methodArgs []reflect.Value
				methodArgs = append(methodArgs, reflect.ValueOf(HttpInfo))

				//把包含http内容的结构体推给控制器
				var CtlBox = reflect.ValueOf(v1).MethodByName(RInfo["ac"])
				CtlBox.Call(methodArgs)

				//完事了就直接退出
				return
			}
		}

	} else {
		library.OutJson(HttpInfo, map[string]interface{}{"code": "0", "msg": "路由不存在"})
	}
}

// LoadRouteCli cli模式加载控制器函数
func LoadRouteCli(ct string, ac string, fBox map[string]string) {
	//初始化http处理结构体，把http信息压入结构体内
	var HttpInfo = library.HttpInfo{}
	HttpInfo.IsCli = true //cli模式
	HttpInfo.Mount = fBox //由于是cli，就直接赋值假数据就完事了

	for k1, v1 := range RegisterMessage {
		if k1 == ct { //找到控制器
			//预创建控制器对象
			var methodArgs []reflect.Value
			methodArgs = append(methodArgs, reflect.ValueOf(HttpInfo))

			//把包含http内容的结构体推给控制器
			var CtlBox = reflect.ValueOf(v1).MethodByName(ac)
			CtlBox.Call(methodArgs)

			//完事了就直接退出
			return
		}
	}
}

func Init() {
	//加载预设服务模块
	var SS = library.ServerS{}
	SS.InitServerS()

	//初始化控制器池
	var ctl = controller.CtlIndex{}
	RegisterMessage = ctl.Init(SS)
}

// InitHttp 以http模式下初始化启动框架
func InitHttp() {

	Init() //集体初始化内容

	//获取启动配置
	deploy := config.Deploy{}
	con := deploy.Run()

	http.HandleFunc("/", LoadRouteHttp)
	err := http.ListenAndServe(con["LISTEN_ADDRESS"]+":"+con["PORT"], nil)
	if err != nil {
		library.SetLog(err, "错误输出")
		return
	}
}

// InitCli 以cli模式下的初始化框架
func InitCli(ct string, ac string, from string) {

	//定义一下先
	fBox := make(map[string]string)

	//分割字符
	f1 := strings.Split(from, "&")
	for _, v1 := range f1 {
		//排除掉那些没有等于号做分割符号的
		isGetData := strings.ContainsRune(v1, '=')
		if isGetData {
			f2 := strings.Split(v1, "=")
			//分割开来再解密，免得解密后的复杂字符影响
			f3, _ := url.QueryUnescape(f2[0])
			f4, _ := url.QueryUnescape(f2[1])
			fBox[f3] = f4 //解密后再赋值，就很稳了
		}
	}

	Init()                     //集体初始化内容
	LoadRouteCli(ct, ac, fBox) //拉起控制器
}
