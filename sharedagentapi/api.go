package sharedagentapi

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/pivotal-cf/cf-redis-broker/redis"
	"github.com/pborman/uuid"
	"github.com/pivotal-cf/cf-redis-broker/redisconf"
	"strconv"
	"encoding/json"
	"fmt"
	"os"
	"github.com/pivotal-cf/cf-redis-broker/sharedagentconfig"
	"path"
)

type redisResetter interface {
	ResetRedis() error
}

func New(config *sharedagentconfig.Config, resetter redisResetter, localRepo *redis.LocalRepository) http.Handler {
	router := mux.NewRouter()

	router.Path("/createDummyRedisConf").Methods(http.MethodPost).HandlerFunc(createDummyRedisConf(localRepo))

	/*	router.Path("/").
			Methods("DELETE").
			HandlerFunc(resetHandler(resetter))*/

	router.Path("/redis/{instance}/").
		Methods("GET").
		HandlerFunc(credentialsHandler(config.ConfBasePath))

	/*router.Path("/keycount").
		Methods("GET").
		HandlerFunc(keyCountHandler(configPath))*/

	return router
}

func createDummyRedisConf(localRepo *redis.LocalRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Print("Hello World")
		instanceID := uuid.NewRandom().String()
		instance := &redis.Instance{
			ID:   instanceID,
			Host: "127.0.0.1",
			Port: 8080,
		}

		err := localRepo.Setup(instance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func credentialsHandler(configBasePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instance := vars["instance"]

		configPath := path.Join(configBasePath, instance, "redis.conf")

		_, err := os.Stat(configPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("no such redis instances '%s'", instance), http.StatusInternalServerError)
			return
		}

		conf, err := redisconf.Load(configPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		port, err := strconv.Atoi(conf.Get("port"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		password := conf.Get("requirepass")

		credentials := struct {
			Port     int    `json:"port"`
			Password string `json:"password"`
		}{
			Port:     port,
			Password: password,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(credentials)
	}
}
