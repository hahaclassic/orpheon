package converter

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrNilUser = errors.New("converter: user is nil")
)

func ProtoToEntityUser(p *proto.User) (*entity.User, error) {
	if p == nil {
		return nil, ErrNilUser
	}
	id, err := uuid.Parse(p.Id)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:               id,
		Name:             p.Name,
		AccessLvl:        entity.AccessLevel(p.AccessLevel),
		RegistrationDate: p.RegistrationDate.AsTime(),
	}, nil
}

func EntityToProtoUser(u *entity.User) *proto.User {
	return &proto.User{
		Id:               u.ID.String(),
		Name:             u.Name,
		AccessLevel:      proto.AccessLevel(u.AccessLvl),
		RegistrationDate: timestamppb.New(u.RegistrationDate),
	}
}
