package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/api/jwt/v1/proto"

	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/entity"
)

func ProtoToEntityClaims(p *proto.Claims) (*entity.Claims, error) {
	id, err := uuid.Parse(p.UserId)
	if err != nil {
		return nil, err
	}

	lvl := entity.AccessLevel(p.AccessLevel)
	if !lvl.IsValid() {
		return nil, fmt.Errorf("%w: invalid access level", ErrProtoConvertation)
	}

	return &entity.Claims{
		UserID:    id,
		AccessLvl: entity.AccessLevel(p.AccessLevel),
	}, nil
}

func EntityToProtoClaims(u *entity.Claims) *proto.Claims {
	return &proto.Claims{
		UserId:      u.UserID.String(),
		AccessLevel: int32(u.AccessLvl),
	}
}
