package register

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create() error
	Update() error
	Register(gin.IRouter)
}

var register = struct {
	sync.RWMutex
	s map[int]map[string]Service
}{
	s: make(map[int]map[string]Service),
}

func Register(version int, name string, s Service) {
	register.Lock()
	defer register.Unlock()

	if _, ok := register.s[version]; !ok {
		register.s[version] = make(map[string]Service)
	}

	if _, ok := register.s[version][name]; ok {
		log.Fatal("duplicate service for same version and name")
	}

	register.s[version][name] = s
}

func Init(version int, router gin.IRouter) error {
	register.Lock()
	defer register.Unlock()
	for i := 0; i <= version; i++ {
		services, ok := register.s[i]
		if !ok {
			continue
		}

		for name, service := range services {
			if err := service.Create(); err != nil {
				log.Printf("[create] service version %d, error: %s", i, err)
				return err
			}

			if err := service.Update(); err != nil {
				log.Printf("[update] service version %d, error: %s", i, err)
				return err
			}

			service.Register(router.Group(fmt.Sprintf("/api/v%d/%s", i, name)))
		}
	}

	return nil
}
