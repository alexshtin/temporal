package interaction

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	enumspb "go.temporal.io/api/enums/v1"
	failurepb "go.temporal.io/api/failure/v1"
	interactionpb "go.temporal.io/api/interaction/v1"

	"go.temporal.io/server/common/future"
)

type (
	Interaction struct {
		state       State
		requestTime time.Time
		meta        *interactionpb.Meta
		in          *interactionpb.Input
		out         *future.FutureImpl[*interactionpb.Output]
	}

	State int32
)

const (
	StatePending State = iota
	StateAccepted
	StateRejected
	StateCompleted
)

func newInteraction(in *interactionpb.Input, requestID string, requestTime time.Time, lastEventID int64) *Interaction {
	return &Interaction{
		state:       StatePending,
		in:          in,
		out:         future.NewFuture[*interactionpb.Output](),
		requestTime: requestTime,
		meta: &interactionpb.Meta{
			Id:              uuid.New(),
			InteractionType: enumspb.INTERACTION_TYPE_WORKFLOW_UPDATE,
			EventId:         lastEventID,
			RequestId:       requestID,
		},
	}
}

// func (u *Interaction) Meta() *interactionpb.Meta {
// 	return u.meta
// }
//
// func (u *Interaction) In() *interactionpb.Input {
// 	return u.in
// }
//
// func (u *Interaction) Out() *future.FutureImpl[*interactionpb.Output] {
// 	return u.out
// }
//
// func (u *Interaction) RequestTime() *time.Time {
// 	return &u.requestTime
// }

func (u *Interaction) State() State {
	return u.state
}

func (u *Interaction) Accept() {
	u.state = StateAccepted
}

func (u *Interaction) SendComplete(out *interactionpb.Output) {
	u.state = StateCompleted
	u.out.Set(out, nil)
}

func (u *Interaction) SendReject(failure *failurepb.Failure) {
	u.state = StateRejected
	u.out.Set(&interactionpb.Output{
		Result: &interactionpb.Output_Failure{
			Failure: failure,
		},
	}, nil)
}

func (u *Interaction) WaitResult(ctx context.Context) (*interactionpb.Output, error) {
	return u.out.Get(ctx)
}
