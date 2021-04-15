package beego

var ValidationMessageTmpls = map[string]string{
	"Required":     "不能为空",
	"Min":          "最小值: %d",
	"Max":          "最大值: %d",
	"Range":        "数值的范围 %d - %d",
	"MinSize":      "最小长度: %d",
	"MaxSize":      "最大长度: %d",
	"Length":       "指定长度: %d",
	"Alpha":        "必须为有效的字母字符",
	"Numeric":      "必须为有效的数字",
	"AlphaNumeric": "必须为有效的字母或数字",
	"Match":        "必须匹配 %s",
	"NoMatch":      "必须不匹配 %s",
	"AlphaDash":    "必须为字母、数字或横杠(-_)",
	"Email":        "必须为邮箱格式",
	"IP":           "必须为有效的IP格式",
	"Base64":       "必须为有效的base64字符",
	"Mobile":       "必须有效的手机号码",
	"Tel":          "必须是有效的电话号码",
	"Phone":        "必须是有效的电话或手机号码",
	"ZipCode":      "必须是有效的邮政编码",
}
