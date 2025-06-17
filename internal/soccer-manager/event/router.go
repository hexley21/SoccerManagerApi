package event

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
)

func RegisterEventHandlers(eventEmitter evbus.BusSubscriber, c *delivery.Components) {
	onSignup := newUserSignUpHandler(
		c.Services.TeamService,
		c.Services.PlayerService,
		c.Cfg.Events.UserSignUp,
		c.Logger,
	)

	eventEmitter.Subscribe(domain.EventONSIGNUP, onSignup.Handle)
}
