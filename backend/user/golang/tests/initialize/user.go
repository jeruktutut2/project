package initialize

import (
	"context"
	"database/sql"
	"log"
)

func CreateTableUser(db *sql.DB, ctx context.Context) {
	query := `CREATE TABLE user (
  		id int NOT NULL AUTO_INCREMENT,
  		username varchar(50) NOT NULL,
  		email varchar(100) NOT NULL,
  		password varchar(100) NOT NULL,
  		utc varchar(6) NOT NULL,
  		created_at bigint NOT NULL,
  		PRIMARY KEY (id),
  		UNIQUE KEY username_unique_index (username),
  		UNIQUE KEY email_unique_index (email)
	) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when creating table user:", err.Error())
	}
	log.Println("create table user succedded")
}

func CreateDataUser(db *sql.DB, ctx context.Context) {
	query := `INSERT INTO user (id,username,email,password,utc,created_at) VALUES (1,"username","email@email.com","$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG","utc",1695095017);`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when creating data user:", err.Error())
	}
	log.Println("create data user succedded")
}

func DropTableUser(db *sql.DB, ctx context.Context) {
	query := `DROP TABLE IF EXISTS user;`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when dropping table user:", err.Error())
	}
	log.Println("drop table user succedded")
}
