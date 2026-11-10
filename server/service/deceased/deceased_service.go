package deceased

import (
	"server/database/connection"
	"server/errors"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
	krmodel "server/model/kinship_relation"
	mrmodel "server/model/memorial_relation"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	kdrrepo "server/repository/kinship_deceased_relation"
	krrepo "server/repository/kinship_relation"
	mrrepo "server/repository/memorial_relation"
	"server/types"

	"gorm.io/gorm"
)

type DeceasedService struct {
	deceasedRepo         *deceasedrepo.DeceasedRepository
	graveRepo            *graverepo.GraveRepository
	memberRepo           *memberrepo.GraveMemberRepository
	krRepo               *krrepo.KinshipRelationRepository
	mrRepo               *mrrepo.MemorialRelationRepository
	deceasedRelationRepo *kdrrepo.DeceasedRelationRepository
}

func NewDeceasedService(deceasedRepo *deceasedrepo.DeceasedRepository, graveRepo *graverepo.GraveRepository, memberRepo *memberrepo.GraveMemberRepository, krRepo *krrepo.KinshipRelationRepository, mrRepo *mrrepo.MemorialRelationRepository, deceasedRelationRepo *kdrrepo.DeceasedRelationRepository) *DeceasedService {
	return &DeceasedService{deceasedRepo: deceasedRepo, graveRepo: graveRepo, memberRepo: memberRepo, krRepo: krRepo, mrRepo: mrRepo, deceasedRelationRepo: deceasedRelationRepo}
}

type CreateDeceasedParams struct {
	Name             string
	Gender           int8
	BirthDate        string
	DeathDate        string
	Biography        string
	PortraitURL      string
	Epitaph          string
	BurialOrder      int
	SortOrder        int
	OriginGraveID    *uint
	UserRelationType int8
	UserRelationNote string
}

func (s *DeceasedService) Create(graveID, userID uint, req CreateDeceasedParams) (*types.DeceasedResp, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	if err := s.checkCapacity(grave); err != nil {
		return nil, err
	}

	if err := validateDates(req.BirthDate, req.DeathDate); err != nil {
		return nil, err
	}

	deceased := &deceasedmodel.Deceased{
		GraveID:       graveID,
		Name:          req.Name,
		Gender:        req.Gender,
		BirthDate:     req.BirthDate,
		DeathDate:     req.DeathDate,
		Biography:     req.Biography,
		PortraitURL:   req.PortraitURL,
		Epitaph:       req.Epitaph,
		BurialOrder:   req.BurialOrder,
		SortOrder:     req.SortOrder,
		OriginGraveID: req.OriginGraveID,
	}
	if err := s.deceasedRepo.Create(deceased); err != nil {
		return nil, err
	}

	if req.UserRelationType != 0 {
		confirmedBy := userID
		if grave.GraveMode == gravemodel.GraveModeMemorial {
			rel := &mrmodel.MemorialRelation{
				GraveID:            graveID,
				UserID:             userID,
				DeceasedID:         deceased.ID,
				RelationType:       req.UserRelationType,
				ConfirmationStatus: mrmodel.MemConfirmationConfirmed,
				ConfirmedBy:        &confirmedBy,
				Note:               req.UserRelationNote,
			}
			if err := s.mrRepo.Create(rel); err != nil {
				return nil, err
			}
		} else {
			rel := &krmodel.KinshipRelation{
				GraveID:            graveID,
				UserID:             userID,
				DeceasedID:         deceased.ID,
				RelationType:       req.UserRelationType,
				ConfirmationStatus: krmodel.ConfirmationConfirmed,
				ConfirmedBy:        &confirmedBy,
				Note:               req.UserRelationNote,
			}
			if err := s.krRepo.Create(rel); err != nil {
				return nil, err
			}
		}
	}

	return toDeceasedResp(deceased), nil
}

func (s *DeceasedService) ListByGrave(graveID, userID uint, page, size int) (*types.DeceasedListResp, int64, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, 0, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.deceasedRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toDeceasedListResp(list), total, nil
}

type UpdateDeceasedParams struct {
	Name          *string
	Gender        *int8
	BirthDate     *string
	DeathDate     *string
	Biography     *string
	PortraitURL   *string
	Epitaph       *string
	BurialOrder   *int
	SortOrder     *int
	OriginGraveID *uint
}

func (s *DeceasedService) Update(id, userID uint, params *UpdateDeceasedParams) (*types.DeceasedResp, error) {
	deceased, err := s.deceasedRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	grave, err := s.graveRepo.GetByID(deceased.GraveID)
	if err != nil {
		return nil, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	birthDate := deceased.BirthDate
	deathDate := deceased.DeathDate
	if params.BirthDate != nil {
		birthDate = *params.BirthDate
	}
	if params.DeathDate != nil {
		deathDate = *params.DeathDate
	}
	if err := validateDates(birthDate, deathDate); err != nil {
		return nil, err
	}

	if params.Name != nil {
		deceased.Name = *params.Name
	}
	if params.Gender != nil {
		deceased.Gender = *params.Gender
	}
	if params.BirthDate != nil {
		deceased.BirthDate = *params.BirthDate
	}
	if params.DeathDate != nil {
		deceased.DeathDate = *params.DeathDate
	}
	if params.Biography != nil {
		deceased.Biography = *params.Biography
	}
	if params.PortraitURL != nil {
		deceased.PortraitURL = *params.PortraitURL
	}
	if params.Epitaph != nil {
		deceased.Epitaph = *params.Epitaph
	}
	if params.BurialOrder != nil {
		deceased.BurialOrder = *params.BurialOrder
	}
	if params.SortOrder != nil {
		deceased.SortOrder = *params.SortOrder
	}
	if params.OriginGraveID != nil {
		deceased.OriginGraveID = params.OriginGraveID
	}

	if err := s.deceasedRepo.Update(deceased); err != nil {
		return nil, err
	}
	return toDeceasedResp(deceased), nil
}

func (s *DeceasedService) Delete(id, userID uint) error {
	deceased, err := s.deceasedRepo.GetByID(id)
	if err != nil {
		return err
	}

	grave, err := s.graveRepo.GetByID(deceased.GraveID)
	if err != nil {
		return err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return errors.ErrGraveAccessDenied
	}

	return connection.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM memorial_relations WHERE deceased_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM kinship_relations WHERE deceased_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM kinship_deceased_relations WHERE deceased_id1 = ? OR deceased_id2 = ?", id, id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM memorial_records WHERE deceased_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&deceasedmodel.Deceased{}, id).Error
	})
}

func (s *DeceasedService) checkCapacity(grave *gravemodel.Grave) error {
	count, err := s.deceasedRepo.CountByGrave(grave.ID)
	if err != nil {
		return err
	}
	switch grave.GraveType {
	case gravemodel.GraveTypeSingle:
		if count >= 1 {
			return errors.ErrDeceasedOverLimit
		}
	case gravemodel.GraveTypeCouple:
		if count >= 2 {
			return errors.ErrDeceasedOverLimit
		}
	case gravemodel.GraveTypeFamily:
		return nil
	default:
		return errors.ErrDeceasedOverLimit
	}
	return nil
}

func validateDates(birthDate, deathDate string) error {
	if birthDate != "" && deathDate != "" && deathDate < birthDate {
		return errors.ErrDeceasedDateInvalid
	}
	return nil
}
