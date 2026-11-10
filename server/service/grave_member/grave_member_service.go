package grave_member

import (
	"server/database/connection"
	"server/errors"
	membermodel "server/model/grave_member"
	usermodel "server/model/user"
	memberrepo "server/repository/grave_member"
	userrepo "server/repository/user"
	"server/types"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GraveMemberService struct {
	memberRepo *memberrepo.GraveMemberRepository
	userRepo   *userrepo.UserRepository
}

func NewGraveMemberService(memberRepo *memberrepo.GraveMemberRepository, userRepo *userrepo.UserRepository) *GraveMemberService {
	return &GraveMemberService{memberRepo: memberRepo, userRepo: userRepo}
}

func (s *GraveMemberService) Add(graveID, operatorID uint, targetUserID uint, role int8, relationDesc string) (*types.GraveMemberResp, error) {
	operator, err := s.memberRepo.GetByUserAndGrave(operatorID, graveID)
	if err != nil {
		return nil, err
	}
	if operator.Role != membermodel.MemberRoleOwner && operator.Role != membermodel.MemberRoleAdmin {
		return nil, errors.ErrMemberRoleInsufficient
	}

	existing, err := s.memberRepo.GetByUserAndGrave(targetUserID, graveID)
	if err == nil && existing != nil {
		return nil, errors.ErrMemberAlreadyExist
	}

	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	member := &membermodel.GraveMember{
		GraveID:      graveID,
		UserID:       targetUserID,
		Role:         role,
		RelationDesc: relationDesc,
	}
	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return toMemberResp(member, targetUser), nil
}

func (s *GraveMemberService) ListByGrave(graveID, userID uint, page, size int) (*types.GraveMemberListResp, int64, error) {
	_, err := s.memberRepo.GetByUserAndGrave(userID, graveID)
	if err != nil {
		return nil, 0, err
	}

	list, total, err := s.memberRepo.ListByGrave(graveID, page, size)
	if err != nil {
		return nil, 0, err
	}

	userIDs := make([]uint, 0, len(list))
	for _, m := range list {
		userIDs = append(userIDs, m.UserID)
	}

	userMap := make(map[uint]*usermodel.User)
	for _, uid := range userIDs {
		u, err := s.userRepo.GetByID(uid)
		if err != nil {
			continue
		}
		userMap[uid] = u
	}

	return toMemberListResp(list, userMap), total, nil
}

func (s *GraveMemberService) UpdateRole(operatorID, memberID uint, newRole int8) (*types.GraveMemberResp, error) {
	target, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		return nil, err
	}

	operator, err := s.memberRepo.GetByUserAndGrave(operatorID, target.GraveID)
	if err != nil {
		return nil, err
	}
	if operator.Role != membermodel.MemberRoleOwner {
		return nil, errors.ErrMemberRoleInsufficient
	}

	if target.Role == membermodel.MemberRoleOwner {
		return nil, errors.ErrMemberCannotDemoteOwner
	}

	target.Role = newRole
	if err := s.memberRepo.Update(target); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(target.UserID)
	if err != nil {
		return nil, err
	}
	return toMemberResp(target, user), nil
}

func (s *GraveMemberService) UpdateRelationDesc(operatorID, memberID uint, relationDesc string) (*types.GraveMemberResp, error) {
	target, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		return nil, err
	}

	if target.UserID != operatorID {
		operator, err := s.memberRepo.GetByUserAndGrave(operatorID, target.GraveID)
		if err != nil {
			return nil, err
		}
		if operator.Role != membermodel.MemberRoleOwner && operator.Role != membermodel.MemberRoleAdmin {
			return nil, errors.ErrMemberRoleInsufficient
		}
	}

	target.RelationDesc = relationDesc
	if err := s.memberRepo.Update(target); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(target.UserID)
	if err != nil {
		return nil, err
	}
	return toMemberResp(target, user), nil
}

func (s *GraveMemberService) Remove(operatorID, memberID uint) error {
	target, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		return err
	}

	if target.Role == membermodel.MemberRoleOwner {
		return errors.ErrMemberCannotRemoveOwner
	}

	if target.UserID != operatorID {
		operator, err := s.memberRepo.GetByUserAndGrave(operatorID, target.GraveID)
		if err != nil {
			return err
		}
		if operator.Role != membermodel.MemberRoleOwner && operator.Role != membermodel.MemberRoleAdmin {
			return errors.ErrMemberRoleInsufficient
		}
	}

	return s.memberRepo.Delete(memberID)
}

func (s *GraveMemberService) TransferOwnership(ownerID, targetMemberID uint) (*types.GraveMemberResp, *types.GraveMemberResp, error) {
	target, err := s.memberRepo.GetByID(targetMemberID)
	if err != nil {
		return nil, nil, err
	}

	owner, err := s.memberRepo.GetByUserAndGrave(ownerID, target.GraveID)
	if err != nil {
		return nil, nil, err
	}
	if owner.Role != membermodel.MemberRoleOwner {
		return nil, nil, errors.ErrMemberRoleInsufficient
	}

	if target.Role != membermodel.MemberRoleAdmin {
		return nil, nil, errors.ErrMemberTransferTargetInvalid
	}

	owner.Role = membermodel.MemberRoleAdmin
	target.Role = membermodel.MemberRoleOwner

	if err := connection.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.memberRepo.Tx(tx).Update(owner); err != nil {
			return err
		}
		return s.memberRepo.Tx(tx).Update(target)
	}); err != nil {
		return nil, nil, err
	}

	ownerUser, err := s.userRepo.GetByID(owner.UserID)
	if err != nil {
		logrus.Warnf("获取原Owner用户信息失败 (id=%d): %v", owner.UserID, err)
	}
	targetUser, err := s.userRepo.GetByID(target.UserID)
	if err != nil {
		logrus.Warnf("获取新Owner用户信息失败 (id=%d): %v", target.UserID, err)
	}

	return toMemberResp(owner, ownerUser), toMemberResp(target, targetUser), nil
}
