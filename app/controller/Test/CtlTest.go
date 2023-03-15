package controller

import (
	"core/library"
	"model"
	"reflect"
)

type CtlTest struct {
	DB  library.MysqlG
	out map[string]string
}

func (r *CtlTest) Test(CH library.HttpInfo) {

	postData := map[string]string{
		"a": CH.R("a", "a"),
		"b": CH.R("b", "b"),
		"c": CH.R("c", "c"),
		"d": CH.R("d", "d"),
		"e": CH.R("e", "e"),
		"f": CH.R("f", "f"),
	}

	library.SetLog(postData, "输出打印")

	//------------- 以下是抽象化调用 --------------//

	//初始化MOD池
	var (
		Mod             = model.ModIndex{}
		RegisterMessage = Mod.Init(r.DB)
	)

	var (
		ModName    = "ModTest"
		ModFunName = "Test"
		methodArgs []reflect.Value
	)

	//压入数据
	methodArgs = append(methodArgs, reflect.ValueOf(postData))

	//反射调用
	var ModBox = reflect.ValueOf(RegisterMessage[ModName]).MethodByName(ModFunName)
	var aaa = ModBox.Call(methodArgs)[0]
	library.SetLog(aaa, "反射调用输出打印")

	//------------- 以上是抽象化调用 --------------//

	//------------- 以下的正常拉起 --------------//

	//Mod2 := ModTest.ModTest{DB: r.DB}
	//var aab = Mod2.Test(postData) //直接调用
	//library.SetLog(aab, "直接调用输出打印")

	//------------- 以上的正常拉起 --------------//

	//library.OutJson(CH.ResponseWriter, aab) //输出到web页面
}

func (r *CtlTest) TestA(CH library.HttpInfo) {

	postData := map[string]string{
		"a":  CH.R("a", "a"),
		"b":  CH.R("b", "b"),
		"c":  CH.R("c", "c"),
		"d":  CH.R("d", "d"),
		"e":  CH.R("e", "e"),
		"f":  CH.R("f", "f"),
		"ct": CH.R("ct", "ctc"),
		"ac": CH.R("ac", "aca"),
	}

	//library.OutJson(CH.ResponseWriter, postData) //输出到web页面

	library.OutHtml(CH.ResponseWriter, "test.html", postData)
}