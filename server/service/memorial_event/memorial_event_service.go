package memorial_event

import (
	"server/errors"
	eventmodel "server/model/memorial_event"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	eventrepo "server/repository/memorial_event"
	"server/types"
)

type MemorialEventService struct {
	eventRepo    *eventrepo.MemorialEventRepository
	graveRepo    *graverepo.GraveRepository
	memberRepo   *memberrepo.GraveMemberRepository
	deceasedRepo *deceasedrepo.DeceasedRepository
}

func NewMemorialEventService(eventRepo *eventrepo.MemorialEventRepository, graveRepo *graverepo.GraveRepository, memberRepo *memberrepo.GraveMemberRepository, deceasedRepo *deceasedrepo.DeceasedRepository) *MemorialEventService {
	return &MemorialEventService{
		eventRepo:    eventRepo,
		graveRepo:    graveRepo,
		memberRepo:   memberRepo,
		deceasedRepo: deceasedRepo,
	}
}

type UpdateEventParams struct {
	DeceasedID *uint
	EventType  *int8
	EventDate  *string
	Title      *string
	Lunar      *int8
}

func (s *MemorialEventService) Create(graveID, userID uint, eventType int8, deceasedID uint, eventDate, title string, lunar int8) (*types.MemorialEventResp, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	deceased, err := s.deceasedRepo.GetByID(deceasedID)
	if err != nil {
		return nil, err
	}
	if deceased.GraveID != graveID {
		return nil, errors.ErrDeceasedNotInGrave
	}

	event := &eventmodel.MemorialEvent{
		GraveID:    graveID,
		DeceasedID: deceasedID,
		EventType:  eventType,
		EventDate:  eventDate,
		Title:      title,
		Lunar:      lunar,
	}
	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}
	return toEventResp(event), nil
}

func (s *MemorialEventService) ListByGrave(graveID, userID uint, page, size int) (*types.MemorialEventListResp, int64, error) {
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, 0, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, 0, errors.ErrGraveAccessDenied
	}

	list, total, err := s.eventRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toEventListResp(list), total, nil
}

func (s *MemorialEventService) ListByDeceased(deceasedID, userID uint, page, size int) (*types.MemorialEventListResp, int64, error) {
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

	list, total, err := s.eventRepo.ListByDeceased(deceasedID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toEventListResp(list), total, nil
}

func (s *MemorialEventService) Update(eventID, userID uint, params *UpdateEventParams) (*types.MemorialEventResp, error) {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	grave, err := s.graveRepo.GetByID(event.GraveID)
	if err != nil {
		return nil, err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	if params.DeceasedID != nil {
		deceased, err := s.deceasedRepo.GetByID(*params.DeceasedID)
		if err != nil {
			return nil, err
		}
		if deceased.GraveID != event.GraveID {
			return nil, errors.ErrDeceasedNotInGrave
		}
		event.DeceasedID = *params.DeceasedID
	}
	if params.EventType != nil {
		event.EventType = *params.EventType
	}
	if params.EventDate != nil {
		event.EventDate = *params.EventDate
	}
	if params.Title != nil {
		event.Title = *params.Title
	}
	if params.Lunar != nil {
		event.Lunar = *params.Lunar
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	return toEventResp(event), nil
}

func (s *MemorialEventService) Delete(eventID, userID uint) error {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	grave, err := s.graveRepo.GetByID(event.GraveID)
	if err != nil {
		return err
	}
	if !s.memberRepo.IsMember(grave.ID, userID) {
		return errors.ErrGraveAccessDenied
	}

	return s.eventRepo.Delete(eventID)
}
