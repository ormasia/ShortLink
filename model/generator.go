package model

// 使用数据库自增主键生成ID
func GenerateID() (int64, error) {
	res, err := db.Exec("INSERT INTO id_generator () VALUES ()")
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

/*
CREATE TABLE id_generator (
    id BIGINT AUTO_INCREMENT PRIMARY KEY
);

CREATE TABLE url_mapping (
    short_url VARCHAR(255) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/
