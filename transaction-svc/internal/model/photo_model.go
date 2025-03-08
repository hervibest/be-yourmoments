package model

// TODO add similarity
type RequestUpdateProcessedPhoto struct {
	Id                     string
	PreviewUrl             string
	PreviewWithBoundingUrl string
	UserId                 []string
}

type RequestClaimPhoto struct {
	Id     string
	UserId string
}
