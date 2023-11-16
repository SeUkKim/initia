package keeper

import (
	goerrors "errors"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	customtypes "github.com/initia-labs/initia/x/gov/types"
)

// SubmitProposal create new proposal given a content
func (keeper Keeper) SubmitProposal(ctx sdk.Context, messages []sdk.Msg, metadata, title, summary string, proposer sdk.AccAddress) (v1.Proposal, error) {
	err := keeper.assertMetadataLength(metadata)
	if err != nil {
		return v1.Proposal{}, err
	}

	// assert summary is no longer than predefined max length of metadata
	err = keeper.assertMetadataLength(summary)
	if err != nil {
		return v1.Proposal{}, err
	}

	// assert title is no longer than predefined max length of metadata
	err = keeper.assertMetadataLength(title)
	if err != nil {
		return v1.Proposal{}, err
	}

	// Will hold a string slice of all Msg type URLs.
	msgs := []string{}

	// Loop through all messages and confirm that each has a handler and the gov module account
	// as the only signer
	for _, msg := range messages {
		msgs = append(msgs, sdk.MsgTypeURL(msg))

		// perform a basic validation of the message
		if err := msg.ValidateBasic(); err != nil {
			return v1.Proposal{}, errors.Wrap(types.ErrInvalidProposalMsg, err.Error())
		}

		signers := msg.GetSigners()
		if len(signers) != 1 {
			return v1.Proposal{}, types.ErrInvalidSigner
		}

		// assert that the governance module account is the only signer of the messages
		if !signers[0].Equals(keeper.GetGovernanceAccount(ctx).GetAddress()) {
			return v1.Proposal{}, errors.Wrapf(types.ErrInvalidSigner, signers[0].String())
		}

		// use the msg service router to see that there is a valid route for that message.
		handler := keeper.router.Handler(msg)
		if handler == nil {
			return v1.Proposal{}, errors.Wrap(types.ErrUnroutableProposalMsg, sdk.MsgTypeURL(msg))
		}

		// Only if it's a MsgExecLegacyContent do we try to execute the
		// proposal in a cached context.
		// For other Msgs, we do not verify the proposal messages any further.
		// They may fail upon execution.
		// ref: https://github.com/cosmos/cosmos-sdk/pull/10868#discussion_r784872842
		if msg, ok := msg.(*v1.MsgExecLegacyContent); ok {
			cacheCtx, _ := ctx.CacheContext()
			if _, err := handler(cacheCtx, msg); err != nil {
				if goerrors.Is(types.ErrNoProposalHandlerExists, err) {
					return v1.Proposal{}, err
				}
				return v1.Proposal{}, errors.Wrap(types.ErrInvalidProposalContent, err.Error())
			}
		}
	}

	proposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return v1.Proposal{}, err
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := keeper.GetParams(ctx).MaxDepositPeriod

	proposal, err := v1.NewProposal(messages, proposalID, submitTime, submitTime.Add(depositPeriod), metadata, title, summary, proposer)
	if err != nil {
		return v1.Proposal{}, err
	}

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposalID, *proposal.DepositEndTime)
	keeper.SetProposalID(ctx, proposalID+1)

	// called right after a proposal is submitted
	keeper.Hooks().AfterProposalSubmission(ctx, proposalID)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalMessages, strings.Join(msgs, ",")),
		),
	)

	return proposal, nil
}

// GetProposal get proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (v1.Proposal, bool) {
	store := ctx.KVStore(keeper.storeKey)

	bz := store.Get(types.ProposalKey(proposalID))
	if bz == nil {
		return v1.Proposal{}, false
	}

	var proposal v1.Proposal
	keeper.MustUnmarshalProposal(bz, &proposal)

	return proposal, true
}

// SetProposal set a proposal to store
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal v1.Proposal) {
	store := ctx.KVStore(keeper.storeKey)

	if proposal.Status == v1.StatusVotingPeriod {
		store.Set(types.VotingPeriodProposalKey(proposal.Id), []byte{1})
	} else {
		store.Delete(types.VotingPeriodProposalKey(proposal.Id))
	}

	bz := keeper.MustMarshalProposal(proposal)
	store.Set(types.ProposalKey(proposal.Id), bz)
}

// DeleteProposal deletes a proposal from store
func (keeper Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		panic(fmt.Sprintf("couldn't find proposal with id#%d", proposalID))
	}

	if proposal.DepositEndTime != nil {
		keeper.RemoveFromInactiveProposalQueue(ctx, proposalID, *proposal.DepositEndTime)
	}
	if proposal.VotingEndTime != nil {
		keeper.RemoveFromActiveProposalQueue(ctx, proposalID, *proposal.VotingEndTime)
		store.Delete(types.VotingPeriodProposalKey(proposalID))
	}

	store.Delete(types.ProposalKey(proposalID))
}

// IterateProposals iterates over the all the proposals and performs a callback function
// Panics when the iterator encounters a proposal which can't be unmarshaled.
func (keeper Keeper) IterateProposals(ctx sdk.Context, cb func(proposal v1.Proposal) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ProposalsKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal v1.Proposal
		keeper.MustUnmarshalProposal(iterator.Value(), &proposal)

		if cb(proposal) {
			break
		}
	}
}

// GetProposals returns all the proposals from store
func (keeper Keeper) GetProposals(ctx sdk.Context) (proposals v1.Proposals) {
	keeper.IterateProposals(ctx, func(proposal v1.Proposal) bool {
		proposals = append(proposals, &proposal)
		return false
	})
	return
}

// GetProposalsFiltered retrieves proposals filtered by a given set of params which
// include pagination parameters along with voter and depositor addresses and a
// proposal status. The voter address will filter proposals by whether or not
// that address has voted on proposals. The depositor address will filter proposals
// by whether or not that address has deposited to them. Finally, status will filter
// proposals by status.
//
// NOTE: If no filters are provided, all proposals will be returned in paginated
// form.
func (keeper Keeper) GetProposalsFiltered(ctx sdk.Context, params v1.QueryProposalsParams) v1.Proposals {
	proposals := keeper.GetProposals(ctx)
	filteredProposals := make([]*v1.Proposal, 0, len(proposals))

	for _, p := range proposals {
		matchVoter, matchDepositor, matchStatus := true, true, true

		// match status (if supplied/valid)
		if v1.ValidProposalStatus(params.ProposalStatus) {
			matchStatus = p.Status == params.ProposalStatus
		}

		// match voter address (if supplied)
		if len(params.Voter) > 0 {
			_, matchVoter = keeper.GetVote(ctx, p.Id, params.Voter)
		}

		// match depositor (if supplied)
		if len(params.Depositor) > 0 {
			_, matchDepositor = keeper.GetDeposit(ctx, p.Id, params.Depositor)
		}

		if matchVoter && matchDepositor && matchStatus {
			filteredProposals = append(filteredProposals, p)
		}
	}

	start, end := client.Paginate(len(filteredProposals), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredProposals = []*v1.Proposal{}
	} else {
		filteredProposals = filteredProposals[start:end]
	}

	return filteredProposals
}

// GetProposalID gets the highest proposal ID
func (keeper Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ProposalIDKey)
	if bz == nil {
		return 0, errors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = types.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

// SetProposalID sets the new proposal ID to the store
func (keeper Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.ProposalIDKey, types.GetProposalIDBytes(proposalID))
}

func (keeper Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal v1.Proposal) {
	startTime := ctx.BlockHeader().Time
	proposal.VotingStartTime = &startTime
	votingPeriod := keeper.GetParams(ctx).VotingPeriod
	endTime := proposal.VotingStartTime.Add(votingPeriod)
	proposal.VotingEndTime = &endTime
	proposal.Status = v1.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.Id, *proposal.DepositEndTime)
	keeper.InsertActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)

	// Emergency proposal registration
	if sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(keeper.GetParams(ctx).EmergencyMinDeposit) {
		keeper.InsertEmergencyProposalQueue(ctx, proposal.Id)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				customtypes.EventTypeEmergencyProposal,
				sdk.NewAttribute(customtypes.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
			),
		)
	}
}

func (keeper Keeper) MarshalProposal(proposal v1.Proposal) ([]byte, error) {
	bz, err := keeper.cdc.Marshal(&proposal)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (keeper Keeper) UnmarshalProposal(bz []byte, proposal *v1.Proposal) error {
	err := keeper.cdc.Unmarshal(bz, proposal)
	if err != nil {
		return err
	}
	return nil
}

func (keeper Keeper) MustMarshalProposal(proposal v1.Proposal) []byte {
	bz, err := keeper.MarshalProposal(proposal)
	if err != nil {
		panic(err)
	}
	return bz
}

func (keeper Keeper) MustUnmarshalProposal(bz []byte, proposal *v1.Proposal) {
	err := keeper.UnmarshalProposal(bz, proposal)
	if err != nil {
		panic(err)
	}
}
