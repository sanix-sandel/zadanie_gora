package dbutils

const image = `
		CREATE TABLE IF NOT EXISTS images (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			TITLE VARCHAR(64) ,
			URL TEXT,
			SIZE INTEGER
		)
`
