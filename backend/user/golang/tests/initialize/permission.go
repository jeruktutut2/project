package initialize

import (
	"context"
	"database/sql"
	"log"
)

func CreateTablePermission(db *sql.DB, ctx context.Context) {
	query := `CREATE TABLE permission (
  		id int NOT NULL AUTO_INCREMENT,
  		permission varchar(50) NOT NULL,
  		PRIMARY KEY (id),
  		UNIQUE KEY permission_unique_index (permission)
	) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when creating table permission:", err.Error())
	}
	log.Println("create table permission succedded")
}

func CreateDataPermission(db *sql.DB, ctx context.Context) {
	query := `INSERT INTO permission (id, permission) VALUES (1, "ADMINISTRATOR"), (2, "CREATE_PERMISSION"), (3, "READ_PERMISSION"), (4, "UPDATE_PERMISSION"), (5, "DELETE_PERMISSION");`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when creating data permission:", err.Error())
	}
	log.Println("create data permission succedded")
}

func DropTablePermission(db *sql.DB, ctx context.Context) {
	query := `DROP TABLE IF EXISTS permission;`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when dropping table permission:", err.Error())
	}
	log.Println("drop table permission succedded")
}
