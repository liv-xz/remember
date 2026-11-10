package main

import (
	"fmt"

	"server/config"
	"server/database/connection"
	initdatabase "server/database/init"
	"server/flags"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
	membermodel "server/model/grave_member"
	gslmodel "server/model/grave_share_link"
	kdrmodel "server/model/kinship_deceased_relation"
	krmodel "server/model/kinship_relation"
	eventmodel "server/model/memorial_event"
	recordmodel "server/model/memorial_record"
	usermodel "server/model/user"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("=== RememberMe 演示数据导入 ===")

	flags.Parse1()
	config.Init(flags.FlagOptions.File)
	if err := connection.Init(); err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}
	if err := initdatabase.CreateTable(); err != nil {
		panic(fmt.Sprintf("创建表失败: %v", err))
	}

	db := connection.DB

	lat := 31.2304
	lng := 121.4737

	// 1. 创建用户
	admin := &usermodel.User{Username: "admin", Nickname: "系统管理员", Email: "admin@example.com"}
	pwd1, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin.Password = string(pwd1)
	db.Create(admin)

	zhangsan := &usermodel.User{Username: "zhangsan", Nickname: "张小明", Email: "zhangsan@example.com"}
	pwd2, _ := bcrypt.GenerateFromPassword([]byte("zs123"), bcrypt.DefaultCost)
	zhangsan.Password = string(pwd2)
	db.Create(zhangsan)

	db.Create(&usermodel.UserConfig{UserID: admin.ID})
	db.Create(&usermodel.UserConfig{UserID: zhangsan.ID})
	fmt.Printf("  ✓ 用户: admin(uid=%d) / admin123, zhangsan(uid=%d) / zs123\n", admin.ID, zhangsan.ID)

	// 2. 创建墓园
	grave := &gravemodel.Grave{
		Name: "张氏家族纪念馆", GraveType: 3, GraveMode: 2, LayoutStyle: 1,
		Description: "张氏家族的数字纪念馆，记录家族先人的生平与故事。这里承载着我们对先辈的思念与敬意。",
		Status:      1, Latitude: &lat, Longitude: &lng, Address: "上海市静安区",
		CreatedBy: admin.ID,
	}
	db.Create(grave)
	fmt.Printf("  ✓ 墓园: 张氏家族纪念馆 (gid=%d)\n", grave.ID)

	// 3. 创建成员
	db.Create(&membermodel.GraveMember{GraveID: grave.ID, UserID: admin.ID, Role: 10, RelationDesc: "墓主"})
	db.Create(&membermodel.GraveMember{GraveID: grave.ID, UserID: zhangsan.ID, Role: 20, RelationDesc: "家族管理员"})
	fmt.Println("  ✓ 成员: admin(墓主), 张小明(管理员)")

	// 4. 创建4位逝者
	d1 := &deceasedmodel.Deceased{
		GraveID: grave.ID, Name: "张德厚", Gender: 1, BirthDate: "1920-03-15", DeathDate: "1995-12-08",
		Biography: "张德厚，1920年生于江苏南京。青年投身革命，解放后在上海纺织厂工作至退休。一生勤劳朴实，抚养三子一女成人。爱好书法，尤擅颜体。享年75岁。",
		Epitaph:   "德厚流光，后世永怀", BurialOrder: 1, SortOrder: 1,
	}
	db.Create(d1)

	d2 := &deceasedmodel.Deceased{
		GraveID: grave.ID, Name: "李秀兰", Gender: 2, BirthDate: "1923-07-20", DeathDate: "2000-03-05",
		Biography: "李秀兰，1923年生于浙江宁波。相夫教子，勤俭持家。做得一手好菜，红烧肉是家族永远的记忆。晚年信佛，每日诵经。享年77岁。",
		Epitaph:   "慈母手中线，游子身上衣", BurialOrder: 2, SortOrder: 2,
	}
	db.Create(d2)

	d3 := &deceasedmodel.Deceased{
		GraveID: grave.ID, Name: "张建国", Gender: 1, BirthDate: "1950-10-01", DeathDate: "2020-02-14",
		Biography: "张建国，1950年生于上海。下乡知青，1978年返城考入同济大学土木工程系。参与多个上海地标建筑建设。为人正直，2020年因病去世，享年70岁。",
		Epitaph:   "为国为家，鞠躬尽瘁", BurialOrder: 3, SortOrder: 3,
	}
	db.Create(d3)

	d4 := &deceasedmodel.Deceased{
		GraveID: grave.ID, Name: "王美玲", Gender: 2, BirthDate: "1952-05-10",
		Biography: "王美玲，1952年生于上海。退休教师，教书育人三十年，桃李满天下。性格开朗，爱好园艺和烹饪。现居上海，与女儿一家同住。",
		Epitaph:   "春风化雨，师恩永存", BurialOrder: 4, SortOrder: 4,
	}
	db.Create(d4)
	fmt.Printf("  ✓ 逝者: %s(did=%d), %s(did=%d), %s(did=%d), %s(did=%d)\n",
		d1.Name, d1.ID, d2.Name, d2.ID, d3.Name, d3.ID, d4.Name, d4.ID)

	// 5. 创建亲缘关系 (用户→逝者)
	confirmedBy := admin.ID
	db.Create(&krmodel.KinshipRelation{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: d1.ID, RelationType: 9,
		ConfirmationStatus: 1, ConfirmedBy: &confirmedBy,
	})
	db.Create(&krmodel.KinshipRelation{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: d2.ID, RelationType: 10,
		ConfirmationStatus: 1, ConfirmedBy: &confirmedBy,
	})
	db.Create(&krmodel.KinshipRelation{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: d3.ID, RelationType: 7,
		ConfirmationStatus: 1, ConfirmedBy: &confirmedBy,
	})
	db.Create(&krmodel.KinshipRelation{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: d4.ID, RelationType: 8,
		ConfirmationStatus: 1, ConfirmedBy: &confirmedBy,
	})
	db.Create(&krmodel.KinshipRelation{
		GraveID: grave.ID, UserID: zhangsan.ID, DeceasedID: d1.ID, RelationType: 9,
		ConfirmationStatus: 1, ConfirmedBy: &confirmedBy,
	})
	db.Create(&krmodel.KinshipRelation{
		GraveID: grave.ID, UserID: zhangsan.ID, DeceasedID: d3.ID, RelationType: 7,
		ConfirmationStatus: 1, ConfirmedBy: &confirmedBy,
	})
	fmt.Println("  ✓ 亲缘关系: 6条 (祖父/祖母/父亲/母亲, 全部已确认)")

	// 6. 创建逝者间关系
	db.Create(&kdrmodel.DeceasedRelation{
		GraveID: grave.ID, DeceasedID1: d1.ID, DeceasedID2: d2.ID,
		RelationType: 1, Direction: 0,
	})
	db.Create(&kdrmodel.DeceasedRelation{
		GraveID: grave.ID, DeceasedID1: d3.ID, DeceasedID2: d4.ID,
		RelationType: 1, Direction: 0,
	})
	db.Create(&kdrmodel.DeceasedRelation{
		GraveID: grave.ID, DeceasedID1: d1.ID, DeceasedID2: d3.ID,
		RelationType: 2, Direction: 1,
	})
	fmt.Println("  ✓ 逝者间关系: 3条 (祖父母夫妻, 父母夫妻, 祖→父 亲子)")

	// 7. 创建纪念记录
	db.Create(&recordmodel.MemorialRecord{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: &d1.ID, Type: 1,
		Message: "爷爷，今年清明又来看您了。您种的桂花树今年开了特别多的花，满院子都是香的。小美今年考上大学了，您一定很高兴吧。我们都很好，您放心。",
	})
	db.Create(&recordmodel.MemorialRecord{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: &d3.ID, Type: 1,
		Message: "父亲，您离开我们六年了。家里一切都好，母亲身体硬朗，每天还去公园打太极拳。我今年也升职了，您一直说要脚踏实地，我会记住的。",
	})
	db.Create(&recordmodel.MemorialRecord{
		GraveID: grave.ID, UserID: admin.ID, DeceasedID: &d2.ID, Type: 1,
		Message: "奶奶，您做的红烧肉味道再也没有人能做出来了。我试着做了好多次，总觉得不对。可能是少了您那份心意吧。",
	})
	db.Create(&recordmodel.MemorialRecord{
		GraveID: grave.ID, UserID: zhangsan.ID, DeceasedID: &d4.ID, Type: 1,
		Message: "妈妈，清明安康。小美考上复旦了，和您当年一样。她说想当老师，和外婆一样教书育人。我们都为您骄傲。",
	})
	fmt.Println("  ✓ 纪念记录: 4条")

	// 8. 创建纪念事件
	db.Create(&eventmodel.MemorialEvent{
		GraveID: grave.ID, DeceasedID: d1.ID, EventType: 0,
		EventDate: "2026-04-05", Title: "清明节扫墓", Lunar: 0,
	})
	db.Create(&eventmodel.MemorialEvent{
		GraveID: grave.ID, DeceasedID: d3.ID, EventType: 0,
		EventDate: "2026-12-21", Title: "冬至家族聚会", Lunar: 0,
	})
	fmt.Println("  ✓ 纪念事件: 2条 (清明节扫墓, 冬至家族聚会)")

	// 9. 创建分享链接
	db.Create(&gslmodel.GraveShareLink{
		GraveID: grave.ID, CreatedBy: admin.ID, ShareCode: "memorial2026",
		MaxUseCount: 0, UsedCount: 0, IsActive: true,
	})
	fmt.Println("  ✓ 分享链接: memorial2026 (已激活, 不限次数)")

	fmt.Println()
	fmt.Println("=== ✅ 演示数据导入完成 ===")
	fmt.Println()
	fmt.Println("  用户端:  http://localhost:3000")
	fmt.Printf("  账号:    admin / admin123 (uid=%d)\n", admin.ID)
	fmt.Printf("          zhangsan / zs123 (uid=%d)\n", zhangsan.ID)
}
