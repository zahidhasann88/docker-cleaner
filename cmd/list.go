package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zahidhasann88/docker-cleaner/pkg/docker"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Docker resources",
	Long:  `List Docker resources including containers, images, volumes, and networks.`,
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer client.Close()

	fmt.Println("📋 Docker Resources Overview")
	fmt.Println(strings.Repeat("=", 50))

	containers, err := client.ListContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	fmt.Printf("\n🐳 Containers (%d total):\n", len(containers))
	if len(containers) == 0 {
		fmt.Println("   No containers found")
	} else {
		fmt.Printf("   %-12s %-20s %-15s %-10s\n", "ID", "Image", "Status", "Names")
		fmt.Println("   " + strings.Repeat("-", 60))
		for _, container := range containers {
			names := strings.Join(container.Names, ", ")
			if len(names) > 20 {
				names = names[:17] + "..."
			}
			fmt.Printf("   %-12s %-20s %-15s %-10s\n",
				container.ID[:12],
				container.Image,
				container.Status,
				names)
		}
	}

	images, err := client.ListImages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	fmt.Printf("\n🖼️  Images (%d total):\n", len(images))
	if len(images) == 0 {
		fmt.Println("   No images found")
	} else {
		fmt.Printf("   %-12s %-30s %-10s %-15s\n", "ID", "Repository", "Tag", "Size")
		fmt.Println("   " + strings.Repeat("-", 70))
		for _, image := range images {
			repo := "<none>"
			tag := "<none>"
			if len(image.RepoTags) > 0 {
				parts := strings.Split(image.RepoTags[0], ":")
				repo = parts[0]
				if len(parts) > 1 {
					tag = parts[1]
				}
			}
			if len(repo) > 30 {
				repo = repo[:27] + "..."
			}
			fmt.Printf("   %-12s %-30s %-10s %-15s\n",
				image.ID[:12],
				repo,
				tag,
				formatSize(image.Size))
		}
	}

	volumes, err := client.ListVolumes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list volumes: %w", err)
	}

	fmt.Printf("\n💾 Volumes (%d total):\n", len(volumes.Volumes))
	if len(volumes.Volumes) == 0 {
		fmt.Println("   No volumes found")
	} else {
		fmt.Printf("   %-20s %-15s %-30s\n", "Name", "Driver", "Mountpoint")
		fmt.Println("   " + strings.Repeat("-", 70))
		for _, volume := range volumes.Volumes {
			name := volume.Name
			if len(name) > 20 {
				name = name[:17] + "..."
			}
			mountpoint := volume.Mountpoint
			if len(mountpoint) > 30 {
				mountpoint = mountpoint[:27] + "..."
			}
			fmt.Printf("   %-20s %-15s %-30s\n", name, volume.Driver, mountpoint)
		}
	}

	return nil
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
