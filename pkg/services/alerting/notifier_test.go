package alerting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/services/rendering"
)

var (
	alertingNotificationTimeout = 100 * time.Millisecond
)

type timeoutRenderService struct{}

func (frs timeoutRenderService) Render(ctx context.Context, opts rendering.Opts) (*rendering.RenderResult, error) {
	time.Sleep(alertingNotificationTimeout) // Make sure to consume all the time while rendering
	return nil, fmt.Errorf("can't render: %v", ctx.Err())
}

func TestImgRenderTimeout(t *testing.T) {

	// sqlstore := sqlstore.InitTestDB(t)

	// Create the notification service with a fake render service
	npt := newNotificationService(timeoutRenderService{})

	// Create the global context
	testContext, testContextCancel := context.WithTimeout(context.Background(), alertingNotificationTimeout)
	defer testContextCancel()

	fmt.Println(npt.SendIfNeeded(&EvalContext{
		Firing:    true,
		IsTestRun: true,
		Rule: &Rule{
			OrgID:         1,
			Notifications: []string{},
		},
		Ctx: testContext,
	}))
}
