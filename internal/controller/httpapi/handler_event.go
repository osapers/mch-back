package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osapers/mch-back/internal/types"
	"golang.org/x/sync/errgroup"
)

type listEventsRes []*eventRes

type eventRes struct {
	*types.Event
	IsParticipant bool `json:"is_participant"`
}

func (s *Server) listEventsHandler(c *gin.Context) {
	var (
		events                      []*types.Event
		eventsWithUserParticipation map[string]struct{}
	)

	eg, egCtx := errgroup.WithContext(c.Request.Context())
	eg.Go(func() error {
		var err error
		events, err = s.es.List(egCtx)
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		var err error
		eventsWithUserParticipation, err = s.us.GetParticipatedEvents(egCtx)
		if err != nil {
			return err
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		return
	}

	res := make(listEventsRes, len(events))

	for i, event := range events {
		_, isParticipant := eventsWithUserParticipation[event.ID]
		res[i] = &eventRes{Event: event, IsParticipant: isParticipant}
	}

	c.JSON(http.StatusOK, jsonResp(res, nil))
}
