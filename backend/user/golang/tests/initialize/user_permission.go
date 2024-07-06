package initialize

import (
	"context"
	"database/sql"
	"log"
)

func CreateTableUserPermission(db *sql.DB, ctx context.Context) {
	query := `CREATE TABLE user_permission (
  		id int NOT NULL AUTO_INCREMENT,
  		user_id int NOT NULL,
  		permission_id int NOT NULL,
  		PRIMARY KEY (id),
  		KEY fk_user_permission_user (user_id),
  		KEY fk_user_permission_permission (permission_id),
  		CONSTRAINT user_permission_ibfk_1 FOREIGN KEY (user_id) REFERENCES user (id),
  		CONSTRAINT user_permission_ibfk_2 FOREIGN KEY (permission_id) REFERENCES permission (id)
	) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when creating table user_permission:", err.Error())
	}
	log.Println("create table user_permission succedded")
}

func CreateDataUserPermission(db *sql.DB, ctx context.Context) {
	query := `INSERT INTO user_permission(user_id, permission_id)
		VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5);`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when creating data user_permission:", err.Error())
	}
	log.Println("create data user_permission succedded")
}

func DropTableUserPermission(db *sql.DB, ctx context.Context) {
	query := `DROP TABLE IF EXISTS user_permission;`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalln("error when dropping table user_permission:", err.Error())
	}
	log.Println("drop table user_permission succedded")
}
