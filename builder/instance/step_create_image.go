package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/errors"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var _ multistep.Step = &StepCreateImage{}

// StepCreateImage represents a Packer build step that creates the image.
type StepCreateImage struct{}

// Cleanup implements multistep.Step.
func (s *StepCreateImage) Cleanup(state multistep.StateBag) {}

// Run implements multistep.Step.
func (s *StepCreateImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)
	zsconf := state.Get("zsconf").(*client.ZSConfig)

	client := client.NewZSClient(zsconf)
	if _, err := client.Login(ctx); err != nil {
		if err != errors.ErrNotSupported {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	if _, ok := state.GetOk("instance_root_volume_uuid"); !ok {
		err := fmt.Errorf("unable to get instance_root_volume_uuid state")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Creating image...")
	image, err := client.CreateRootVolumeTemplateFromRootVolume(
		state.Get("instance_root_volume_uuid").(string),
		param.CreateRootVolumeTemplateFromRootVolumeParam{
			Params: param.CreateRootVolumeTemplateFromRootVolumeParamDetail{
				Name:               config.ImageName,
				Description:        &config.ImageDescription,
				BackupStorageUuids: []string{config.ImageBackupStorage},
			},
		})

	if err != nil {
		err := fmt.Errorf("error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Image has been created!")
	state.Put("image_uuid", image.UUID)
	state.Put("image_name", image.Name)
	return multistep.ActionContinue
}
