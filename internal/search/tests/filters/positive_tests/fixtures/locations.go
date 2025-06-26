package fixtures

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"

var LocationsFixture = []dto.Location{
	{
		Id:            1,
		ZipCode:       "40601",
		Country:       dto.Country{Code: "US"},
		State:         dto.State{Code: "KY"},
		City:          "Frankfort",
		Address:       dto.Address{Line1: "2001 Leestown Rd"},
		IsResidential: false,
		Name:          "Frankfort, KY 40601",
		Phone:         "",
	},
	{
		Id:            97,
		ZipCode:       "N3V 6T1",
		Country:       dto.Country{Name: "CA", Code: "CA"},
		State:         dto.State{Name: "ON", Code: "ON"},
		City:          "Brantford",
		Address:       dto.Address{Line1: "470 Hardy Rd"},
		IsResidential: false,
		Phone:         "",
		Name:          "Brantford, ON N3V 6T1",
	},
	{
		Id:            64,
		ZipCode:       "76140",
		Country:       dto.Country{Name: "US", Code: "US"},
		State:         dto.State{Name: "TX", Code: "TX"},
		City:          "Fort Worth",
		Address:       dto.Address{Line1: "7100 Oak Grove Road"},
		IsResidential: false,
		Phone:         "",
		Name:          "Fort Worth D, TX 76140",
	},
	{
		Id:            10,
		ZipCode:       "72713",
		Country:       dto.Country{Name: "US", Code: "US"},
		State:         dto.State{Name: "AR", Code: "AR"},
		City:          "Bentonville",
		Address:       dto.Address{Line1: "4900 SW Regional Airport Blvd"},
		IsResidential: false,
		Phone:         "",
		Name:          "Bentonville, AR 72713",
	},
}
