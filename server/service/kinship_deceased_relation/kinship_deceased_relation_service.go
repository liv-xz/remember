package kinship_deceased_relation

import (
	"server/errors"
	kdr "server/model/kinship_deceased_relation"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	kdrrepo "server/repository/kinship_deceased_relation"
	"server/types"
)

type DeceasedRelationService struct {
	relationRepo *kdrrepo.DeceasedRelationRepository
	graveRepo    *graverepo.GraveRepository
	memberRepo   *memberrepo.GraveMemberRepository
	deceasedRepo *deceasedrepo.DeceasedRepository
}

func NewDeceasedRelationService(relationRepo *kdrrepo.DeceasedRelationRepository, graveRepo *graverepo.GraveRepository, memberRepo *memberrepo.GraveMemberRepository, deceasedRepo *deceasedrepo.DeceasedRepository) *DeceasedRelationService {
	return &DeceasedRelationService{
		relationRepo: relationRepo,
		graveRepo:    graveRepo,
		memberRepo:   memberRepo,
		deceasedRepo: deceasedRepo,
	}
}

func (s *DeceasedRelationService) Create(graveID, userID, deceasedID1, deceasedID2 uint, relationType, direction int8) (*types.DeceasedRelationResp, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	if deceasedID1 == deceasedID2 {
		return nil, errors.ErrRelationSelfInvalid
	}

	if !kdr.IsValidDirection(relationType, direction) {
		return nil, errors.ErrRelationDirectionInvalid
	}

	deceased1, err := s.deceasedRepo.GetByID(deceasedID1)
	if err != nil {
		return nil, err
	}
	if deceased1.GraveID != graveID {
		return nil, errors.ErrDeceasedNotInGrave
	}

	deceased2, err := s.deceasedRepo.GetByID(deceasedID2)
	if err != nil {
		return nil, err
	}
	if deceased2.GraveID != graveID {
		return nil, errors.ErrDeceasedNotInGrave
	}

	existing, _, err := s.relationRepo.ListByDeceased(deceasedID1, 1, 0)
	if err != nil {
		return nil, err
	}
	for _, r := range existing {
		if r.DeceasedID1 == deceasedID2 || r.DeceasedID2 == deceasedID2 {
			return nil, errors.ErrRelationAlreadyExist
		}
	}

	relation := &kdr.DeceasedRelation{
		GraveID:      graveID,
		DeceasedID1:  deceasedID1,
		DeceasedID2:  deceasedID2,
		RelationType: relationType,
		Direction:    direction,
	}
	if err := s.relationRepo.Create(relation); err != nil {
		return nil, err
	}
	return toRelationResp(relation), nil
}

func (s *DeceasedRelationService) ListByGrave(graveID, userID uint, page, size int) (*types.DeceasedRelationListResp, int64, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, 0, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.relationRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toRelationListResp(list), total, nil
}

func (s *DeceasedRelationService) ListByDeceased(deceasedID, userID uint, page, size int) (*types.DeceasedRelationListResp, int64, error) {
	deceased, err := s.deceasedRepo.GetByID(deceasedID)
	if err != nil {
		return nil, 0, err
	}

	grave, err := s.graveRepo.GetByID(deceased.GraveID)
	if err != nil {
		return nil, 0, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.relationRepo.ListByDeceased(deceasedID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toRelationListResp(list), total, nil
}

func (s *DeceasedRelationService) Delete(relationID, userID uint) error {
	relation, err := s.relationRepo.GetByID(relationID)
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

	return s.relationRepo.Delete(relationID)
}
