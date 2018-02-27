package sharedagentapi

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/pivotal-cf/cf-redis-broker/redis"
	"github.com/pborman/uuid"
)

type redisResetter interface {
	ResetRedis() error
}

func New(resetter redisResetter, configPath string, localRepo *redis.LocalRepository) http.Handler {
	router := mux.NewRouter()

	router.Path("/createDummyRedisConf").Methods(http.MethodPost).HandlerFunc(createDummyRedisConf(localRepo))

	/*	router.Path("/").
			Methods("DELETE").
			HandlerFunc(resetHandler(resetter))

		router.Path("/").
			Methods("GET").
			HandlerFunc(credentialsHandler(configPath))

		router.Path("/keycount").
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
