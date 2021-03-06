package container

import (
	"path"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/bitmovin/bitmovin-api-sdk-go/query"
	"github.com/cbsinteractive/transcode-orchestrator/provider"
	"github.com/cbsinteractive/transcode-orchestrator/provider/bitmovin/internal/storage"
	"github.com/pkg/errors"
)

// ProgressiveWebMAssembler is an assembler that creates ProgressiveWebM outputs based on a cfg
type ProgressiveWebMAssembler struct {
	api *bitmovin.BitmovinApi
}

// NewProgressiveWebMAssembler creates and returns an ProgressiveWebMAssembler
func NewProgressiveWebMAssembler(api *bitmovin.BitmovinApi) *ProgressiveWebMAssembler {
	return &ProgressiveWebMAssembler{api: api}
}

// Assemble creates ProgressiveWebM outputs
func (a *ProgressiveWebMAssembler) Assemble(cfg AssemblerCfg) error {
	_, err := a.api.Encoding.Encodings.Muxings.ProgressiveWebm.Create(cfg.EncID, model.ProgressiveWebmMuxing{
		Filename:             path.Base(cfg.OutputFilename),
		Streams:              streamsFrom(cfg),
		StreamConditionsMode: model.StreamConditionsMode_DROP_STREAM,
		Outputs: []model.EncodingOutput{
			storage.EncodingOutputFrom(cfg.OutputID, path.Dir(path.Join(cfg.DestPath, cfg.OutputFilename))),
		},
	})
	if err != nil {
		return errors.Wrap(err, "creating progressive webm muxing")
	}

	return nil
}

// ProgressiveWebMStatusEnricher is responsible for adding ProgressiveWebM output info to a job status
type ProgressiveWebMStatusEnricher struct {
	api *bitmovin.BitmovinApi
}

// NewProgressiveWebMStatusEnricher creates and returns an ProgressiveWebMStatusEnricher
func NewProgressiveWebMStatusEnricher(api *bitmovin.BitmovinApi) *ProgressiveWebMStatusEnricher {
	return &ProgressiveWebMStatusEnricher{api: api}
}

// Enrich populates information about ProgressiveWebM outputs if they exist
func (e *ProgressiveWebMStatusEnricher) Enrich(s provider.JobStatus) (provider.JobStatus, error) {
	var totalCount int64 = 1
	var muxings []model.ProgressiveWebmMuxing
	for int64(len(muxings)) < totalCount {
		resp, err := e.api.Encoding.Encodings.Muxings.ProgressiveWebm.List(s.ProviderJobID, func(params *query.ProgressiveWebmMuxingListQueryParams) {
			params.Offset = int32(len(muxings))
			params.Limit = 100
		})
		if err != nil {
			return s, errors.Wrap(err, "retrieving progressive webm muxings from the Bitmovin API")
		}

		totalCount = int64Value(resp.TotalCount)
		muxings = append(muxings, resp.Items...)
	}

	for _, muxing := range muxings {
		info, err := e.api.Encoding.Encodings.Muxings.ProgressiveWebm.Information.Get(s.ProviderJobID, muxing.Id)
		if err != nil {
			return s, errors.Wrapf(err, "retrieving muxing information with ID %q", muxing.Id)
		}

		var (
			height, width int64
			videoCodec    string
		)
		if len(info.VideoTracks) > 0 {
			track := info.VideoTracks[0]
			height, width = dimensionToInt64(track.FrameHeight), dimensionToInt64(track.FrameWidth)
			videoCodec = track.Codec
		}

		s.Output.Files = append(s.Output.Files, provider.OutputFile{
			Path:       s.Output.Destination + muxing.Filename,
			Container:  info.ContainerFormat,
			FileSize:   int64Value(info.FileSize),
			VideoCodec: videoCodec,
			Width:      width,
			Height:     height,
		})
	}

	return s, nil
}
