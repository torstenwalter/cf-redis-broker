package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pivotal-cf/cf-redis-broker/backup/s3bucket"
	"github.com/pivotal-cf/cf-redis-broker/brokerconfig"
	"github.com/pivotal-cf/cf-redis-broker/redis"
	"github.com/pivotal-cf/cf-redis-broker/redis/client"
	"github.com/pivotal-cf/cf-redis-broker/redisconf"
	"github.com/pivotal-golang/lager"
)

type Backup struct {
	Config *brokerconfig.Config
	Logger lager.Logger
}

func (backup Backup) Create(instanceID string) error {
	bucket, err := backup.getOrCreateBucket()
	if err != nil {
		return err
	}

	if err = backup.createSnapshot(instanceID); err != nil {
		return err
	}

	pathToRdbFile := filepath.Join(backup.Config.RedisConfiguration.InstanceDataDirectory, instanceID, "db", "dump.rdb")

	if !fileExists(pathToRdbFile) {
		backup.Logger.Info("dump.rdb not found, skipping instance backup", lager.Data{
			"Local file": pathToRdbFile,
		})
		return nil
	}

	return backup.uploadToS3(instanceID, pathToRdbFile, bucket)
}

func (backup Backup) getOrCreateBucket() (s3bucket.Bucket, error) {
	s3Client := s3bucket.NewClient(
		backup.Config.RedisConfiguration.BackupConfiguration.EndpointUrl,
		backup.Config.RedisConfiguration.BackupConfiguration.S3Region,
		backup.Config.RedisConfiguration.BackupConfiguration.AccessKeyId,
		backup.Config.RedisConfiguration.BackupConfiguration.SecretAccessKey,
	)

	return s3Client.GetOrCreate(backup.Config.RedisConfiguration.BackupConfiguration.BucketName)
}

func (backup Backup) createSnapshot(instanceID string) error {
	client, err := backup.buildRedisClient(instanceID)
	if err != nil {
		return err
	}

	return client.CreateSnapshot(backup.Config.RedisConfiguration.BackupConfiguration.BGSaveTimeoutSeconds)
}

func (backup Backup) buildRedisClient(instanceID string) (*client.Client, error) {
	localRepo := redis.LocalRepository{RedisConf: backup.Config.RedisConfiguration}
	instance, err := localRepo.FindByID(instanceID)
	if err != nil {
		return nil, err
	}

	instanceConf, err := redisconf.Load(localRepo.InstanceConfigPath(instanceID))
	if err != nil {
		return nil, err
	}

	return client.Connect(instance.Host, uint(instance.Port), instance.Password, instanceConf)
}

func (backup Backup) uploadToS3(instanceID, pathToRdbFile string, bucket s3bucket.Bucket) error {
	rdbBytes, err := ioutil.ReadFile(pathToRdbFile)
	if err != nil {
		return err
	}

	remotePath := fmt.Sprintf("%s/%s", backup.Config.RedisConfiguration.BackupConfiguration.Path, instanceID)

	backup.Logger.Info("Backing up instance", lager.Data{
		"Local file":  pathToRdbFile,
		"Remote file": remotePath,
	})

	return bucket.Upload(rdbBytes, remotePath)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}