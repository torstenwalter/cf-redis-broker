package restoreintegration_test

// restore <INSTANCE_ID> /path/to/dump.rdb

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf/cf-redis-broker/restoreconfig"
)

var _ = Describe("restore dedicated", func() {
	var restoreCommand *exec.Cmd

	var instanceID string
	var sourceRdbPath string
	var testDataDir string
	var monitLogFile string
	var monitLogDir string
	var redisSession *gexec.Session

	var config restoreconfig.Config

	BeforeEach(func() {
		err := copyFile(filepath.Join("..", "brokerintegration", "assets", "redis.conf"), "/tmp/redis.conf")
		Ω(err).ShouldNot(HaveOccurred())
		err = copyFile(filepath.Join("assets", "monit"), "/tmp/monit")
		Ω(err).ShouldNot(HaveOccurred())
		err = os.Chmod("/tmp/monit", 0755)
		Ω(err).ShouldNot(HaveOccurred())

		configPath := filepath.Join("assets", "restore-dedicated.yml")
		config, _ = restoreconfig.Load(configPath)

		instanceID = "test_instance"
		testDataDir = filepath.Join(config.RedisDataDirectory)
		os.RemoveAll(testDataDir)
		os.MkdirAll(testDataDir, 0777)

		monitLogDir, err = ioutil.TempDir("", "monit-test-logs")
		Expect(err).NotTo(HaveOccurred())

		monitLogFile = filepath.Join(monitLogDir, "monit.log")

		sourceRdbPath = filepath.Join("assets", "dump.rdb")
		restoreCommand = exec.Command(restoreExecutablePath, instanceID, sourceRdbPath)
		restoreCommand.Env = append(os.Environ(), "RESTORE_CONFIG_PATH="+configPath)
		restoreCommand.Env = append(restoreCommand.Env, "MONIT_LOG_FILE="+monitLogFile)

		fakeChownPath := "assets"
		for i, envVar := range restoreCommand.Env {
			parts := strings.Split(envVar, "=")
			if parts[0] == "PATH" {
				path := fakeChownPath + ":" + parts[1]

				restoreCommand.Env[i] = "PATH=" + path
			}
		}

		pidfilePath := filepath.Join(testDataDir, "redis-server.pid")
		redisCmd := exec.Command("redis-server",
			"--pidfile", pidfilePath,
			"--daemonize", "yes",
		)

		redisSession, err = gexec.Start(redisCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		pidFileWritten := make(chan bool)
		go func(c chan<- bool) {
			for {
				if _, err := os.Stat(pidfilePath); !os.IsNotExist(err) {
					c <- true
					break
				}
				time.Sleep(50 * time.Millisecond)
			}
		}(pidFileWritten)

		// wait for redis to write pid file
		select {
		case <-pidFileWritten:
			break
		case <-time.After(30 * time.Second):
			Fail("Test timed out waiting for redis to write PID file.")
		}
	})

	AfterEach(func() {
		pid, err := config.InstancePid(instanceID)
		if err == nil {
			syscall.Kill(pid, syscall.SIGKILL)
		}

		os.RemoveAll(monitLogDir)

		Eventually(redisSession, "20s").Should(gexec.Exit(0))
	})

	// It("exits with a non zero status if no arguments are provided", func() {
	// 	restoreCommand.Args = []string{}
	// 	session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	Eventually(session, "20s").Should(gexec.Exit(1))
	// 	Eventually(session.Err).Should(gbytes.Say("usage: restore <instance_id> <rdb_path>"))
	// })

	// It("exits with a non zero status if the instance directory does not exist", func() {
	// 	restoreCommand.Args = []string{restoreCommand.Args[0], "foo", sourceRdbPath}
	// 	session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	Eventually(session.Err).Should(gbytes.Say("Instance not found"))
	// 	Eventually(session, "20s").Should(gexec.Exit(1))
	// })

	// It("exits with a non zero status if the RDB file does not exist", func() {
	// 	restoreCommand.Args = []string{restoreCommand.Args[0], instanceID, "bar"}
	// 	session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	Eventually(session.Err).Should(gbytes.Say("RDB file not found"))
	// 	Eventually(session, "20s").Should(gexec.Exit(1))
	// })

	// It("exits with a non zero status if the config cannot be loaded", func() {
	// 	restoreCommand.Env = []string{"RESTORE_CONFIG_PATH=foo"}
	// 	session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	Eventually(session.Err).Should(gbytes.Say("Could not load config"))
	// 	Eventually(session, "20s").Should(gexec.Exit(1))
	// })

	It("exits successfully if the instance and the RDB file exist", func() {
		session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "20s").Should(gexec.Exit(0))
	})

	It("stops redis", func() {
		session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(redisSession, "20s").Should(gexec.Exit(0))
		Eventually(session, "20s").Should(gexec.Exit(0))
	})

	It("stops and then starts the process-watcher", func() {
		session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "20s").Should(gexec.Exit(0))

		monitLogBytes, err := ioutil.ReadFile(monitLogFile)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(monitLogBytes)).To(ContainSubstring("stopping process-watcher"))
		Expect(string(monitLogBytes)).To(ContainSubstring("starting process-watcher"))
	})

	It("creates a new RDB file in the instance directory", func() {
		newRdbPath := filepath.Join(testDataDir, "dump.rdb")

		_, err := os.Stat(newRdbPath)
		Expect(os.IsNotExist(err)).To(BeTrue())

		session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "20s").Should(gexec.Exit(0))

		copiedFileContents, err := ioutil.ReadFile(newRdbPath)
		Expect(err).NotTo(HaveOccurred())
		sourceFileContents, err := ioutil.ReadFile(sourceRdbPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(copiedFileContents).To(Equal(sourceFileContents))
	})

	It("creates a new AOF file in the instance directory", func() {
		aofPath := filepath.Join(testDataDir, "appendonly.aof")

		_, err := os.Stat(aofPath)
		Expect(os.IsNotExist(err)).To(BeTrue())

		session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "20s").Should(gexec.Exit(0))

		fileContents, err := ioutil.ReadFile(aofPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(fileContents)).To(ContainSubstring("TEST_KEY"))
	})

	It("does not leave redis running", func() {
		session, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, "20s").Should(gexec.Exit(0))

		pkillCommand := exec.Command("pkill", "redis-server")
		pkillSession, err := gexec.Start(pkillCommand, GinkgoWriter, GinkgoWriter)
		// pkill returns 1 if there is nothing for it to kill
		Eventually(pkillSession).Should(gexec.Exit(1))
	})
})
