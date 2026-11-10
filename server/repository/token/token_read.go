package repository

import (
	"server/errors"
	rtmodel "server/model/token"
)

func (r *RefreshTokenRepository) GetByHash(hash string) (*rtmodel.RefreshToken, error) {
	var rt rtmodel.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&rt).Error
	if err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, errors.ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return &rt, nil
}
