package configuration

import (
	"fmt"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/cbsinteractive/transcode-orchestrator/db"
	"github.com/cbsinteractive/transcode-orchestrator/provider/bitmovin/internal/configuration/codec"
	"github.com/pkg/errors"
)

// H264 is a configuration service for content in h.264
type H264 struct {
	api  *bitmovin.BitmovinApi
	repo db.PresetSummaryRepository
}

// NewH264 returns a service for managing h.264 configurations
func NewH264(api *bitmovin.BitmovinApi, repo db.PresetSummaryRepository) *H264 {
	return &H264{api: api, repo: repo}
}

// Create will create a new H264 configuration based on a preset
func (c *H264) Create(preset db.Preset) (string, error) {
	vidCfgID, err := codec.NewH264(c.api, preset)
	if err != nil {
		return "", err
	}

	err = c.repo.CreatePresetSummary(&db.PresetSummary{
		Name:          preset.Name,
		Container:     preset.Container,
		VideoCodec:    string(model.CodecConfigType_H264),
		VideoConfigID: vidCfgID,
	})
	if err != nil {
		return "", err
	}

	return preset.Name, nil
}

// Get retrieves a stored db.PresetSummary by its name
func (c *H264) Get(presetName string) (db.PresetSummary, error) {
	return c.repo.GetPresetSummary(presetName)
}

// Delete removes the video configuration
func (c *H264) Delete(presetName string) error {
	summary, err := c.Get(presetName)
	if err != nil {
		return err
	}

	_, err = c.api.Encoding.Configurations.Video.H264.Delete(summary.VideoConfigID)
	if err != nil {
		return errors.Wrap(err, "removing the video config")
	}

	err = c.repo.DeletePresetSummary(presetName)
	if err != nil {
		return fmt.Errorf("deleting preset summary: %w", err)
	}

	return nil
}
