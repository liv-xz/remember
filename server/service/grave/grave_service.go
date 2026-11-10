package grave

import (
	"server/database/connection"
	"server/errors"
	gravemodel "server/model/grave"
	membermodel "server/model/grave_member"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	userRepo "server/repository/user"
	"server/types"

	"gorm.io/gorm"
)

// GraveService 使用第二种查询方式：
// Service 通过持有多个 Repository（graveRepo + userRepo），
// 在业务方法中按需组合调用，而不是依赖 GORM Preload 自动关联。
// 优点：职责清晰，每个查询显式可控，不增加 Repo 层的 Preload 方法。
type GraveService struct {
	graveRepo  *graverepo.GraveRepository
	userRepo   *userRepo.UserRepository
	memberRepo *memberrepo.GraveMemberRepository
}

func NewGraveService(graveRepo *graverepo.GraveRepository, userRepo *userRepo.UserRepository, memberRepo *memberrepo.GraveMemberRepository) *GraveService {
	return &GraveService{graveRepo: graveRepo, userRepo: userRepo, memberRepo: memberRepo}
}

type UpdateGraveParams struct {
	Name          *string
	GraveType     *int8
	GraveMode     *int8
	LayoutStyle   *int8
	Description   *string
	CoverImage    *string
	Status        *int8
	Latitude      *float64
	Longitude     *float64
	Address       *string
	ClearLocation bool
}

func (s *GraveService) canAccess(userID uint, graveID uint) bool {
	return s.memberRepo.IsMember(graveID, userID)
}

func (s *GraveService) Create(userID uint, name string, graveType int8, graveMode int8, layoutStyle int8, description string, coverImage string, latitude *float64, longitude *float64, address string) (*types.GraveResp, error) {
	var result *types.GraveResp

	err := connection.DB.Transaction(func(tx *gorm.DB) error {
		grave := &gravemodel.Grave{
			Name:        name,
			GraveType:   graveType,
			GraveMode:   graveMode,
			LayoutStyle: layoutStyle,
			Description: description,
			CoverImage:  coverImage,
			Status:      gravemodel.GraveStatusActive,
			Latitude:    latitude,
			Longitude:   longitude,
			Address:     address,
			CreatedBy:   userID,
		}
		if err := s.graveRepo.Tx(tx).Create(grave); err != nil {
			return err
		}

		member := &membermodel.GraveMember{
			GraveID: grave.ID,
			UserID:  userID,
			Role:    membermodel.MemberRoleOwner,
		}
		if err := s.memberRepo.Tx(tx).Create(member); err != nil {
			return err
		}

		creatorName := s.getCreatorName(userID)
		result = toGraveResp(grave, creatorName)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *GraveService) GetByID(userID, graveID uint) (*types.GraveResp, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}
	if !s.canAccess(userID, grave.ID) {
		return nil, errors.ErrGraveAccessDenied
	}

	creatorName := s.getCreatorName(grave.CreatedBy)
	return toGraveResp(grave, creatorName), nil
}

func (s *GraveService) ListMy(userID uint, page, size int) (*types.GraveListResp, int64, error) {
	graves, total, err := s.graveRepo.ListByUserOrMember(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	items := make([]types.GraveResp, 0, len(graves))
	for _, g := range graves {
		creatorName := s.getCreatorName(g.CreatedBy)
		items = append(items, *toGraveResp(&g, creatorName))
	}
	return &types.GraveListResp{Items: items}, total, nil
}

func (s *GraveService) Update(userID, graveID uint, params *UpdateGraveParams) (*types.GraveResp, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}
	if !s.canAccess(userID, grave.ID) {
		return nil, errors.ErrGraveAccessDenied
	}

	if params.Name != nil {
		grave.Name = *params.Name
	}
	if params.GraveType != nil {
		grave.GraveType = *params.GraveType
	}
	if params.GraveMode != nil {
		grave.GraveMode = *params.GraveMode
	}
	if params.LayoutStyle != nil {
		grave.LayoutStyle = *params.LayoutStyle
	}
	if params.Description != nil {
		grave.Description = *params.Description
	}
	if params.CoverImage != nil {
		grave.CoverImage = *params.CoverImage
	}
	if params.Status != nil {
		grave.Status = *params.Status
	}
	if params.Latitude != nil {
		v := *params.Latitude
		grave.Latitude = &v
	}
	if params.Longitude != nil {
		v := *params.Longitude
		grave.Longitude = &v
	}
	if params.Address != nil {
		grave.Address = *params.Address
	}
	if params.ClearLocation {
		grave.Latitude = nil
		grave.Longitude = nil
		grave.Address = ""
	}

	if err := s.graveRepo.Update(grave); err != nil {
		return nil, err
	}

	creatorName := s.getCreatorName(grave.CreatedBy)
	return toGraveResp(grave, creatorName), nil
}

func (s *GraveService) Delete(userID, graveID uint) error {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return err
	}
	if !s.canAccess(userID, grave.ID) {
		return errors.ErrGraveAccessDenied
	}

	return connection.DB.Transaction(func(tx *gorm.DB) error {
		tx.Exec("DELETE FROM kinship_relations WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM memorial_relations WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM kinship_deceased_relations WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM memorial_records WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM memorial_events WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM grave_share_links WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM deceased WHERE grave_id = ?", graveID)
		tx.Exec("DELETE FROM grave_members WHERE grave_id = ?", graveID)
		return tx.Delete(&gravemodel.Grave{}, graveID).Error
	})
}

func (s *GraveService) getCreatorName(userID uint) string {
	creator, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ""
	}
	if creator.Nickname != "" {
		return creator.Nickname
	}
	return creator.Username
}
