package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Book struct {
	ID             uint64  `gorm:"primaryKey;autoIncrement"`
	UserID         uint64  `gorm:"index;not null"`
	ISBN           string  `gorm:"size:20;index"`
	Title          string  `gorm:"size:255;not null"`
	Author         string  `gorm:"size:255"`
	Publisher      string  `gorm:"size:100"`
	CoverImage     string  `gorm:"size:255"`
	Description    string  `gorm:"type:text"`
	Category       string  `gorm:"size:50"`
	Mode           int8    `gorm:"not null"`
	Status         int8    `gorm:"default:1;index"`
	DailyRent      *float64 `gorm:"type:decimal(10,2)"`
	WeeklyRent     *float64 `gorm:"type:decimal(10,2)"`
	Deposit        *float64 `gorm:"type:decimal(10,2)"`
	MinRentDays    int     `gorm:"default:1"`
	SellPrice      *float64 `gorm:"type:decimal(10,2)"`
	Images         string  `gorm:"type:json;not null"`
	PickupLocation string  `gorm:"size:255"`
	ViewCount      int     `gorm:"default:0"`
}

func (Book) TableName() string {
	return "books"
}

func floatPtr(v float64) *float64 {
	return &v
}

func main() {
	dsn := "root:24235096cyf@tcp(129.146.105.41:3306)/bookshare?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	var count int64
	db.Model(&Book{}).Count(&count)
	fmt.Printf("Current book count: %d\n", count)

	if count > 0 {
		fmt.Println("Deleting existing books...")
		db.Exec("DELETE FROM books")
	}

	type BookData struct {
		ISBN, Title, Author, Publisher, Description, Category string
		Mode int8
		DailyRent, WeeklyRent, Deposit, SellPrice float64
		MinRentDays int
		CoverImage, Images, PickupLocation string
		ViewCount int
	}

	books := []BookData{
		{"978-7530211538", "活着", "余华", "北京十月文艺出版社", "讲述农民福贵悲惨的一生，是中国当代文学的经典之作，销量超过2000万册。", "文学", 2, 3.00, 15.00, 5.00, 35.00, 1, "https://picsum.photos/seed/huozhe123/400/600", `["https://picsum.photos/seed/huozhe123/400/600"]`, "北京大学东门", 45890},
		{"978-7229100640", "三体", "刘慈欣", "重庆出版社", "中国科幻巅峰之作，亚洲首部雨果奖长篇小说，销量超过1900万册。", "科技", 2, 5.00, 35.00, 7.00, 68.00, 1, "https://picsum.photos/seed/santi456/400/600", `["https://picsum.photos/seed/santi456/400/600"]`, "清华大学北门", 52340},
		{"978-7530211533", "平凡的世界", "路遥", "北京十月文艺出版社", "激励无数青年的现实主义巨著，描写普通人在大时代中的奋斗历程。", "文学", 2, 4.00, 25.00, 5.00, 45.00, 1, "https://picsum.photos/seed/pingfan789/400/600", `["https://picsum.photos/seed/pingfan789/400/600"]`, "复旦大学邯郸校区", 38920},
		{"978-7530211534", "额尔古纳河右岸", "迟子建", "北京十月文艺出版社", "描写鄂温克族最后一个酋长的女人一生的史诗般的作品。", "文学", 2, 4.00, 22.00, 5.00, 42.00, 1, "https://picsum.photos/seed/eerguena/400/600", `["https://picsum.photos/seed/eerguena/400/600"]`, "上海交通大学闵行校区", 25670},
		{"978-7544263926", "我的阿勒泰", "李娟", "云南人民出版社", "2024年京东图书畅销榜冠军，描写新疆阿勒泰地区生活的散文集。", "文学", 2, 3.50, 20.00, 5.00, 38.00, 1, "https://picsum.photos/seed/alatai123/400/600", `["https://picsum.photos/seed/alatai123/400/600"]`, "浙江大学紫金港校区", 41230},
		{"978-7530211535", "我与地坛", "史铁生", "人民文学出版社", "史铁生的代表作，关于生命、残疾和母爱的深刻思考。", "文学", 2, 3.00, 18.00, 5.00, 32.00, 1, "https://picsum.photos/seed/ditan321/400/600", `["https://picsum.photos/seed/ditan321/400/600"]`, "南京大学仙林校区", 29870},
		{"978-7530211536", "一句顶一万句", "刘震云", "长江文艺出版社", "获得茅盾文学奖的作品，描写底层人物的孤独与挣扎。", "文学", 2, 4.00, 25.00, 5.00, 45.00, 1, "https://picsum.photos/seed/yijudi/400/600", `["https://picsum.photos/seed/yijudi/400/600"]`, "中山大学南校区", 18760},
		{"978-7572506927", "长安的荔枝", "马伯庸", "湖南文艺出版社", "历史悬疑题材，以紧凑剧情吸引读者，大唐盛世下小人物的传奇。", "历史", 2, 4.00, 22.00, 5.00, 40.00, 1, "https://picsum.photos/seed/lizhi123/400/600", `["https://picsum.photos/seed/lizhi123/400/600"]`, "中国科学技术大学", 32140},
		{"978-7020047549", "红高粱家族", "莫言", "人民文学出版社", "诺贝尔文学奖获奖作品，展现山东高密乡的传奇故事。", "文学", 2, 4.00, 25.00, 5.00, 42.00, 1, "https://picsum.photos/seed/honggaoliang/400/600", `["https://picsum.photos/seed/honggaoliang/400/600"]`, "武汉大学", 21560},
		{"978-7020047543", "四世同堂", "老舍", "人民文学出版社", "老舍先生的百万字巨著，描写抗战时期北平胡同里的人生百态。", "文学", 2, 5.00, 30.00, 7.00, 55.00, 1, "https://picsum.photos/seed/sishitong/400/600", `["https://picsum.photos/seed/sishitong/400/600"]`, "西安交通大学", 17680},
		{"978-7020047544", "骆驼祥子", "老舍", "人民文学出版社", "中国现代文学的经典之作，描写人力车夫祥子的悲剧命运。", "文学", 2, 3.00, 15.00, 5.00, 28.00, 1, "https://picsum.photos/seed/luotuoxiang/400/600", `["https://picsum.photos/seed/luotuoxiang/400/600"]`, "同济大学四平路校区", 19870},
		{"978-7530211537", "白鹿原", "陈忠实", "人民文学出版社", "茅盾文学奖获奖作品，渭河平原五十年的沧桑变迁。", "文学", 2, 4.00, 25.00, 5.00, 48.00, 1, "https://picsum.photos/seed/bailuyuan/400/600", `["https://picsum.photos/seed/bailuyuan/400/600"]`, "哈尔滨工业大学", 22340},
		{"978-7530211538", "围城", "钱钟书", "人民文学出版社", "中国现代文学的讽刺喜剧巅峰之作。", "文学", 1, 3.00, 18.00, 5.00, 32.00, 1, "https://picsum.photos/seed/weicheng123/400/600", `["https://picsum.photos/seed/weicheng123/400/600"]`, "北京外国语大学", 28760},
		{"978-7530211539", "蛙", "莫言", "上海文艺出版社", "诺贝尔文学奖获奖作品，描写中国计划生育政策的复杂影响。", "文学", 2, 4.00, 22.00, 5.00, 40.00, 1, "https://picsum.photos/seed/wa123456/400/600", `["https://picsum.photos/seed/wa123456/400/600"]`, "厦门大学", 15430},
		{"978-7530211540", "丰乳肥臀", "莫言", "作家出版社", "莫言的代表作之一，通过上官鲁氏一家的命运展现中国百年历史。", "文学", 2, 4.00, 25.00, 5.00, 45.00, 1, "https://picsum.photos/seed/fengru123/400/600", `["https://picsum.photos/seed/fengru123/400/600"]`, "华南理工大学", 12340},
		{"978-7533944319", "文化苦旅", "余秋雨", "浙江古籍出版社", "余秋雨的经典散文集，将文化历史与山水游记融为一体。", "历史", 1, 3.50, 20.00, 5.00, 38.00, 1, "https://picsum.photos/seed/wenhuaku/400/600", `["https://picsum.photos/seed/wenhuaku/400/600"]`, "四川大学望江校区", 19870},
		{"978-7533944318", "千年一叹", "余秋雨", "作家出版社", "余秋雨考察人类重要文明遗址的随笔集。", "历史", 1, 3.50, 20.00, 5.00, 36.00, 1, "https://picsum.photos/seed/qianniantan/400/600", `["https://picsum.photos/seed/qianniantan/400/600"]`, "天津大学", 11230},
		{"978-7020047540", "红楼梦", "曹雪芹", "人民文学出版社", "中国古典四大名著之首，被誉为中国封建社会的百科全书。", "文学", 1, 4.00, 22.00, 5.00, 48.00, 1, "https://picsum.photos/seed/hongloumeng/400/600", `["https://picsum.photos/seed/hongloumeng/400/600"]`, "北京师范大学", 35670},
		{"978-7020047541", "西游记", "吴承恩", "人民文学出版社", "中国古典四大名著之一，孙悟空保唐僧西天取经的神话故事。", "文学", 1, 3.50, 18.00, 5.00, 42.00, 1, "https://picsum.photos/seed/xiyouji123/400/600", `["https://picsum.photos/seed/xiyouji123/400/600"]`, "华东师范大学", 32140},
		{"978-7020047542", "水浒传", "施耐庵", "人民文学出版社", "中国古典四大名著之一，描写梁山好汉替天行道的故事。", "文学", 1, 3.50, 18.00, 5.00, 40.00, 1, "https://picsum.photos/seed/shuihuzhuan/400/600", `["https://picsum.photos/seed/shuihuzhuan/400/600"]`, "南开大学", 28760},
		{"978-7020047545", "三国演义", "罗贯中", "人民文学出版社", "中国古典四大名著之一，魏蜀吴三国争霸的历史演义。", "历史", 1, 3.50, 18.00, 5.00, 42.00, 1, "https://picsum.photos/seed/sanguoyanyi/400/600", `["https://picsum.photos/seed/sanguoyanyi/400/600"]`, "东南大学九龙湖校区", 31250},
		{"978-7510855353", "繁花", "金宇澄", "上海文艺出版社", "茅盾文学奖获奖作品，用上海方言书写上海的市井生活。", "文学", 2, 4.00, 25.00, 5.00, 48.00, 1, "https://picsum.photos/seed/fanhua123/400/600", `["https://picsum.photos/seed/fanhua123/400/600"]`, "华中科技大学", 16540},
		{"978-7540487647", "云边有个小卖部", "张嘉佳", "湖南文艺出版社", "写给每个人心中的山和海，感人至深的治愈系小说。", "文学", 2, 3.50, 20.00, 5.00, 38.00, 1, "https://picsum.photos/seed/yunbianxmd/400/600", `["https://picsum.photos/seed/yunbianxmd/400/600"]`, "北京航空航天大学", 27890},
		{"978-7540487648", "天堂旅行团", "张嘉佳", "湖南文艺出版社", "讲述一个绝望的人寻找活下去的理由的感人故事。", "文学", 2, 3.50, 20.00, 5.00, 36.00, 1, "https://picsum.photos/seed/tiantanglxt/400/600", `["https://picsum.photos/seed/tiantanglxt/400/600"]`, "南京航空航天大学", 19870},
		{"978-7544273956", "解忧杂货店", "东野圭吾", "南海出版公司", "日本作家东野圭吾的治愈系代表作，跨越时空的温暖故事。", "文学", 2, 3.00, 18.00, 5.00, 35.00, 1, "https://picsum.photos/seed/jieyouzhd/400/600", `["https://picsum.photos/seed/jieyouzhd/400/600"]`, "重庆大学", 25670},
		{"978-7229100641", "三体2：黑暗森林", "刘慈欣", "重庆出版社", "三体系列的第二部，宇宙文明之间的博弈与生存法则。", "科技", 2, 5.00, 35.00, 7.00, 68.00, 1, "https://picsum.photos/seed/santi2heian/400/600", `["https://picsum.photos/seed/santi2heian/400/600"]`, "电子科技大学清水河校区", 31240},
		{"978-7229100642", "三体3：死神永生", "刘慈欣", "重庆出版社", "三体系列的最终章，跨越亿万年的宇宙史诗。", "科技", 2, 5.00, 35.00, 7.00, 72.00, 1, "https://picsum.photos/seed/santi3siming/400/600", `["https://picsum.photos/seed/santi3siming/400/600"]`, "西北工业大学长安校区", 28760},
		{"978-7542667894", "金庸全集", "金庸", "广州出版社", "武侠小说大师金庸的经典作品合集。", "文学", 2, 6.00, 40.00, 7.00, 680.00, 1, "https://picsum.photos/seed/jinyong123/400/600", `["https://picsum.photos/seed/jinyong123/400/600"]`, "大连理工大学", 43210},
		{"978-7021161960", "暂坐", "贾平凹", "作家出版社", "贾平凹的最新长篇小说，描写西安城中一群独立女性的都市生活。", "文学", 2, 4.00, 25.00, 5.00, 48.00, 1, "https://picsum.photos/seed/zanzuo123/400/600", `["https://picsum.photos/seed/zanzuo123/400/600"]`, "北京理工大学", 12350},
		{"978-7542667893", "蛙", "莫言", "上海文艺出版社", "诺贝尔文学奖获奖作品，关于中国计划生育政策的深刻文学作品。", "文学", 2, 4.00, 22.00, 5.00, 42.00, 1, "https://picsum.photos/seed/wa789abc/400/600", `["https://picsum.photos/seed/wa789abc/400/600"]`, "山东大学中心校区", 18760},
	}

	fmt.Println("Inserting 30 books with proper cover images...")

	for i, b := range books {
		book := Book{
			UserID:         1,
			ISBN:           b.ISBN,
			Title:          b.Title,
			Author:         b.Author,
			Publisher:      b.Publisher,
			CoverImage:     b.CoverImage,
			Description:    b.Description,
			Category:       b.Category,
			Mode:           b.Mode,
			Status:         1,
			DailyRent:      floatPtr(b.DailyRent),
			WeeklyRent:     floatPtr(b.WeeklyRent),
			Deposit:        floatPtr(b.Deposit),
			MinRentDays:    b.MinRentDays,
			SellPrice:      floatPtr(b.SellPrice),
			Images:         b.Images,
			PickupLocation: b.PickupLocation,
			ViewCount:      b.ViewCount,
		}

		if err := db.Create(&book).Error; err != nil {
			log.Printf("Failed to insert book %d (%s): %v", i+1, b.Title, err)
		} else {
			fmt.Printf("Inserted: %s by %s (cover: %s)\n", b.Title, b.Author, b.CoverImage)
		}
	}

	var newCount int64
	db.Model(&Book{}).Count(&newCount)
	fmt.Printf("\nTotal books inserted: %d\n", newCount)
}
