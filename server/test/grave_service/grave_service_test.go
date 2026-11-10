package grave_service

import (
	"testing"
	"time"

	"server/database/connection"
	"server/errors"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	userRepo "server/repository/user"
	graveservice "server/service/grave"
	"server/test/testutil"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	reGraveByID    = `SELECT .+ FROM \x60graves\x60 WHERE \x60graves\x60.\x60id\x60 = \? .+`
	reMemberCount  = `SELECT count\(\*\) FROM \x60grave_members\x60 WHERE grave_id = \? AND user_id = \?`
	reUserByID     = `SELECT .+ FROM \x60users\x60 WHERE \x60users\x60.\x60id\x60 = \? .+`
	reInsGrave     = `INSERT INTO \x60graves\x60`
	reInsMember    = `INSERT INTO \x60grave_members\x60`
	reUpdGrave     = `UPDATE \x60graves\x60 SET .+ WHERE .+`
	reDelRelations = `DELETE FROM .+ WHERE grave_id = \?`
	reDelGrave     = `UPDATE \x60graves\x60 SET .+ WHERE .+`
)

func graveRow(id uint, name string, createdBy uint) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "grave_type", "grave_mode", "layout_style",
		"description", "cover_image", "status", "latitude", "longitude",
		"address", "created_by", "created_at", "updated_at"}).
		AddRow(id, name, 1, 1, 1, "", "", 1, nil, nil, "", createdBy, time.Now(), time.Now())
}

func userRow(id uint, username, nickname string) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "nickname", "avatar", "email", "password", "created_at", "updated_at"}).
		AddRow(id, username, nickname, "", "", "$2a$10$test", time.Now(), time.Now())
}

func TestGetByID_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)
	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectQuery(reGraveByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(graveRow(1, "测试墓园", 2))
	mockDB.Mock.ExpectQuery(reMemberCount).WithArgs(uint(1), uint(1)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(2), sqlmock.AnyArg()).WillReturnRows(userRow(2, "creator", "创建者"))

	resp, err := svc.GetByID(1, 1)
	if err != nil {
		t.Fatalf("GetByID 失败: %v", err)
	}
	if resp.Name != "测试墓园" {
		t.Errorf("Name = %s, want 测试墓园", resp.Name)
	}
}

func TestGetByID_AccessDenied(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)
	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectQuery(reGraveByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(graveRow(1, "测试墓园", 2))
	mockDB.Mock.ExpectQuery(reMemberCount).WithArgs(uint(1), uint(3)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	_, err := svc.GetByID(3, 1)
	if err != errors.ErrGraveAccessDenied {
		t.Errorf("期望 ErrGraveAccessDenied, got %v", err)
	}
}

func TestListMy_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)
	svc := graveservice.NewGraveService(gr, ur, mr)

	rows := sqlmock.NewRows([]string{"id", "name", "grave_type", "grave_mode", "layout_style",
		"description", "cover_image", "status", "latitude", "longitude",
		"address", "created_by", "created_at", "updated_at"}).
		AddRow(1, "墓园A", 1, 1, 1, "", "", 1, nil, nil, "", 1, time.Now(), time.Now()).
		AddRow(2, "墓园B", 1, 1, 1, "", "", 1, nil, nil, "", 2, time.Now(), time.Now())
	mockDB.Mock.ExpectQuery("SELECT count\\(\\*\\) FROM `graves` WHERE .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mockDB.Mock.ExpectQuery("SELECT .+ FROM `graves` WHERE .+").WillReturnRows(rows)
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(userRow(1, "u1", "n1"))
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(2), sqlmock.AnyArg()).WillReturnRows(userRow(2, "u2", "n2"))

	resp, _, err := svc.ListMy(1, 1, 10)
	if err != nil {
		t.Fatalf("ListMy 失败: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Errorf("len(Items) = %d, want 2", len(resp.Items))
	}
}

func TestUpdate_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)
	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectQuery(reGraveByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(graveRow(1, "旧名称", 1))
	mockDB.Mock.ExpectQuery(reMemberCount).WithArgs(uint(1), uint(1)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mockDB.Mock.ExpectExec(reUpdGrave).WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(userRow(1, "creator", "创建者"))

	newName := "新名称"
	resp, err := svc.Update(1, 1, &graveservice.UpdateGraveParams{Name: &newName})
	if err != nil {
		t.Fatalf("Update 失败: %v", err)
	}
	if resp.Name != "新名称" {
		t.Errorf("Name = %s, want 新名称", resp.Name)
	}
}

func TestUpdate_AccessDenied(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)
	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectQuery(reGraveByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(graveRow(1, "测试", 2))
	mockDB.Mock.ExpectQuery(reMemberCount).WithArgs(uint(1), uint(3)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	newName := "新名称"
	_, err := svc.Update(3, 1, &graveservice.UpdateGraveParams{Name: &newName})
	if err != errors.ErrGraveAccessDenied {
		t.Errorf("期望 ErrGraveAccessDenied, got %v", err)
	}
}

func TestCreate_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)

	oldDB := connection.DB
	connection.DB = mockDB.DB
	defer func() { connection.DB = oldDB }()

	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectBegin()
	mockDB.Mock.ExpectExec(reInsGrave).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.Mock.ExpectExec(reInsMember).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.Mock.ExpectQuery(reUserByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(userRow(1, "creator", "创建者"))
	mockDB.Mock.ExpectCommit()

	resp, err := svc.Create(1, "新墓园", 1, 1, 1, "描述", "", nil, nil, "")
	if err != nil {
		t.Fatalf("Create 失败: %v", err)
	}
	if resp.Name != "新墓园" {
		t.Errorf("Name = %s, want 新墓园", resp.Name)
	}
}

func TestDelete_Success(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)

	oldDB := connection.DB
	connection.DB = mockDB.DB
	defer func() { connection.DB = oldDB }()

	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectQuery(reGraveByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(graveRow(1, "待删除", 1))
	mockDB.Mock.ExpectQuery(reMemberCount).WithArgs(uint(1), uint(1)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mockDB.Mock.ExpectBegin()
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelRelations).WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.Mock.ExpectExec(reDelGrave).WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.Mock.ExpectCommit()

	if err := svc.Delete(1, 1); err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}
}

func TestDelete_AccessDenied(t *testing.T) {
	mockDB := testutil.NewMockDB(t)
	defer mockDB.Close()

	gr := graverepo.NewGraveRepository(mockDB.DB)
	ur := userRepo.NewUserRepository(mockDB.DB)
	mr := memberrepo.NewGraveMemberRepository(mockDB.DB)
	svc := graveservice.NewGraveService(gr, ur, mr)

	mockDB.Mock.ExpectQuery(reGraveByID).WithArgs(uint(1), sqlmock.AnyArg()).WillReturnRows(graveRow(1, "测试", 2))
	mockDB.Mock.ExpectQuery(reMemberCount).WithArgs(uint(1), uint(3)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	if err := svc.Delete(3, 1); err != errors.ErrGraveAccessDenied {
		t.Errorf("期望 ErrGraveAccessDenied, got %v", err)
	}
}
