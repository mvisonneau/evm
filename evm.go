package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jpillora/backoff"
	"github.com/shirou/gopsutil/disk"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time

func mount(c *cli.Context) error {
	start = time.Now()
	logConfig.configure()

	if user, err := user.Current(); err != nil {
		return exit(cli.NewExitError("Unable to determine current user", 1))
	} else if user.Username != "root" {
		return exit(cli.NewExitError("You have to run this function as root", 1))
	}

	if c.String("block-device-name") == "" ||
		 c.String("filesystem-type") == "" ||
		 c.String("mount-point") == "" ||
		 c.String("volume-name") == "" {
		cli.ShowSubcommandHelp(c)
		return cli.NewExitError("", 1)
	}

	mdsClient, err := getAWSMDSClient()
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	instanceID, err := getInstanceID(mdsClient)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	 az, err := getInstanceAZ(mdsClient)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	region := computeRegionFromAZ(az)
	ec2Client := getAWSEC2Client(region)

	volume, err := getVolume(ec2Client, c.String("volume-name"), az)
	if err != nil {
		return exit(cli.NewExitError(analyzeEC2APIErrors(err), 1))
	}

	if attached, err := isVolumeAttached(volume, instanceID); err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
  } else if attached {
		log.Infof("Volume is already attached to the instance")
	} else {
		log.Infof("Volume is available, attaching it to the instance")
		err = attachVolume(ec2Client, volume, instanceID, c.String("block-device-name"))
		if err != nil {
			return exit(cli.NewExitError(analyzeEC2APIErrors(err), 1))
		}
		log.Infof("Volume attached!")
	}

	if err := configureVolume(c.String("block-device-name"), c.String("filesystem-type"), c.String("mount-point")); err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	return exit(nil)
}

func getAWSMDSClient() (*ec2metadata.EC2Metadata, error) {
	log.Debug("Starting AWS MDS API session")
	client := ec2metadata.New(session.New())

	if !client.Available() {
		return client, errors.New("Unable to access the metadata service, are you running this binary from an AWS EC2 instance?")
	}

	return client, nil
}

func getAWSEC2Client(region string) (client *ec2.EC2) {
	log.Debug("Starting AWS EC2 API session")
	client = ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	return
}

func getInstanceAZ(c *ec2metadata.EC2Metadata) (az string, err error) {
	log.Debug("Fetching current AZ from MDS API")
	az, err = c.GetMetadata("placement/availability-zone")
	log.Infof("Found AZ: '%s'", az)
	return
}

func computeRegionFromAZ(az string) string {
	log.Infof("Computed region : '%s'", az[:len(az)-1])
	return az[:len(az)-1]
}

func getInstanceID(c *ec2metadata.EC2Metadata) (id string, err error) {
	log.Debug("Fetching current instance-id from MDS API")
	id, err = c.GetMetadata("instance-id")
	log.Infof("Found instance-id : '%s'", id)
	return
}

func analyzeEC2APIErrors(err error) string {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr.Error()
		}
		return err.Error()
	}
	return ""
}

func getVolume(c *ec2.EC2, volume string, az string) (*ec2.Volume, error) {
	log.Debugf("Looking up volume '%s' in '%s'", volume, az)

	volumes, err := c.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(volume)},
			},
			{
				Name:   aws.String("availability-zone"),
				Values: []*string{aws.String(az)},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(volumes.Volumes) != 1 {
		return nil, fmt.Errorf("Unexpected amount of volumes retrieved : '%d',  expected 1", len(volumes.Volumes))
	}

	log.Printf("Found volume '%s' in instance AZ (%s) with ID : '%s'!", volume, az, volumes.Volumes[0].VolumeId)

	return volumes.Volumes[0], nil
}

func isVolumeAttached(volume *ec2.Volume, instanceID string) (bool, error) {
	log.Debugf("Checking if the volume is attached")

	if len(volume.Attachments) > 1 {
		return false, fmt.Errorf("Unexpected amount of attachments : '%d',  expected 0 or 1", len(volume.Attachments))
	}

	if len(volume.Attachments) > 0 {
		if *volume.Attachments[0].InstanceId != instanceID {
			return true, fmt.Errorf("Volume is attached onto another instance : '%s'", *volume.Attachments[0].InstanceId)
		} else {
			return true, nil
		}
	}

	return false, nil
}

func attachVolume(c *ec2.EC2, volume *ec2.Volume, instanceID string, blockDeviceName string) error {
	_, err := c.AttachVolume(&ec2.AttachVolumeInput{
		Device:     aws.String(blockDeviceName),
		InstanceId: aws.String(instanceID),
		VolumeId:   volume.VolumeId,
	})

	if err != nil {
		return err
	}

	b := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    100 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	for {
		volumes, err := c.DescribeVolumes(&ec2.DescribeVolumesInput{
			VolumeIds: []*string{volume.VolumeId},
		})
		if err != nil {
			return err
		}

		if len(volumes.Volumes) == 0 {
			continue
		}

		if len(volumes.Volumes[0].Attachments) == 0 {
			continue
		}

		if *volumes.Volumes[0].Attachments[0].State == ec2.VolumeAttachmentStateAttached {
			break
		}

		log.Debugf("Waiting for attachment to complete. '%s'", *volumes.Volumes[0].Attachments[0].State)
		b.Duration()
	}

	return nil
}

func configureVolume(blockDeviceName, fileSystemType, mountPoint string) error {
	log.Printf("Checking for existing filesystem on block device '%s'", blockDeviceName)

	disks, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	deviceFound := false
	for _, disk := range disks {
    if disk.Device == blockDeviceName {
			deviceFound = true

			if disk.Fstype == fileSystemType {
				log.Infof("Volume already formatted to '%s'", fileSystemType)
			} else {
				log.Infof("Formatting volume.. expected '%s' - current '%s'", fileSystemType, disk.Fstype)
				cmd := exec.Command("/usr/sbin/mkfs."+fileSystemType, blockDeviceName)
				cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
				log.Infof("Volume formatted to '%s'", fileSystemType)
			}

			if disk.Mountpoint == mountPoint {
				log.Infof("Volume already mounted to '%s'", mountPoint)
			} else {
				log.Infof("Mounting '%s' to '%s'", blockDeviceName, mountPoint)
				if err := syscall.Mount(blockDeviceName, mountPoint, fileSystemType, 2, ""); err != nil {
					return err
				}
				log.Infof("Volume mounted to '%s'", mountPoint)
			}

      break
    }
	}

	if ! deviceFound {
		return fmt.Errorf("Block device '%s' not found on the system", blockDeviceName)
	}

	return nil
}

func exit(err error) error {
	log.Debugf("Executed in %s, exiting..", time.Since(start))
	return err
}
