package sqlite

import (
	"database/sql"

	"github.com/Tanibox/tania-core/src/growth/query"
	"github.com/gofrs/uuid"
)

type FarmReadQueryMysql struct {
	DB *sql.DB
}

func NewFarmReadQueryMysql(db *sql.DB) query.FarmReadQuery {
	return FarmReadQueryMysql{DB: db}
}

type farmReadResult struct {
	UID  []byte
	Name string
}

func (s FarmReadQueryMysql) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		farmRead := query.CropFarmQueryResult{}
		rowsData := farmReadResult{}

		err := s.DB.QueryRow("SELECT UID, NAME FROM FARM_READ WHERE UID = ?", uid.Bytes()).Scan(
			&rowsData.UID,
			&rowsData.Name,
		)

		if err != nil && err != sql.ErrNoRows {
			result <- query.Result{Error: err}
		}

		if err == sql.ErrNoRows {
			result <- query.Result{Result: farmRead}
		}

		farmUID, err := uuid.FromBytes(rowsData.UID)
		if err != nil {
			result <- query.Result{Error: err}
		}

		farmRead.UID = farmUID
		farmRead.Name = rowsData.Name

		result <- query.Result{Result: farmRead}
		close(result)
	}()

	return result
}
