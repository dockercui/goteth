package custom_spec

import (
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
)

type Phase0Spec struct {
	BState     spec.VersionedBeaconState
	Committees map[string]bitfield.Bitlist
}

func NewPhase0Spec(bstate *spec.VersionedBeaconState) Phase0Spec {
	phase0Obj := Phase0Spec{
		BState:     *bstate,
		Committees: make(map[string]bitfield.Bitlist),
	}
	phase0Obj.CalculatePreviousEpochAttestations()

	return phase0Obj

}

func (p Phase0Spec) ObtainCurrentSlot() uint64 {
	return p.BState.Phase0.Slot
}

func (p Phase0Spec) ObtainCurrentEpoch() uint64 {
	return uint64(p.ObtainCurrentSlot() / 32)
}

func (p *Phase0Spec) CalculatePreviousEpochAttestations() {

	totalAttPreviousEpoch := 0
	totalAttestingVals := 0

	previousAttestatons := p.BState.Phase0.PreviousEpochAttestations
	// currentAttestations := bstate.Phase0.CurrentEpochAttestations
	doubleVotes := 0
	vals := p.BState.Phase0.Validators

	// TODO: check validator active, slashed or exiting
	for _, item := range vals {
		if item.ActivationEligibilityEpoch < phase0.Epoch(p.ObtainCurrentEpoch()) {
			totalAttestingVals += 1
		}
	}

	for _, item := range previousAttestatons {
		slot := item.Data.Slot // Block that is being attested, not included
		committeeIndex := item.Data.Index
		mapKey := strconv.Itoa(int(slot)) + "_" + strconv.Itoa(int(committeeIndex))

		resultBits := bitfield.NewBitlist(0)

		if val, ok := p.Committees[mapKey]; ok {
			// the committeeIndex for the given slot already had an aggregation
			// TODO: check error
			allZero, err := val.And(item.AggregationBits) // if the same position of the aggregations has a 1 in the same position, there was a double vote
			if err != nil {
				fmt.Println(err)
			}
			resultBitstmp, err := val.Or(item.AggregationBits) // to join all aggregation we do Or operation

			if allZero.Count() > 0 {
				// there was a double vote
				doubleVotes += int(allZero.Count())
			}
			if err == nil {
				resultBits = resultBitstmp
			} else {
				fmt.Println(err)
			}

		} else {
			// we had not received any aggregation for this committeeIndex at the given slot
			resultBits = item.AggregationBits
		}
		p.Committees[mapKey] = resultBits
		attPreviousEpoch := int(item.AggregationBits.Count())
		totalAttPreviousEpoch += attPreviousEpoch // we are counting bits set to 1 aggregation by aggregation
		// if we do the And at the same committee we can catch the double votes
		// doing that the number of votes is less than the number of validators

	}

}

func (p Phase0Spec) ObtainPreviousEpochAttestations() uint64 {

	numOf1Bits := 0 // it should be equal to the number of validators that attested

	for _, val := range p.Committees {

		numOf1Bits += int(val.Count())
	}

	return uint64(numOf1Bits)
}

func (p Phase0Spec) ObtainPreviousEpochValNum() uint64 {

	numOfBits := 0 // it should be equal to the number of validators

	for _, val := range p.Committees {

		numOfBits += int(val.Len())
	}

	return uint64(numOfBits)
}

func (p Phase0Spec) ObtainBalance(valIdx uint64) (uint64, error) {
	if uint64(len(p.BState.Phase0.Balances)) < valIdx {
		err := fmt.Errorf("phase0 - validator index %d wasn't activated in slot %d", valIdx, p.BState.Phase0.Slot)
		return 0, err
	}
	balance := p.BState.Phase0.Balances[valIdx]

	return balance, nil
}
