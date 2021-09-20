// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package builtin

import (
	"fmt"
	"io"

	address "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/go-state-types/abi"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf

var lengthBufMinerAddrs = []byte{131}

func (t *MinerAddrs) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufMinerAddrs); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Owner (address.Address) (struct)
	if err := t.Owner.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Worker (address.Address) (struct)
	if err := t.Worker.MarshalCBOR(w); err != nil {
		return err
	}

	// t.ControlAddrs ([]address.Address) (slice)
	if len(t.ControlAddrs) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.ControlAddrs was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.ControlAddrs))); err != nil {
		return err
	}
	for _, v := range t.ControlAddrs {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}
	return nil
}

func (t *MinerAddrs) UnmarshalCBOR(r io.Reader) error {
	*t = MinerAddrs{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Owner (address.Address) (struct)

	{

		if err := t.Owner.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Owner: %w", err)
		}

	}
	// t.Worker (address.Address) (struct)

	{

		if err := t.Worker.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Worker: %w", err)
		}

	}
	// t.ControlAddrs ([]address.Address) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.ControlAddrs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.ControlAddrs = make([]address.Address, extra)
	}

	for i := 0; i < int(extra); i++ {

		var v address.Address
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}

		t.ControlAddrs[i] = v
	}

	return nil
}

var lengthBufConfirmSectorProofsParams = []byte{133}

func (t *ConfirmSectorProofsParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufConfirmSectorProofsParams); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Sectors ([]abi.SectorNumber) (slice)
	if len(t.Sectors) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Sectors was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Sectors))); err != nil {
		return err
	}
	for _, v := range t.Sectors {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}

	// t.PrecomputeRewardPowerStats (bool) (bool)
	if err := cbg.WriteBool(w, t.PrecomputeRewardPowerStats); err != nil {
		return err
	}

	// t.RewardStatsThisEpochRewardSmoothed (smoothing.FilterEstimate) (struct)
	if err := t.RewardStatsThisEpochRewardSmoothed.MarshalCBOR(w); err != nil {
		return err
	}

	// t.RewardStatsThisEpochBaselinePower (big.Int) (struct)
	if err := t.RewardStatsThisEpochBaselinePower.MarshalCBOR(w); err != nil {
		return err
	}

	// t.PwrTotalQualityAdjPowerSmoothed (smoothing.FilterEstimate) (struct)
	if err := t.PwrTotalQualityAdjPowerSmoothed.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *ConfirmSectorProofsParams) UnmarshalCBOR(r io.Reader) error {
	*t = ConfirmSectorProofsParams{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 5 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Sectors ([]abi.SectorNumber) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Sectors: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Sectors = make([]abi.SectorNumber, extra)
	}

	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Sectors slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Sectors was not a uint, instead got %d", maj)
		}

		t.Sectors[i] = abi.SectorNumber(val)
	}

	// t.PrecomputeRewardPowerStats (bool) (bool)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.PrecomputeRewardPowerStats = false
	case 21:
		t.PrecomputeRewardPowerStats = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.RewardStatsThisEpochRewardSmoothed (smoothing.FilterEstimate) (struct)

	{

		if err := t.RewardStatsThisEpochRewardSmoothed.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.RewardStatsThisEpochRewardSmoothed: %w", err)
		}

	}
	// t.RewardStatsThisEpochBaselinePower (big.Int) (struct)

	{

		if err := t.RewardStatsThisEpochBaselinePower.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.RewardStatsThisEpochBaselinePower: %w", err)
		}

	}
	// t.PwrTotalQualityAdjPowerSmoothed (smoothing.FilterEstimate) (struct)

	{

		if err := t.PwrTotalQualityAdjPowerSmoothed.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.PwrTotalQualityAdjPowerSmoothed: %w", err)
		}

	}
	return nil
}
