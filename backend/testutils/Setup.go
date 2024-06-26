package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/arthur-tragante/liven-code-test/models"
)

type TestDBSetup struct {
	DB        *gorm.DB
	Container testcontainers.Container
}

func SetupTestDB(assert *assert.Assertions) *TestDBSetup {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(5 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(err)

	host, err := container.Host(ctx)
	assert.NoError(err)

	port, err := container.MappedPort(ctx, "5432")
	assert.NoError(err)

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=testdb sslmode=disable", host, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(err)

	err = db.AutoMigrate(&models.User{}, &models.Address{})
	assert.NoError(err)

	return &TestDBSetup{
		DB:        db,
		Container: container,
	}
}

func (setup *TestDBSetup) TearDown(assert *assert.Assertions) {
	err := setup.Container.Terminate(context.Background())
	assert.NoError(err)
}