package postgres

import (
	"fmt"
)

func SetupDB() {
	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    login VARCHAR(256) NOT NULL UNIQUE, 
    password VARCHAR(256) NOT NULL,     
    email VARCHAR(256)                   
);

	CREATE TABLE IF NOT EXISTS tokens (
    user_id UUID NOT NULL,             
    token VARCHAR NOT NULL,             
    UNIQUE(user_id),                    
    UNIQUE(token),                      
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE              
);
`

	_, err := Connection.conn.Exec(query)
	if err != nil {
		fmt.Println("Что-то не так с бд или запросом при ее формировании, возможно, косяк в тесте" + err.Error())
	}
}
