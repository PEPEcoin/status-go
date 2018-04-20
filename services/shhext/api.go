package shhext

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/status-im/status-go/services/shhext/dedup"
)

// NewPublicAPI returns instance of the public API.
func NewPublicAPI(w *whisper.Whisper, tracker *tracker, deduplicator *dedup.Deduplicator) *PublicAPI {
	return &PublicAPI{
		w:         w,
		publicAPI: whisper.NewPublicWhisperAPI(w),
		tracker:   tracker,
		d:         deduplicator,
	}
}

// PublicAPI extends whisper public API.
type PublicAPI struct {
	w         *whisper.Whisper
	publicAPI *whisper.PublicWhisperAPI
	tracker   *tracker
	d         *dedup.Deduplicator
}

// Post shamelessly copied from whisper codebase with slight modifications.
func (api *PublicAPI) Post(ctx context.Context, req whisper.NewMessage) (hash hexutil.Bytes, err error) {
	hash, err = api.publicAPI.Post(ctx, req)
	if err == nil {
		var envHash common.Hash
		copy(envHash[:], hash[:]) // slice can't be used as key
		api.tracker.Add(envHash)
	}
	return hash, err
}

// GetNewFilterMessages is a prototype method with deduplication
func (api *PublicAPI) GetNewFilterMessages(filterID string) ([]*whisper.Message, error) {
	msgs, err := api.publicAPI.GetFilterMessages(filterID)
	if err != nil {
		return nil, err
	}
	return api.d.Deduplicate(msgs), err
}
