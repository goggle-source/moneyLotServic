package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(user string, password string, dbName string, port int, hostDB string) error {

	conn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, hostDB, port, dbName)

	// ищет в корне проекта
	m, err := migrate.New("file://migrations", conn)

	if err != nil {
		fmt.Println("newWithDatabaseInstance", err)
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("up", err)
		return err
	}

	return nil
}
