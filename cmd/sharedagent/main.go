package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi/auth"
	"github.com/pivotal-cf/cf-redis-broker/agentconfig"
	"github.com/pivotal-cf/cf-redis-broker/availability"
	"github.com/pivotal-cf/cf-redis-broker/redisconf"
	"github.com/pivotal-cf/cf-redis-broker/resetter"
	"github.com/pivotal-cf/cf-redis-broker/sharedagentapi"
	"github.com/pivotal-cf/cf-redis-broker/brokerconfig"
	"github.com/pivotal-cf/cf-redis-broker/redis"
)

type portChecker struct{}

func (portChecker) Check(address *net.TCPAddr, timeout time.Duration) error {
	return availability.Check(address, timeout)
}

func main() {
	configPath := flag.String("sharedAgentConfig", "", "Shared agent config yaml")
	flag.Parse()

	logger := lager.NewLogger("redis-shared-agent")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	config, err := agentconfig.Load(*configPath)
	if err != nil {
		logger.Fatal("Error loading config file", err, lager.Data{
			"path": *configPath,
		})
	}

	templateRedisConf(config, logger)

	redisResetter := resetter.New(
		config.DefaultConfPath,
		config.ConfPath,
		portChecker{},
	)
	redisResetter.Monit.SetExecutable(config.MonitExecutablePath)

	// config.ConfPath => points to a single redis.conf (not used for shared node as we have multiple redis.conf files)

	handler := auth.NewWrapper(
		config.AuthConfiguration.Username,
		config.AuthConfiguration.Password,
	).Wrap(
		sharedagentapi.New(redisResetter, createLocalRepo(logger, configPath)),
	)

	http.Handle("/", handler)
	logger.Fatal("http-listen", http.ListenAndServe("localhost:"+config.Port, nil))
}

// demonstrate how one could create a localRepo (if needed later)
func createLocalRepo(logger lager.Logger, configPath *string) *redis.LocalRepository {
	brokerConfigPath := "/Users/torsten/go/src/github.com/pivotal-cf/cf-redis-broker/configmigratorintegration/assets/broker.yml"
	logger.Info("Config File: " + brokerConfigPath)
	brokerConfig, err := brokerconfig.ParseConfig(brokerConfigPath)
	if err != nil {
		logger.Fatal("Error parsing config file", err, lager.Data{
			"path": *configPath,
		})
	}
	localRepo := redis.NewLocalRepository(brokerConfig.RedisConfiguration, logger)
	return localRepo
}

func templateRedisConf(config *agentconfig.Config, logger lager.Logger) {
	newConfig, err := redisconf.Load(config.DefaultConfPath)
	if err != nil {
		logger.Fatal("Error loading default redis.conf", err, lager.Data{
			"path": config.DefaultConfPath,
		})
	}

	if fileExists(config.ConfPath) {
		existingConf, err := redisconf.Load(config.ConfPath)
		if err != nil {
			logger.Fatal("Error loading existing redis.conf", err, lager.Data{
				"path": config.ConfPath,
			})
		}
		err = newConfig.InitForDedicatedNode(existingConf.Password())
	} else {
		err = newConfig.InitForDedicatedNode()
	}

	if err != nil {
		logger.Fatal("Error initializing redis conf for dedicated node", err)
	}

	err = newConfig.Save(config.ConfPath)
	if err != nil {
		logger.Fatal("Error saving redis.conf", err, lager.Data{
			"path": config.ConfPath,
		})
	}

	logger.Info("Finished writing redis.conf", lager.Data{
		"path": config.ConfPath,
		"conf": newConfig,
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
