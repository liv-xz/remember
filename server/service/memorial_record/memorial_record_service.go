package memorial_record

import (
	"server/errors"
	recordmodel "server/model/memorial_record"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	recordrepo "server/repository/memorial_record"
	"server/types"
)

type MemorialRecordService struct {
	recordRepo   *recordrepo.MemorialRecordRepository
	graveRepo    *graverepo.GraveRepository
	memberRepo   *memberrepo.GraveMemberRepository
	deceasedRepo *deceasedrepo.DeceasedRepository
}

func NewMemorialRecordService(recordRepo *recordrepo.MemorialRecordRepository, graveRepo *graverepo.GraveRepository, memberRepo *memberrepo.GraveMemberRepository, deceasedRepo *deceasedrepo.DeceasedRepository) *MemorialRecordService {
	return &MemorialRecordService{
		recordRepo:   recordRepo,
		graveRepo:    graveRepo,
		memberRepo:   memberRepo,
		deceasedRepo: deceasedRepo,
	}
}

func (s *MemorialRecordService) Create(graveID, userID uint, recordType int8, deceasedID *uint, message string) (*types.MemorialRecordResp, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	if deceasedID != nil {
		deceased, err := s.deceasedRepo.GetByID(*deceasedID)
		if err != nil {
			return nil, err
		}
		if deceased.GraveID != graveID {
			return nil, errors.ErrDeceasedNotInGrave
		}
	}

	record := &recordmodel.MemorialRecord{
		GraveID:    graveID,
		DeceasedID: deceasedID,
		UserID:     userID,
		Type:       recordType,
		Message:    message,
	}
	if err := s.recordRepo.Create(record); err != nil {
		return nil, err
	}
	return toRecordResp(record), nil
}

func (s *MemorialRecordService) ListByGrave(graveID, userID uint, page, size int) (*types.MemorialRecordListResp, int64, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, 0, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.recordRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toRecordListResp(list), total, nil
}

func (s *MemorialRecordService) ListByDeceased(deceasedID, userID uint, page, size int) (*types.MemorialRecordListResp, int64, error) {
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

	list, total, err := s.recordRepo.ListByDeceased(deceasedID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toRecordListResp(list), total, nil
}

func (s *MemorialRecordService) Delete(recordID, userID uint) error {
	return s.recordRepo.Delete(recordID, userID)
}
