package handler

import (
	"context"
	"property-service/internal/proto"
	"property-service/internal/repository"
	"property-service/internal/service"
)

type GRPCHandler struct {
	proto.UnimplementedPropertyServiceServer
	service *service.PropertyService
}

func New(service *service.PropertyService) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (h *GRPCHandler) CreateProperty(ctx context.Context, req *proto.CreatePropertyRequest) (*proto.CreatePropertyResponse, error) {
	p := &repository.Property{
		OwnerID:      req.OwnerId,
		Title:        req.Title,
		Description:  req.Description,
		City:         req.City,
		AddressLine:  req.AddressLine,
		Lat:          float64(req.Lat),
		Lng:          float64(req.Lng),
		PropertyType: req.PropertyType,
		Rooms:        req.Rooms,
		Area:         float64(req.Area),
		Floor:        req.Floor,
		TotalFloors:  req.TotalFloors,
		PricePerMonth: int32(req.PricePerMonth),
		Currency:     req.Currency,
		MainImageURL: req.MainImageUrl,
		ImageURLs:    req.ImageUrls,
		HasWiFi:      req.HasWifi,
		HasParking:   req.HasParking,
		HasElevator:  req.HasElevator,
		IsVerified:   req.IsVerified,
		Availability: mapAvailability(req.Availability),
		Status:       req.Status,
	}
	id, err := h.service.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	return &proto.CreatePropertyResponse{Id: id}, nil
}

func (h *GRPCHandler) UpdateProperty(ctx context.Context, req *proto.UpdatePropertyRequest) (*proto.EmptyResponse, error) {
	p := &repository.Property{
		OwnerID:      req.OwnerId,
		Title:        req.Title,
		Description:  req.Description,
		City:         req.City,
		AddressLine:  req.AddressLine,
		Lat:          float64(req.Lat),
		Lng:          float64(req.Lng),
		PropertyType: req.PropertyType,
		Rooms:        req.Rooms,
		Area:         float64(req.Area),
		Floor:        req.Floor,
		TotalFloors:  req.TotalFloors,
		PricePerMonth: int32(req.PricePerMonth),
		Currency:     req.Currency,
		MainImageURL: req.MainImageUrl,
		ImageURLs:    req.ImageUrls,
		HasWiFi:      req.HasWifi,
		HasParking:   req.HasParking,
		HasElevator:  req.HasElevator,
		IsVerified:   req.IsVerified,
		Availability: mapAvailability(req.Availability),
		Status:       req.Status,
	}
	return &proto.EmptyResponse{}, h.service.Update(ctx, req.Id, p)
}

func (h *GRPCHandler) GetProperty(ctx context.Context, req *proto.GetPropertyRequest) (*proto.GetPropertyResponse, error) {
	p, err := h.service.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetPropertyResponse{
		Property: &proto.Property{
			Id:           p.ID,
			OwnerId:      p.OwnerID,
			Title:        p.Title,
			Description:  p.Description,
			City:         p.City,
			AddressLine:  p.AddressLine,
			Lat:          float32(p.Lat),
			Lng:          float32(p.Lng),
			PropertyType: p.PropertyType,
			Rooms:        p.Rooms,
			Area:         float32(p.Area),
			Floor:        p.Floor,
			TotalFloors:  p.TotalFloors,
			PricePerMonth: int32(p.PricePerMonth),
			Currency:     p.Currency,
			MainImageUrl: p.MainImageURL,
			ImageUrls:    p.ImageURLs,
			HasWifi:      p.HasWiFi,
			HasParking:   p.HasParking,
			HasElevator:  p.HasElevator,
			IsVerified:   p.IsVerified,
			Rating:       float32(p.Rating),
			ReviewsCount: p.ReviewsCount,
			Availability: mapAvailabilityProto(p.Availability),
			Status:       p.Status,
			CreatedAt:    p.CreatedAt.String(),
			UpdatedAt:    p.UpdatedAt.String(),
			},
	}, nil
}

func (h *GRPCHandler) DeleteProperty(ctx context.Context, req *proto.DeletePropertyRequest) (*proto.EmptyResponse, error) {
	return &proto.EmptyResponse{}, h.service.Delete(ctx, req.Id)
}

func mapAvailabilityProto(a []repository.AvailabilityPeriod) []*proto.AvailabilityPeriod {
	res := make([]*proto.AvailabilityPeriod, 0, len(a))
	for _, v := range a {
		res = append(res, &proto.AvailabilityPeriod{
			FromDate: v.FromDate,
			ToDate:   v.ToDate,
		})
	}
	return res
}

func mapAvailability(availability []*proto.AvailabilityPeriod) []repository.AvailabilityPeriod {
	availabilityPeriods := make([]repository.AvailabilityPeriod, len(availability))
	for i, period := range availability {
		availabilityPeriods[i] = repository.AvailabilityPeriod{
			FromDate: period.FromDate,
			ToDate:   period.ToDate,
		}
	}
	return availabilityPeriods
}