package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

type CleanupStats struct {
	ContainersRemoved int64
	ImagesRemoved     int64
	VolumesRemoved    int64
	NetworksRemoved   int64
	SpaceReclaimed    uint64
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Client{cli: cli}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) CleanContainers(ctx context.Context, all bool) (*CleanupStats, error) {
	stats := &CleanupStats{}

	if all {
		containers, err := c.cli.ContainerList(ctx, container.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list containers: %w", err)
		}

		for _, cont := range containers {
			timeout := 10 // 10 seconds timeout
			if err := c.cli.ContainerStop(ctx, cont.ID, container.StopOptions{Timeout: &timeout}); err != nil {
				fmt.Printf("Warning: failed to stop container %s: %v\n", cont.ID[:12], err)
			}
		}
	}

	pruneFilters := filters.NewArgs()

	report, err := c.cli.ContainersPrune(ctx, pruneFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to prune containers: %w", err)
	}

	stats.ContainersRemoved = int64(len(report.ContainersDeleted))
	stats.SpaceReclaimed = report.SpaceReclaimed

	return stats, nil
}

func (c *Client) ForceCleanContainers(ctx context.Context) (*CleanupStats, error) {
	stats := &CleanupStats{}

	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	for _, cont := range containers {
		if cont.State == "running" {
			timeout := 5
			if err := c.cli.ContainerStop(ctx, cont.ID, container.StopOptions{Timeout: &timeout}); err != nil {
				fmt.Printf("Warning: failed to stop container %s: %v\n", cont.ID[:12], err)
			}
		}
	}

	for _, cont := range containers {
		if err := c.cli.ContainerKill(ctx, cont.ID, "SIGKILL"); err != nil {
		}
	}

	report, err := c.cli.ContainersPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("failed to prune containers: %w", err)
	}

	stats.ContainersRemoved = int64(len(report.ContainersDeleted))
	stats.SpaceReclaimed = report.SpaceReclaimed

	return stats, nil
}

func (c *Client) CleanImages(ctx context.Context, dangling bool) (*CleanupStats, error) {
	stats := &CleanupStats{}

	pruneFilters := filters.NewArgs()
	if dangling {
		pruneFilters.Add("dangling", "true")
	}

	report, err := c.cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to prune images: %w", err)
	}

	stats.ImagesRemoved = int64(len(report.ImagesDeleted))
	stats.SpaceReclaimed += report.SpaceReclaimed

	return stats, nil
}

func (c *Client) CleanVolumes(ctx context.Context) (*CleanupStats, error) {
	stats := &CleanupStats{}

	report, err := c.cli.VolumesPrune(ctx, filters.Args{})
	if err != nil {
		return nil, fmt.Errorf("failed to prune volumes: %w", err)
	}

	stats.VolumesRemoved = int64(len(report.VolumesDeleted))
	stats.SpaceReclaimed += report.SpaceReclaimed

	return stats, nil
}

func (c *Client) CleanNetworks(ctx context.Context) (*CleanupStats, error) {
	stats := &CleanupStats{}

	report, err := c.cli.NetworksPrune(ctx, filters.Args{})
	if err != nil {
		return nil, fmt.Errorf("failed to prune networks: %w", err)
	}

	stats.NetworksRemoved = int64(len(report.NetworksDeleted))

	return stats, nil
}

func (c *Client) CleanAll(ctx context.Context) (*CleanupStats, error) {
	totalStats := &CleanupStats{}

	if containerStats, err := c.ForceCleanContainers(ctx); err != nil {
		fmt.Printf("Warning: container cleanup failed: %v\n", err)
	} else {
		totalStats.ContainersRemoved = containerStats.ContainersRemoved
		totalStats.SpaceReclaimed += containerStats.SpaceReclaimed
	}

	if imageStats, err := c.CleanImages(ctx, false); err != nil {
		fmt.Printf("Warning: image cleanup failed: %v\n", err)
	} else {
		totalStats.ImagesRemoved = imageStats.ImagesRemoved
		totalStats.SpaceReclaimed += imageStats.SpaceReclaimed
	}

	if volumeStats, err := c.CleanVolumes(ctx); err != nil {
		fmt.Printf("Warning: volume cleanup failed: %v\n", err)
	} else {
		totalStats.VolumesRemoved = volumeStats.VolumesRemoved
		totalStats.SpaceReclaimed += volumeStats.SpaceReclaimed
	}

	if networkStats, err := c.CleanNetworks(ctx); err != nil {
		fmt.Printf("Warning: network cleanup failed: %v\n", err)
	} else {
		totalStats.NetworksRemoved = networkStats.NetworksRemoved
	}

	return totalStats, nil
}

func (c *Client) SystemPrune(ctx context.Context, all bool) (*CleanupStats, error) {
	stats := &CleanupStats{}

	containerReport, err := c.cli.ContainersPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("failed to prune containers: %w", err)
	}
	stats.ContainersRemoved = int64(len(containerReport.ContainersDeleted))
	stats.SpaceReclaimed += containerReport.SpaceReclaimed

	imageFilters := filters.NewArgs()
	if !all {
		imageFilters.Add("dangling", "false")
	}
	imageReport, err := c.cli.ImagesPrune(ctx, imageFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to prune images: %w", err)
	}
	stats.ImagesRemoved = int64(len(imageReport.ImagesDeleted))
	stats.SpaceReclaimed += imageReport.SpaceReclaimed

	volumeReport, err := c.cli.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("failed to prune volumes: %w", err)
	}
	stats.VolumesRemoved = int64(len(volumeReport.VolumesDeleted))
	stats.SpaceReclaimed += volumeReport.SpaceReclaimed

	networkReport, err := c.cli.NetworksPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("failed to prune networks: %w", err)
	}
	stats.NetworksRemoved = int64(len(networkReport.NetworksDeleted))

	return stats, nil
}

func (c *Client) ListContainers(ctx context.Context) ([]types.Container, error) {
	return c.cli.ContainerList(ctx, container.ListOptions{All: true})
}

func (c *Client) ListImages(ctx context.Context) ([]image.Summary, error) {
	return c.cli.ImageList(ctx, image.ListOptions{All: true})
}

func (c *Client) ListVolumes(ctx context.Context) (volume.ListResponse, error) {
	return c.cli.VolumeList(ctx, volume.ListOptions{})
}

func (c *Client) GetSystemInfo(ctx context.Context) (system.Info, error) {
	return c.cli.Info(ctx)
}

func (c *Client) GetDiskUsage(ctx context.Context) (types.DiskUsage, error) {
	return c.cli.DiskUsage(ctx, types.DiskUsageOptions{})
}
