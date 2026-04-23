package handler

import (
	"bookshare/pkg/response"

	"github.com/gin-gonic/gin"
)

type GuideHandler struct{}

func NewGuideHandler() *GuideHandler {
	return &GuideHandler{}
}

func (h *GuideHandler) GetGuide(c *gin.Context) {
	guide := gin.H{
		"sections": []gin.H{
			{
				"id":          "publish",
				"title":       "如何发布图书",
				"icon":        "book",
				"description": "将闲置图书发布到平台，让知识流动起来",
				"steps": []gin.H{
					{"step": 1, "title": "登录账号", "content": "首先登录您的账号，如果没有账号请先注册"},
					{"step": 2, "title": "点击发布", "content": "在导航栏点击「发布」按钮，进入发布页面"},
					{"step": 3, "title": "填写信息", "content": "填写图书信息，可以扫描ISBN自动获取图书信息"},
					{"step": 4, "title": "选择模式", "content": "选择流转方式：租赁、出售或赠送"},
					{"step": 5, "title": "提交发布", "content": "确认信息无误后点击发布，图书将展示在平台"},
				},
			},
			{
				"id":          "rent",
				"title":       "如何租赁图书",
				"icon":        "rent",
				"description": "以较低的价格租借需要的图书",
				"steps": []gin.H{
					{"step": 1, "title": "浏览图书", "content": "在图书列表中筛选「租赁」模式的图书"},
					{"step": 2, "title": "查看详情", "content": "点击感兴趣的图书查看详细信息"},
					{"step": 3, "title": "选择天数", "content": "选择您需要租赁的天数"},
					{"step": 4, "title": "支付押金", "content": "支付租金和押金后等待对方确认"},
					{"step": 5, "title": "取书阅读", "content": "对方确认后，按约定地点取书"},
					{"step": 6, "title": "归还图书", "content": "阅读完毕后按时归还，押金将退还给您"},
				},
			},
			{
				"id":          "buy",
				"title":       "如何购买图书",
				"icon":        "buy",
				"description": "直接购买心仪的二手图书",
				"steps": []gin.H{
					{"step": 1, "title": "浏览图书", "content": "在图书列表中筛选「出售」模式的图书"},
					{"step": 2, "title": "查看详情", "content": "点击感兴趣的图书查看详细信息和实物图片"},
					{"step": 3, "title": "立即购买", "content": "点击「立即购买」按钮"},
					{"step": 4, "title": "支付款项", "content": "支付购买款项，平台将代为保管"},
					{"step": 5, "title": "等待交付", "content": "卖家确认后按约定方式交付图书"},
					{"step": 6, "title": "确认收货", "content": "收到图书后确认收货，款项将打给卖家"},
				},
			},
			{
				"id":          "gift",
				"title":       "如何领取赠送图书",
				"icon":        "gift",
				"description": "免费领取他人赠送的图书",
				"steps": []gin.H{
					{"step": 1, "title": "浏览图书", "content": "在图书列表中筛选「赠送」模式的图书"},
					{"step": 2, "title": "查看详情", "content": "点击感兴趣的图书查看详细信息"},
					{"step": 3, "title": "申请领取", "content": "点击「申请领取」按钮，可以附上留言"},
					{"step": 4, "title": "等待确认", "content": "赠送人将审核申请并选择领取人"},
					{"step": 5, "title": "按约取书", "content": "被选中后按约定地点和时间取书"},
				},
			},
			{
				"id":          "tips",
				"title":       "使用小贴士",
				"icon":        "tips",
				"description": "让您的使用体验更好",
				"items": []gin.H{
					{"title": "完善个人信息", "content": "填写学校和校区信息，可以获得更精准的匹配"},
					{"title": "上传清晰图片", "content": "发布图书时上传清晰的实物图片，更容易被选中"},
					{"title": "及时沟通", "content": "通过平台与对方保持良好沟通，交易更顺畅"},
					{"title": "按时归还", "content": "租赁图书请按时归还，维护良好信用"},
					{"title": "诚实描述", "content": "发布图书时如实描述图书状况，建立信任"},
				},
			},
		},
		"faq": []gin.H{
			{
				"question": "押金什么时候退还？",
				"answer":   "租赁图书归还并验收通过后，押金将自动退还到您的账户余额。",
			},
			{
				"question": "如何提高信用分？",
				"answer":   "按时归还租赁图书、及时确认收货、获得好评等都可以提高信用分。",
			},
			{
				"question": "遇到纠纷怎么办？",
				"answer":   "如遇到交易纠纷，可以通过订单详情页联系客服介入处理。",
			},
			{
				"question": "可以取消订单吗？",
				"answer":   "在对方确认前可以取消订单。确认后取消需要双方协商。",
			},
		},
	}

	response.Success(c, guide)
}

func (h *GuideHandler) SearchBooks(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.ParamError(c, nil)
		return
	}

	response.Success(c, gin.H{
		"keyword": keyword,
		"hint":    "请使用 /api/v1/books 接口，传入 keyword 参数进行搜索",
	})
}
