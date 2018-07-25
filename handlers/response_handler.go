package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"net/url"
	"time"
)

func toJSTString(now time.Time) string {
	return now.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05")
}

func toEmptyResponse() *pb.Empty {
	return &pb.Empty{}
}

func toUploadURLResponse(from *url.URL) *pb.UploadURL {
	return &pb.UploadURL{Url: from.String()}
}

func toCurrencyResponse(from *models.Currency) *pb.CurrencyResponse {
	return &pb.CurrencyResponse{
		Code:    string(from.Code),
		Name:    from.Name,
		IconURL: from.IconURL.String(),
	}
}

func toCurrencyListResponse(from []*models.Currency) *pb.CurrencyListResponse {
	var items []*pb.CurrencyResponse
	for _, v := range from {
		items = append(items, toCurrencyResponse(v))
	}
	return &pb.CurrencyListResponse{Items: items}
}

func toCurrencyPriceResponse(from *models.CurrencyPrice) *pb.CurrencyPriceResponse {
	return &pb.CurrencyPriceResponse{
		Id:           string(from.Id),
		CurrencyCode: string(from.CurrencyCode),
		Usd:          from.USD,
		Jpy:          int64(from.JPY),
		Datetime:     toJSTString(from.Datetime),
	}
}

func toCurrencyPriceListResponse(from []*models.CurrencyPrice) *pb.CurrencyPriceListResponse {
	var items []*pb.CurrencyPriceResponse
	for _, v := range from {
		items = append(items, toCurrencyPriceResponse(v))
	}
	return &pb.CurrencyPriceListResponse{Items: items}
}

func toUserResponse(from *models.User) *pb.UserResponse {
	u := "https://storage.cloud.google.com/crypto-assistant-dev.appspot.com/users/default-profile.jpg"
	if from.IconURL != nil && from.IconURL.String() != "" {
		u = from.IconURL.String()
	}

	return &pb.UserResponse{
		Id:       string(from.Id),
		Username: from.Username,
		IconURL:  u,
	}
}

func toUserListResponse(from []*models.User) *pb.UserListResponse {
	var items []*pb.UserResponse
	for _, v := range from {
		items = append(items, toUserResponse(v))
	}
	return &pb.UserListResponse{Items: items}
}

func toAddressResponse(from *models.Address) *pb.AddressResponse {
	return &pb.AddressResponse{
		Id:           string(from.Id),
		UserId:       string(from.UserId),
		CurrencyCode: string(from.CurrencyCode),
		Value:        from.Value,
	}
}

func toAddressListResponse(from []*models.Address) *pb.AddressListResponse {
	var items []*pb.AddressResponse
	for _, v := range from {
		items = append(items, toAddressResponse(v))
	}
	return &pb.AddressListResponse{Items: items}
}

func toAssetResponse(amount float64) *pb.AssetResponse {
	return &pb.AssetResponse{
		Amount: float32(amount),
	}
}

func toPortfolioListResponse(from []*models.Portfolio) *pb.PortfolioListResponse {
	var items []*pb.PortfolioResponse
	for _, v := range from {
		items = append(items, toPortfolioResponse(v))
	}
	return &pb.PortfolioListResponse{Items: items}
}

func toPortfolioResponse(from *models.Portfolio) *pb.PortfolioResponse {
	return &pb.PortfolioResponse{
		CurrencyCode: string(from.CurrencyCode),
		Amount:       float32(from.Amount),
		JpyAsset:     float32(from.JPYAsset),
	}
}

func toChartResponse(from []*models.DataPoint) *pb.ChartResponse {
	var items []*pb.DataPointResponse
	for _, v := range from {
		items = append(items, toDataPointResponse(v))
	}
	return &pb.ChartResponse{Items: items}
}

func toDataPointResponse(from *models.DataPoint) *pb.DataPointResponse {
	return &pb.DataPointResponse{
		Datetime: toJSTString(from.Datetime),
		Usd:      from.USD,
		Jpy:      int64(from.JPY),
	}
}
