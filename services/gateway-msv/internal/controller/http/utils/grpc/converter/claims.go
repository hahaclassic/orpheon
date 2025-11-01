package converter

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/api/jwt/v1/proto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
)

var (
	ErrProtoConvertation = errors.New("converter: failed to convert proto to dto")
)

func ProtoToDTOClaims(p *proto.Claims) (*dto.Claims, error) {
	id, err := uuid.Parse(p.UserId)
	if err != nil {
		return nil, err
	}

	lvl := dto.AccessLevel(p.AccessLevel)
	if !lvl.IsValid() {
		return nil, fmt.Errorf("%w: invalid access level", ErrProtoConvertation)
	}

	return &dto.Claims{
		UserID:    id,
		AccessLvl: dto.AccessLevel(p.AccessLevel),
	}, nil
}

func DTOToProtoClaims(u *dto.Claims) *proto.Claims {
	return &proto.Claims{
		UserId:      u.UserID.String(),
		AccessLevel: int32(u.AccessLvl),
	}
}
