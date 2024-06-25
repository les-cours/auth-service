package server

import (
	context "context"
	"github.com/les-cours/auth-service/api/auth"
	"github.com/les-cours/auth-service/api/users"
	"github.com/les-cours/auth-service/types"
	"github.com/les-cours/auth-service/utils"
)

func (s *Server) Signup(ctx context.Context, in *auth.SignUpRequest) (*auth.SignUpResponse, error) {

	var validUser, err = s.userClient.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: in.AccountID,
		UserRole:  in.UserRole,
	})

	if err != nil {
		//s.Logger.Error(err.Error())
		return nil, err
	}
	var token *types.AuthToken

	switch in.UserRole {
	case "teacher":
		token, err = utils.GenerateTeacherAccessToken(validUser)
	default:
		var grad string
		var gender string
		err = s.DB.QueryRow(`SELECT grade_id, gender FROM students where student_id  = $1;`, validUser.Id).Scan(&grad, &gender)
		if err != nil {
			s.Logger.Error(err.Error())
		}
		token, err = utils.GenerateAccessToken(validUser, grad, gender)
	}

	if err != nil {
		//s.Logger.Error(err.Error())
		return nil, err
	}

	return &auth.SignUpResponse{
		AccessToken: &auth.AuthToken{
			Token:     token.Token,
			ExpiresAt: token.ExpiresIn,
			TokenType: token.TokenType,
		},
	}, nil

}
