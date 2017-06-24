package commands

import (
	"fmt"

	"github.com/docker-slim/docker-slim/master/config"
	"github.com/docker-slim/docker-slim/master/docker/dockerclient"
	"github.com/docker-slim/docker-slim/master/inspectors/image"
	"github.com/docker-slim/docker-slim/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
)

// OnInfo implements the 'info' docker-slim command
func OnInfo(statePath string, clientConfig *config.DockerClient, imageRef string) {

	fmt.Println("docker-slim: [info] image=", imageRef)

	client := dockerclient.New(clientConfig)

	imageInspector, err := image.NewInspector(client, imageRef)
	utils.FailOn(err)

	if imageInspector.NoImage() {
		fmt.Println("docker-slim: [info] target image not found -", imageRef)
		return
	}

	log.Info("docker-slim: inspecting 'fat' image metadata...")
	err = imageInspector.Inspect()
	utils.FailOn(err)

	_, artifactLocation := utils.PrepareSlimDirs(statePath, imageInspector.ImageInfo.ID)
	imageInspector.ArtifactLocation = artifactLocation

	log.Infof("docker-slim: [%v] 'fat' image size => %v (%v)",
		imageInspector.ImageInfo.ID,
		imageInspector.ImageInfo.VirtualSize,
		humanize.Bytes(uint64(imageInspector.ImageInfo.VirtualSize)))

	log.Info("docker-slim: processing 'fat' image info...")
	err = imageInspector.ProcessCollectedData()
	utils.FailOn(err)

	fmt.Println("docker-slim: [info] done.")
}
