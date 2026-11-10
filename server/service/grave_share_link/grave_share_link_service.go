package grave_share_link

import (
	"crypto/rand"
	"encoding/hex"
	"server/errors"
	membermodel "server/model/grave_member"
	gslmodel "server/model/grave_share_link"
	krmodel "server/model/kinship_relation"
	mrmodel "server/model/memorial_relation"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	gslrepo "server/repository/grave_share_link"
	krrepo "server/repository/kinship_relation"
	mrrepo "server/repository/memorial_relation"
	"server/types"
	"server/utils"
	"time"
)

type GraveShareLinkService struct {
	linkRepo     *gslrepo.GraveShareLinkRepository
	graveRepo    *graverepo.GraveRepository
	memberRepo   *memberrepo.GraveMemberRepository
	deceasedRepo *deceasedrepo.DeceasedRepository
	mrRepo       *mrrepo.MemorialRelationRepository
	krRepo       *krrepo.KinshipRelationRepository
}

func NewGraveShareLinkService(
	linkRepo *gslrepo.GraveShareLinkRepository,
	graveRepo *graverepo.GraveRepository,
	memberRepo *memberrepo.GraveMemberRepository,
	deceasedRepo *deceasedrepo.DeceasedRepository,
	mrRepo *mrrepo.MemorialRelationRepository,
	krRepo *krrepo.KinshipRelationRepository,
) *GraveShareLinkService {
	return &GraveShareLinkService{
		linkRepo:     linkRepo,
		graveRepo:    graveRepo,
		memberRepo:   memberRepo,
		deceasedRepo: deceasedRepo,
		mrRepo:       mrRepo,
		krRepo:       krRepo,
	}
}

// Create 创建分享链接
// 只有墓地创建者可以创建分享链接
// 分享链接用于邀请家人访问墓园并建立与逝者的关系
// 参数:
//   - graveID: 墓地ID
//   - userID: 用户ID（必须是墓地创建者）
//   - expireDays: 过期天数（0表示永久有效）
//
// 返回:
//   - GraveShareLinkResp: 创建的分享链接响应数据
//   - error: 错误信息
func (s *GraveShareLinkService) Create(graveID, userID uint, expireDays int) (*types.GraveShareLinkResp, error) {
	// 1. 验证墓地是否存在
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}

	// 2. 权限验证：只有墓地创建者可以创建分享链接
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	// 3. 生成唯一的分享码（32位十六进制字符串）
	shareCode, err := s.generateShareCode()
	if err != nil {
		return nil, err
	}

	// 4. 计算过期时间
	var expireAt *time.Time
	if expireDays > 0 {
		// 如果指定了过期天数，计算具体的过期时间
		t := time.Now().AddDate(0, 0, expireDays)
		expireAt = &t
	}

	// 5. 创建分享链接记录
	link := &gslmodel.GraveShareLink{
		GraveID:     graveID,
		CreatedBy:   userID,
		ShareCode:   shareCode,
		ExpireAt:    expireAt,
		MaxUseCount: 0, // 0表示无使用次数限制
		UsedCount:   0,
		IsActive:    true,
	}

	// 6. 保存到数据库
	if err := s.linkRepo.Create(link); err != nil {
		return nil, err
	}

	// 7. 转换为响应格式返回
	return s.toLinkResp(link, grave.Name), nil
}

// Validate 验证分享链接
// 这是公开接口，用户通过分享链接访问时调用
// 验证链接是否有效（存在、未过期、未达到使用上限）
// 参数:
//   - shareCode: 分享码（32位十六进制字符串）
//
// 返回:
//   - GraveShareLinkValidResp: 验证通过的墓园信息
//   - error: 错误信息
func (s *GraveShareLinkService) Validate(shareCode string) (*types.GraveShareLinkValidResp, error) {
	// 1. 根据分享码查询分享链接
	link, err := s.linkRepo.GetByShareCode(shareCode)
	if err != nil {
		return nil, err
	}

	// 2. 检查链接是否被停用
	if !link.IsActive {
		return nil, errors.ErrShareLinkExpired
	}

	// 3. 检查是否已过期
	if link.ExpireAt != nil && time.Now().After(*link.ExpireAt) {
		return nil, errors.ErrShareLinkExpired
	}

	// 4. 检查是否超过使用次数限制
	if link.MaxUseCount > 0 && link.UsedCount >= link.MaxUseCount {
		return nil, errors.ErrShareLinkExpired
	}

	// 5. 获取墓园信息
	grave, err := s.graveRepo.GetByID(link.GraveID)
	if err != nil {
		return nil, err
	}

	// 6. 获取逝者列表（分享链接即预授权，允许查看逝者列表）
	deceasedList, _, err := s.deceasedRepo.ListByGrave(link.GraveID, 1, 0)
	if err != nil {
		return nil, err
	}
	deceased := make([]types.DeceasedBrief, 0, len(deceasedList))
	for _, d := range deceasedList {
		deceased = append(deceased, types.DeceasedBrief{
			ID:     d.ID,
			Name:   d.Name,
			Gender: d.Gender,
		})
	}

	// 7. 增加使用计数（异步更新，不阻塞返回）
	_ = s.linkRepo.UpdateUsedCount(link.ID)

	// 8. 返回墓园基本信息
	return &types.GraveShareLinkValidResp{
		GraveID:   link.GraveID,
		GraveName: grave.Name,
		GraveDesc: grave.Description,
		CreatedBy: link.CreatedBy,
		Deceased:  deceased,
	}, nil
}

// Deactivate 停用分享链接
// 只有墓地创建者可以停用分享链接
// 停用后，该链接将无法再被使用
// 参数:
//   - linkID: 分享链接ID
//   - userID: 用户ID（必须是墓地创建者）
//
// 返回:
//   - error: 错误信息
func (s *GraveShareLinkService) Deactivate(linkID, userID uint) error {
	// 1. 获取分享链接记录
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return err
	}

	// 2. 获取墓园信息（用于权限验证）
	grave, err := s.graveRepo.GetByID(link.GraveID)
	if err != nil {
		return err
	}

	// 3. 权限验证：只有墓地创建者可以停用链接
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return errors.ErrGraveAccessDenied
	}

	// 4. 停用链接
	return s.linkRepo.Deactivate(linkID)
}

// ListByGrave 获取墓园的所有分享链接
// 只有墓地创建者可以查看
// 用于在管理页面显示所有有效的分享链接
// 参数:
//   - graveID: 墓地ID
//   - userID: 用户ID（必须是墓地创建者）
//
// 返回:
//   - GraveShareLinkListResp: 分享链接列表响应
//   - error: 错误信息
func (s *GraveShareLinkService) ListByGrave(graveID, userID uint, page, size int) (*types.GraveShareLinkListResp, int64, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, 0, err
	}

	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.linkRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}

	items := make([]types.GraveShareLinkResp, 0, len(list))
	for _, l := range list {
		resp := s.toLinkResp(&l, grave.Name)
		items = append(items, *resp)
	}

	return &types.GraveShareLinkListResp{Items: items}, total, nil
}

// Join 通过分享链接加入墓园
// 分享链接即 Owner 对该用户的提前背书，加入后直接成为成员并确认关系，无需审批
// 参数:
//   - shareCode: 分享码
//   - userID: 用户ID
//   - req: 加入请求（逝者ID、关系类型、备注）
//
// 返回:
//   - GraveShareLinkValidResp: 墓园基本信息
//   - error: 错误信息
func (s *GraveShareLinkService) Join(shareCode string, userID uint, req types.JoinShareLinkReq) (*types.GraveShareLinkValidResp, error) {
	link, err := s.linkRepo.GetByShareCode(shareCode)
	if err != nil {
		return nil, err
	}

	if !link.IsActive {
		return nil, errors.ErrShareLinkExpired
	}

	if link.ExpireAt != nil && time.Now().After(*link.ExpireAt) {
		return nil, errors.ErrShareLinkExpired
	}

	if link.MaxUseCount > 0 && link.UsedCount >= link.MaxUseCount {
		return nil, errors.ErrShareLinkExpired
	}

	grave, err := s.graveRepo.GetByID(link.GraveID)
	if err != nil {
		return nil, err
	}

	_, err = s.memberRepo.GetByUserAndGrave(userID, link.GraveID)
	if err != nil && err != errors.ErrMemberNotFound {
		return nil, err
	}

	if err == errors.ErrMemberNotFound {
		member := &membermodel.GraveMember{
			GraveID: link.GraveID,
			UserID:  userID,
			Role:    membermodel.MemberRoleMember,
		}
		if err := s.memberRepo.Create(member); err != nil {
			return nil, err
		}
	}

	deceased, err := s.deceasedRepo.GetByID(req.DeceasedID)
	if err != nil {
		return nil, errors.ErrDeceasedNotFound
	}
	if deceased.GraveID != link.GraveID {
		return nil, errors.ErrDeceasedNotInGrave
	}

	if grave.GraveMode == 2 {
		exists, err := s.krRepo.CheckExists(userID, req.DeceasedID)
		if err != nil {
			return nil, err
		}
		if !exists {
			relation := &krmodel.KinshipRelation{
				GraveID:            link.GraveID,
				UserID:             userID,
				DeceasedID:         req.DeceasedID,
				RelationType:       req.RelationType,
				ConfirmationStatus: krmodel.ConfirmationConfirmed,
				Note:               req.Note,
			}
			if err := s.krRepo.Create(relation); err != nil {
				return nil, err
			}
		}
	} else {
		exists, err := s.mrRepo.CheckExists(userID, req.DeceasedID)
		if err != nil {
			return nil, err
		}
		if !exists {
			relation := &mrmodel.MemorialRelation{
				GraveID:            link.GraveID,
				UserID:             userID,
				DeceasedID:         req.DeceasedID,
				RelationType:       req.RelationType,
				ConfirmationStatus: mrmodel.MemConfirmationConfirmed,
				Note:               req.Note,
			}
			if err := s.mrRepo.Create(relation); err != nil {
				return nil, err
			}
		}
	}

	_ = s.linkRepo.UpdateUsedCount(link.ID)

	return &types.GraveShareLinkValidResp{
		GraveID:   link.GraveID,
		GraveName: grave.Name,
		GraveDesc: grave.Description,
		CreatedBy: link.CreatedBy,
	}, nil
}

// generateShareCode 生成唯一的分享码
// 使用加密随机数生成器生成32位十六进制字符串
// 返回:
//   - string: 分享码
//   - error: 错误信息
func (s *GraveShareLinkService) generateShareCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// toLinkResp 将数据库模型转换为响应格式
// 内部方法，用于统一处理数据转换
// 参数:
//   - link: 分享链接数据库模型
//   - graveName: 墓园名称
//
// 返回:
//   - GraveShareLinkResp: 响应格式的分享链接数据
func (s *GraveShareLinkService) toLinkResp(link *gslmodel.GraveShareLink, graveName string) *types.GraveShareLinkResp {
	return &types.GraveShareLinkResp{
		ID:          link.ID,
		GraveID:     link.GraveID,
		GraveName:   graveName,
		ShareCode:   link.ShareCode,
		ShareURL:    "", // 前端会根据 shareCode 生成完整URL
		ExpireAt:    link.ExpireAt,
		MaxUseCount: link.MaxUseCount,
		UsedCount:   link.UsedCount,
		IsActive:    link.IsActive,
		CreatedAt:   link.CreatedAt.Format(utils.DateTimeFormat),
	}
}
