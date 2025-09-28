package postgres

const (
	// register and login
	createNewUser          = "INSERT INTO users(login, password) VALUES ($1, $2) RETURNING user_id;"
	checkUserIsExists      = "SELECT user_id FROM users WHERE login = $1;"
	getUserPasswordByLogin = "SELECT user_id, password FROM users WHERE login = $1;"

	// credentials
	createCredentials              = "INSERT INTO credentials(login, password, user_id, meta) VALUES ($1, $2, $3, $4) RETURNING credentials_id;"
	getCredentialsByUserID         = "SELECT * FROM credentials WHERE user_id = $1;"
	getCredentialsByIDAndUserID    = "SELECT * FROM credentials WHERE credentials_id = $1 AND user_id = $2;"
	updateCredentialsByIDAndUserID = "UPDATE credentials SET login = $1, password = $2, meta = $3 WHERE credentials_id = $4 AND user_id = $5;"
	deleteCredentialsByIDAndUserID = "DELETE FROM credentials WHERE credentials_id = $1 AND user_id = $2;"

	// text data
	createTextData              = "INSERT INTO text_data(text, user_id, meta) VALUES ($1, $2, $3) RETURNING text_data_id;"
	getTextDataByUserID         = "SELECT * FROM text_data WHERE user_id = $1;"
	getTextDataByIDAndUserID    = "SELECT * FROM text_data WHERE text_data_id = $1 AND user_id = $2;"
	updateTextDataByIDAndUserID = "UPDATE text_data SET text = $1, meta = $2 WHERE text_data_id = $3 AND user_id = $4;"
	deleteTextDataByIDAndUserID = "DELETE FROM text_data WHERE text_data_id = $1 AND user_id = $2;"

	// bank cards
	createBankCards              = "INSERT INTO bank_cards(card_number, holder, cvc, expiration_date, user_id, meta) VALUES ($1, $2, $3, $4, $5, $6) RETURNING credentials_id;"
	getBankCardsByUserID         = "SELECT * FROM bank_cards WHERE user_id = $1;"
	getBankCardsByIDAndUserID    = "SELECT * FROM bank_cards WHERE bank_cards_id = $1 AND user_id = $2;"
	updateBankCardsByIDAndUserID = "UPDATE bank_cards SET card_number = $1, holder = $2, cvc = $3, expiration_date = $4, meta = $5 WHERE bank_cards_id = $6 AND user_id = $7;"
	deleteBankCardsByIDAndUserID = "DELETE FROM bank_cards WHERE bank_cards_id = $1 AND user_id = $2;"

	// binaries
	createBinariesData              = "INSERT INTO binaries_data(binary_data, user_id, meta) VALUES ($1, $2, $3) RETURNING credentials_id;"
	getBinariesDataByUserID         = "SELECT * FROM binaries_data WHERE user_id = $1;"
	getBinariesDataByIDAndUserID    = "SELECT * FROM binaries_data WHERE binaries_data_id = $1 AND user_id = $2;"
	updateBinariesDataByIDAndUserID = "UPDATE binaries_data SET binary_data = $1, meta = $2 WHERE bank_cards_id = $3 AND user_id = $4;"
	deleteBinariesDataByIDAndUserID = "DELETE FROM binaries_data WHERE binaries_data_id = $1 AND user_id = $2;"
)
