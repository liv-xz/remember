package memorial_relation

import (
	"server/errors"
	mrmodel "server/model/memorial_relation"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	mrrepo "server/repository/memorial_relation"
	"server/types"
	"server/utils"
)

type MemorialRelationService struct {
	mrRepo       *mrrepo.MemorialRelationRepository
	graveRepo    *graverepo.GraveRepository
	memberRepo   *memberrepo.GraveMemberRepository
	deceasedRepo *deceasedrepo.DeceasedRepository
}

func NewMemorialRelationService(
	mrRepo *mrrepo.MemorialRelationRepository,
	graveRepo *graverepo.GraveRepository,
	memberRepo *memberrepo.GraveMemberRepository,
	deceasedRepo *deceasedrepo.DeceasedRepository,
) *MemorialRelationService {
	return &MemorialRelationService{
		mrRepo:       mrRepo,
		graveRepo:    graveRepo,
		memberRepo:   memberRepo,
		deceasedRepo: deceasedRepo,
	}
}

func (s *MemorialRelationService) Create(graveID, userID, deceasedID uint, relationType int8, note string) (*types.MemorialRelationResp, error) {
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

	exists, err := s.mrRepo.CheckExists(userID, deceasedID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrRelationAlreadyExist
	}

	initialStatus := mrmodel.MemConfirmationPending
	if mrmodel.IsAutoConfirmed(relationType) {
		initialStatus = mrmodel.MemConfirmationConfirmed
	}

	relation := &mrmodel.MemorialRelation{
		GraveID:            graveID,
		UserID:             userID,
		DeceasedID:         deceasedID,
		RelationType:       relationType,
		ConfirmationStatus: initialStatus,
		Note:               note,
	}

	if err := s.mrRepo.Create(relation); err != nil {
		return nil, err
	}

	return s.toRelationResp(relation)
}

func (s *MemorialRelationService) Confirm(relationID, userID uint) error {
	relation, err := s.mrRepo.GetByID(relationID)
	if err != nil {
		return err
	}

	grave, err := s.graveRepo.GetByID(relation.GraveID)
	if err != nil {
		return err
	}

	if !s.memberRepo.IsMember(grave.ID, userID) {
		return errors.ErrGraveAccessDenied
	}

	update := &mrmodel.MemorialRelation{
		ConfirmationStatus: mrmodel.MemConfirmationConfirmed,
		ConfirmedBy:        &userID,
	}

	return s.mrRepo.Update(relationID, update)
}

func (s *MemorialRelationService) Reject(relationID, userID uint) error {
	relation, err := s.mrRepo.GetByID(relationID)
	if err != nil {
		return err
	}

	grave, err := s.graveRepo.GetByID(relation.GraveID)
	if err != nil {
		return err
	}

	if !s.memberRepo.IsMember(grave.ID, userID) {
		return errors.ErrGraveAccessDenied
	}

	update := &mrmodel.MemorialRelation{
		ConfirmationStatus: mrmodel.MemConfirmationRejected,
	}

	return s.mrRepo.Update(relationID, update)
}

func (s *MemorialRelationService) ListByGrave(graveID, userID uint, page, size int) (*types.MemorialRelationListResp, int64, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, 0, err
	}

	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.mrRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}

	return s.toRelationListResp(list), total, nil
}

func (s *MemorialRelationService) toRelationResp(relation *mrmodel.MemorialRelation) (*types.MemorialRelationResp, error) {
	userName := relation.User.Nickname
	if userName == "" {
		userName = relation.User.Username
	}
	return &types.MemorialRelationResp{
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

func (s *MemorialRelationService) toRelationListResp(relations []mrmodel.MemorialRelation) *types.MemorialRelationListResp {
	items := make([]types.MemorialRelationResp, 0, len(relations))
	for _, r := range relations {
		resp, err := s.toRelationResp(&r)
		if err != nil {
			continue
		}
		items = append(items, *resp)
	}
	return &types.MemorialRelationListResp{Items: items}
}
