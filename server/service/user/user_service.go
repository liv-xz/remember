package user

import (
	"server/errors"
	userRepo "server/repository/user"
	"server/types"
)

type UserService struct {
	userRepo *userRepo.UserRepository
}

func NewUserService(userRepo *userRepo.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID uint) (*types.UserProfileResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &types.UserProfileResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Abstract: user.Abstract,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, nickname, avatar, abstract string) (*types.UserProfileResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if abstract != "" {
		user.Abstract = abstract
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.ErrUserNotFound
	}

	return &types.UserProfileResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Abstract: user.Abstract,
	}, nil
}
