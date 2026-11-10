package kinship_relation

import (
	"server/database/connection"
	"server/errors"
	membermodel "server/model/grave_member"
	krmodel "server/model/kinship_relation"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	krrepo "server/repository/kinship_relation"
	"server/types"
	"server/utils"

	"gorm.io/gorm"
)

type KinshipRelationService struct {
	krRepo       *krrepo.KinshipRelationRepository
	graveRepo    *graverepo.GraveRepository
	deceasedRepo *deceasedrepo.DeceasedRepository
	memberRepo   *memberrepo.GraveMemberRepository
}

func NewKinshipRelationService(
	krRepo *krrepo.KinshipRelationRepository,
	graveRepo *graverepo.GraveRepository,
	deceasedRepo *deceasedrepo.DeceasedRepository,
	memberRepo *memberrepo.GraveMemberRepository,
) *KinshipRelationService {
	return &KinshipRelationService{
		krRepo:       krRepo,
		graveRepo:    graveRepo,
		deceasedRepo: deceasedRepo,
		memberRepo:   memberRepo,
	}
}

func (s *KinshipRelationService) CreatePublic(graveID, userID, deceasedID uint, relationType int8, note string) (*types.KinshipRelationResp, error) {
	_, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}

	deceased, err := s.deceasedRepo.GetByID(deceasedID)
	if err != nil {
		return nil, errors.ErrDeceasedNotFound
	}
	if deceased.GraveID != graveID {
		return nil, errors.ErrDeceasedNotInGrave
	}

	exists, err := s.krRepo.CheckExists(userID, deceasedID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrRelationAlreadyExist
	}

	relation := &krmodel.KinshipRelation{
		GraveID:            graveID,
		UserID:             userID,
		DeceasedID:         deceasedID,
		RelationType:       relationType,
		ConfirmationStatus: krmodel.ConfirmationPending,
		Note:               note,
	}

	if err := s.krRepo.Create(relation); err != nil {
		return nil, err
	}

	return s.toRelationResp(relation)
}

func (s *KinshipRelationService) Confirm(relationID, userID uint) error {
	relation, err := s.krRepo.GetByID(relationID)
	if err != nil {
		return err
	}

	if !s.memberRepo.IsMember(relation.GraveID, userID) {
		return errors.ErrGraveAccessDenied
	}

	return connection.DB.Transaction(func(tx *gorm.DB) error {
		update := &krmodel.KinshipRelation{
			ConfirmationStatus: krmodel.ConfirmationConfirmed,
			ConfirmedBy:        &userID,
		}
		if err := s.krRepo.Tx(tx).Update(relationID, update); err != nil {
			return err
		}

		existing, err := s.memberRepo.Tx(tx).GetByUserAndGrave(relation.UserID, relation.GraveID)
		if err != nil {
			return err
		}
		if existing != nil {
			return nil
		}

		member := &membermodel.GraveMember{
			GraveID:      relation.GraveID,
			UserID:       relation.UserID,
			Role:         membermodel.MemberRoleMember,
			RelationDesc: "",
		}
		return s.memberRepo.Tx(tx).Create(member)
	})
}

func (s *KinshipRelationService) Reject(relationID, userID uint) error {
	relation, err := s.krRepo.GetByID(relationID)
	if err != nil {
		return err
	}

	if !s.memberRepo.IsMember(relation.GraveID, userID) {
		return errors.ErrGraveAccessDenied
	}

	update := &krmodel.KinshipRelation{
		ConfirmationStatus: krmodel.ConfirmationRejected,
	}

	return s.krRepo.Update(relationID, update)
}

func (s *KinshipRelationService) ListByGrave(graveID, userID uint, page, size int) (*types.KinshipRelationListResp, int64, error) {
	if !s.memberRepo.IsMember(graveID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.krRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}

	return s.toRelationListResp(list), total, nil
}

func (s *KinshipRelationService) ListByDeceased(deceasedID, userID uint, page, size int) (*types.KinshipRelationListResp, int64, error) {
	deceased, err := s.deceasedRepo.GetByID(deceasedID)
	if err != nil {
		return nil, 0, err
	}

	if !s.memberRepo.IsMember(deceased.GraveID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.krRepo.ListByDeceased(deceasedID, page, size)
	if err != nil {
		return nil, 0, err
	}

	return s.toRelationListResp(list), total, nil
}

func (s *KinshipRelationService) toRelationResp(relation *krmodel.KinshipRelation) (*types.KinshipRelationResp, error) {
	userName := relation.User.Nickname
	if userName == "" {
		userName = relation.User.Username
	}
	return &types.KinshipRelationResp{
		ID:                 relation.ID,
		GraveID:            relation.GraveID,
		UserID:             relation.UserID,
		UserName:           userName,
		UserAvatar:         relation.User.Avatar,
		DeceasedID:         relation.DeceasedID,
		DeceasedName:       relation.Deceased.Name,
		DeceasedAvatar:     relation.Deceased.PortraitURL,
		RelationType:       relation.RelationType,
		ConfirmationStatus: relation.ConfirmationStatus,
		ConfirmedBy:        relation.ConfirmedBy,
		Note:               relation.Note,
		CreatedAt:          relation.CreatedAt.Format(utils.DateTimeFormat),
		UpdatedAt:          relation.UpdatedAt.Format(utils.DateTimeFormat),
	}, nil
}

func (s *KinshipRelationService) toRelationListResp(relations []krmodel.KinshipRelation) *types.KinshipRelationListResp {
	items := make([]types.KinshipRelationResp, 0, len(relations))
	for _, r := range relations {
		resp, err := s.toRelationResp(&r)
		if err != nil {
			continue
		}
		items = append(items, *resp)
	}
	return &types.KinshipRelationListResp{Items: items}
}
