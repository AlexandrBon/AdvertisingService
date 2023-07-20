package grpc

import (
	"advertisingService/internal/adApp"
	"advertisingService/internal/ads"
	"advertisingService/internal/userApp"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"strings"
)

func getGRPCStatus(err error) codes.Code {
	if errors.Is(err, adApp.ErrForbidden) {
		return codes.PermissionDenied
	} else if errors.Is(err, adApp.ErrBadRequest) {
		return codes.Internal
	} else if errors.Is(err, userApp.ErrEmailConflict) {
		return codes.AlreadyExists
	} else if errors.Is(err, userApp.ErrBadRequest) {
		return codes.Internal
	}
	return codes.Unknown
}

type AdService struct {
	adApp   adApp.App
	userApp userApp.App
}

func NewServiceServer(adApp adApp.App, userApp userApp.App) AdServiceServer {
	return &AdService{adApp: adApp, userApp: userApp}
}

func (as *AdService) CreateAd(_ context.Context, req *CreateAdRequest) (*AdResponse, error) {
	ad, err := as.adApp.CreateAd(req.Title, req.Text, req.UserId)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title, Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published}, status.Error(codes.OK, "")
}

func (as *AdService) ChangeAdStatus(_ context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := as.adApp.ChangeAdStatus(req.AdId, req.UserId, req.Published)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title, Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published}, status.Error(codes.OK, "")
}

func (as *AdService) UpdateAd(_ context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	ad, err := as.adApp.UpdateAd(req.AdId, req.UserId, req.Title, req.Text)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title, Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published}, status.Error(codes.OK, "")
}

func (as *AdService) ListAds(_ context.Context, req *ListAdRequest) (*ListAdResponse, error) {
	var filter []func(ads.Ad) bool

	if value, ok := req.Filter["published"]; ok {
		filter = append(filter, func(ad ads.Ad) bool {
			pub, _ := strconv.ParseBool(value)
			return ad.Published == pub
		})
	}

	if value, ok := req.Filter["authorID"]; ok {
		filter = append(filter, func(ad ads.Ad) bool {
			userID, _ := strconv.ParseInt(value, 10, 64)
			return ad.AuthorID == userID
		})
	}

	if value, ok := req.Filter["creationDate"]; ok {
		filter = append(filter, func(ad ads.Ad) bool {
			return ad.CreationDate >= value
		})
	}

	if value, ok := req.Filter["title"]; ok {
		filter = append(filter, func(ad ads.Ad) bool {
			return strings.Contains(ad.Title, value)
		})
	}

	adList, err := as.adApp.ListAds(filter)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}

	var adsResponse []*AdResponse
	for _, ad := range adList {
		adsResponse = append(adsResponse, &AdResponse{
			Id:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationTime: ad.CreationDate,
		})
	}

	return &ListAdResponse{List: adsResponse}, status.Error(codes.OK, "")
}

func (as *AdService) CreateUser(_ context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user, err := as.userApp.CreateUser(req.Name, req.Email)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &UserResponse{Id: user.ID, Name: user.Nickname, Email: user.Email}, status.Error(codes.OK, "")
}

func (as *AdService) GetUser(_ context.Context, req *GetUserRequest) (*UserResponse, error) {
	user, err := as.userApp.GetUser(req.Id)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &UserResponse{Id: user.ID, Name: user.Nickname, Email: user.Email}, status.Error(codes.OK, "")
}

func (as *AdService) DeleteUser(_ context.Context, req *DeleteUserRequest) (*emptypb.Empty, error) {
	err := as.userApp.DeleteUser(req.Id)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &emptypb.Empty{}, status.Error(codes.OK, "")
}

func (as *AdService) DeleteAd(_ context.Context, req *DeleteAdRequest) (*emptypb.Empty, error) {
	err := as.adApp.DeleteAd(req.AdId, req.AuthorId)
	if err != nil {
		return nil, status.Error(getGRPCStatus(err), err.Error())
	}
	return &emptypb.Empty{}, status.Error(codes.OK, "")
}
