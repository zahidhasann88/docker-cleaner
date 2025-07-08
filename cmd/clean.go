package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zahidhasann88/docker-cleaner/pkg/docker"
)

var (
	cleanAll        bool
	cleanContainers bool
	cleanImages     bool
	cleanVolumes    bool
	cleanNetworks   bool
	force           bool
	dangling        bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean Docker resources",
	Long: `Clean Docker resources including containers, images, volumes, and networks.
By default, only removes stopped containers and dangling images.`,
	RunE: runClean,
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean all resources (containers, images, volumes, networks)")
	cleanCmd.Flags().BoolVarP(&cleanContainers, "containers", "c", false, "Clean containers")
	cleanCmd.Flags().BoolVarP(&cleanImages, "images", "i", false, "Clean images")
	cleanCmd.Flags().BoolVarP(&cleanVolumes, "volumes", "v", false, "Clean volumes")
	cleanCmd.Flags().BoolVarP(&cleanNetworks, "networks", "n", false, "Clean networks")
	cleanCmd.Flags().BoolVarP(&force, "force", "f", false, "Force removal without confirmation")
	cleanCmd.Flags().BoolVar(&dangling, "dangling", true, "Only remove dangling images (default: true)")
}

func runClean(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer client.Close()

	if !cleanAll && !cleanContainers && !cleanImages && !cleanVolumes && !cleanNetworks {
		cleanContainers = true
		cleanImages = true
	}

	if cleanAll {
		cleanContainers = true
		cleanImages = true
		cleanVolumes = true
		cleanNetworks = true
	}

	if !force {
		fmt.Println("This will remove Docker resources. Continue? (y/N)")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	totalStats := &docker.CleanupStats{}

	if cleanContainers {
		fmt.Println("🧹 Cleaning containers...")
		stats, err := client.CleanContainers(ctx, cleanAll)
		if err != nil {
			return fmt.Errorf("failed to clean containers: %w", err)
		}
		totalStats.ContainersRemoved += stats.ContainersRemoved
		fmt.Printf("   ✓ Removed %d containers\n", stats.ContainersRemoved)
	}

	if cleanImages {
		fmt.Println("🖼️  Cleaning images...")
		stats, err := client.CleanImages(ctx, dangling)
		if err != nil {
			return fmt.Errorf("failed to clean images: %w", err)
		}
		totalStats.ImagesRemoved += stats.ImagesRemoved
		totalStats.SpaceReclaimed += stats.SpaceReclaimed
		fmt.Printf("   ✓ Removed %d images\n", stats.ImagesRemoved)
	}

	if cleanVolumes {
		fmt.Println("💾 Cleaning volumes...")
		stats, err := client.CleanVolumes(ctx)
		if err != nil {
			return fmt.Errorf("failed to clean volumes: %w", err)
		}
		totalStats.VolumesRemoved += stats.VolumesRemoved
		totalStats.SpaceReclaimed += stats.SpaceReclaimed
		fmt.Printf("   ✓ Removed %d volumes\n", stats.VolumesRemoved)
	}

	if cleanNetworks {
		fmt.Println("🌐 Cleaning networks...")
		stats, err := client.CleanNetworks(ctx)
		if err != nil {
			return fmt.Errorf("failed to clean networks: %w", err)
		}
		totalStats.NetworksRemoved += stats.NetworksRemoved
		fmt.Printf("   ✓ Removed %d networks\n", stats.NetworksRemoved)
	}

	fmt.Println("\n📊 Cleanup Summary:")
	fmt.Printf("   Containers: %d\n", totalStats.ContainersRemoved)
	fmt.Printf("   Images: %d\n", totalStats.ImagesRemoved)
	fmt.Printf("   Volumes: %d\n", totalStats.VolumesRemoved)
	fmt.Printf("   Networks: %d\n", totalStats.NetworksRemoved)
	fmt.Printf("   Space reclaimed: %.2f MB\n", float64(totalStats.SpaceReclaimed)/(1024*1024))

	return nil
}
