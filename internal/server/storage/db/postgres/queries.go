package postgres

const (
	// register and login
	createNewUser          = "INSERT INTO users(login, password) VALUES ($1, $2) RETURNING user_id;"
	checkUserIsExists      = "SELECT user_id FROM users WHERE login = $1;"
	getUserPasswordByLogin = "SELECT user_id, password FROM users WHERE login = $1;"

	// credentials
	createCredentials              = "INSERT INTO credentials(login, password, user_id) VALUES ($1, $2, $3) RETURNING credentials_id;"
	getCredentialsByUserID         = "SELECT * FROM credentials WHERE user_id = $1;"
	getCredentialsByIDAndUserID    = "SELECT * FROM credentials WHERE credentials_id = $1 AND user_id = $2;"
	updateCredentialsByIDAndUserID = "UPDATE credentials SET login = $1, password = $2 WHERE credentials_id = $3 AND user_id = $4;"
	deleteCredentialsByIDAndUserID = "DELETE FROM credentials WHERE credentials_id = $1 AND user_id = $2;"

	// text data
	createTextData              = "INSERT INTO text_data(text, user_id) VALUES ($1, $2) RETURNING text_data_id;"
	getTextDataByUserID         = "SELECT * FROM text_data WHERE user_id = $1;"
	getTextDataByIDAndUserID    = "SELECT * FROM text_data WHERE text_data_id = $1 AND user_id = $2;"
	updateTextDataByIDAndUserID = "UPDATE text_data SET text = $1 WHERE text_data_id = $2 AND user_id = $3;"
	deleteTextDataByIDAndUserID = "DELETE FROM text_data WHERE text_data_id = $1 AND user_id = $2;"
)
