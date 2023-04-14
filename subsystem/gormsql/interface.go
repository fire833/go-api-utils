package gormsql

import (
	"sync"

	"github.com/fire833/go-api-utils/manager"
	"gorm.io/gorm"
)

const GormSQLSubsystemName = "gormsql"

type GormSQLManager struct {
	manager.DefaultSubsystem

	db *gorm.DB
}

func New() *GormSQLManager {
	return &GormSQLManager{
		db: nil,
	}
}

func (g *GormSQLManager) Initialize(wg *sync.WaitGroup, reg manager.AppRegistration) error {
	defer wg.Done()
	return nil
}

func (g *GormSQLManager) Name() string { return GormSQLSubsystemName }

func (g *GormSQLManager) Reload(wg *sync.WaitGroup) {
	defer wg.Done()
}

func (g *GormSQLManager) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
}
