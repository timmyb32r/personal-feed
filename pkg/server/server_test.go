package server

import (
	"context"
	"personal-feed/pkg/model"
	"testing"
	"time"
)

func TestQ(t *testing.T) {
	server := &Server{}
	lastRunTime, _ := time.Parse(time.RFC3339, "2023-05-03T12:00:01Z")
	currentTime, _ := time.Parse(time.RFC3339, "2023-05-03T16:00:01Z")
	server.runIteration(context.Background(), nil, &model.Source{Schedule: "0 */6 * * *"}, &lastRunTime, currentTime)
}
