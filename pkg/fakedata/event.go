package fakedata

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/osapers/mch-back/internal/types"
)

func Events(count int) []*types.Event {
	events := make([]*types.Event, count)

	startTime := time.Unix(1623301410, 0)
	endTime := time.Unix(1624597410, 0)

	for i := 0; i < count; i++ {
		events[i] = &types.Event{
			Name:             gofakeit.BeerName(),
			Image:            gofakeit.ImageURL(600, 600),
			Category:         randomCategory(),
			Date:             gofakeit.DateRange(startTime, endTime).Unix(),
			ShortDescription: gofakeit.Sentence(4),
			Description:      gofakeit.HackerPhrase(),
			Address: types.Address{
				Lat: gofakeit.Latitude(),
				Lon: gofakeit.Longitude(),
				Raw: gofakeit.Street(),
			},
			Website: fmt.Sprintf("https://%s.com", gofakeit.HipsterSentence(1)),
			Email:   gofakeit.Email(),
		}
	}

	return events
}

func randomCategory() string {
	for k := range types.EventCategories {
		return k
	}

	return ""
}
