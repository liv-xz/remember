package init

import (
	"server/database/connection"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
	gravemembermodel "server/model/grave_member"
	gslmodel "server/model/grave_share_link"
	kdr "server/model/kinship_deceased_relation"
	kr "server/model/kinship_relation"
	loginlogmodel "server/model/login_log"
	memorialeventmodel "server/model/memorial_event"
	memorialrecordmodel "server/model/memorial_record"
	mr "server/model/memorial_relation"
	oplogmodel "server/model/operation_log"
	rtmodel "server/model/token"
	model "server/model/user"

	"gorm.io/gorm"
)

func CreateTable() error {
	db := connection.DB
	if db == nil {
		return gorm.ErrInvalidDB
	}

	if db.Migrator().HasColumn(&model.User{}, "open_id") {
		db.Exec("UPDATE users SET open_id = NULL WHERE open_id = ''")
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")
	err := db.AutoMigrate(
		&model.User{},
		&model.UserConfig{},
		&gravemodel.Grave{},
		&deceasedmodel.Deceased{},
		&kdr.DeceasedRelation{},
		&gravemembermodel.GraveMember{},
		&memorialrecordmodel.MemorialRecord{},
		&memorialeventmodel.MemorialEvent{},
		&rtmodel.RefreshToken{},
		&loginlogmodel.LoginLog{},
		&oplogmodel.OperationLog{},
		&gslmodel.GraveShareLink{},
		&kr.KinshipRelation{},
		&mr.MemorialRelation{},
	)
	if err != nil {
		return err
	}

	return nil
}
