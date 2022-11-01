package interaction

import (
	"sync"
	"time"

	commandpb "go.temporal.io/api/command/v1"
	enumspb "go.temporal.io/api/enums/v1"
	failurepb "go.temporal.io/api/failure/v1"
	interactionpb "go.temporal.io/api/interaction/v1"
)

type (
	Registry interface {
		Add(in *interactionpb.Input, requestID string, requestTime time.Time, lastEventID int64) (*Interaction, RemoveFunc)
		PendingInvocations() []*interactionpb.Invocation

		HasPending(filterCommands ...*commandpb.Command) bool
		Pending(id string) *Interaction
		Accepted(id string) *Interaction

		// Clear() // TODO: not sure if this is needed
	}

	RemoveFunc func()

	RegistryImpl struct {
		sync.RWMutex
		interactions map[string]*Interaction
	}
)

func NewRegistry() *RegistryImpl {
	return &RegistryImpl{
		interactions: make(map[string]*Interaction),
	}
}

func (r *RegistryImpl) Pending(id string) *Interaction {
	r.RLock()
	defer r.RUnlock()
	if inter, ok := r.interactions[id]; ok && inter.State() == StatePending {
		return inter
	}
	return nil
}

func (r *RegistryImpl) Accepted(id string) *Interaction {
	r.RLock()
	defer r.RUnlock()
	if inter, ok := r.interactions[id]; ok && inter.State() == StateAccepted {
		return inter
	}
	return nil
}

func (r *RegistryImpl) Add(in *interactionpb.Input, requestID string, requestTime time.Time, lastEventID int64) (*Interaction, RemoveFunc) {
	r.Lock()
	defer r.Unlock()
	inter := newInteraction(in, requestID, requestTime, lastEventID)
	r.interactions[inter.meta.Id] = inter
	return inter, func() { r.remove(inter.meta.Id) }
}

func (r *RegistryImpl) HasPending(filterCommands ...*commandpb.Command) bool {
	// Filter out interactions which will be accepted or rejected by current commands.
	// This interaction has Pending state in the registry but in fact, they are already accepted or rejected
	// and should not be counted as pending.
	notPendingInteractions := make(map[string]struct{})
	for _, command := range filterCommands {
		if command.GetCommandType() == enumspb.COMMAND_TYPE_REJECT_WORKFLOW_UPDATE {
			notPendingInteractions[command.GetRejectWorkflowUpdateCommandAttributes().GetMeta().GetId()] = struct{}{}
		}
		if command.GetCommandType() == enumspb.COMMAND_TYPE_ACCEPT_WORKFLOW_UPDATE {
			notPendingInteractions[command.GetAcceptWorkflowUpdateCommandAttributes().GetMeta().GetId()] = struct{}{}
		}
	}

	r.RLock()
	defer r.RUnlock()

	for _, inter := range r.interactions {
		if _, notPending := notPendingInteractions[inter.meta.GetId()]; !notPending && inter.state == StatePending {
			return true
		}
	}
	return false
}

func (r *RegistryImpl) PendingInvocations() []*interactionpb.Invocation {
	r.RLock()
	defer r.RUnlock()
	invocations := make([]*interactionpb.Invocation, 0, len(r.interactions))
	for _, interaction := range r.interactions {
		if interaction.state == StatePending {
			invocations = append(invocations, &interactionpb.Invocation{
				Meta:  interaction.meta,
				Input: interaction.in,
			})
		}
	}
	return invocations
}

func (r *RegistryImpl) Clear() {
	r.Lock()
	defer r.Unlock()
	for _, inter := range r.interactions {
		inter.SendReject(r.clearFailure())
	}
	r.interactions = make(map[string]*Interaction)
}

func (r *RegistryImpl) clearFailure() *failurepb.Failure {
	return &failurepb.Failure{
		Message: "update cleared, please retry",
		FailureInfo: &failurepb.Failure_ServerFailureInfo{
			ServerFailureInfo: &failurepb.ServerFailureInfo{
				NonRetryable: true,
			},
		},
	}
}

func (r *RegistryImpl) remove(id string) {
	r.Lock()
	defer r.Unlock()
	delete(r.interactions, id)
}
