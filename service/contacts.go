package service

import "database/sql"

// Contact describes a contact in our database.
type Contact struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ===== ADD CONTACT ===================================================================================================

// AddContact inserts a new contact into the database.
func (db *Database) AddContact(c Contact) (int64, error) {
	var contactID int64
	err := db.Write(func(tx *Transaction) {
		contactID = tx.AddContact(c)
	})

	return contactID, err
}

// AddContact inserts a new contact within the transaction.
func (tx *Transaction) AddContact(c Contact) int64 {
	rs, err := tx.Exec(
		"INSERT INTO contacts (email, name) VALUES (?, ?)",
		c.Email,
		c.Name,
	)
	if err != nil {
		panic(err)
	}

	id, err := rs.LastInsertId()
	if err != nil {
		panic(err)
	}

	return id
}

// ===== GET CONTACT ===================================================================================================

// GetContactByEmail reads a Contact from the Database.
func (db *Database) GetContactByEmail(email string) (*Contact, error) {
	var contact *Contact
	err := db.Read(func(tx *Transaction) {
		contact = tx.GetContactByEmail(email)
	})

	return contact, err
}

// GetContactByEmail finds a contact given an email address. `nil` is returned if the Contact doesn't exist in the DB.
func (tx *Transaction) GetContactByEmail(email string) *Contact {
	row := tx.QueryRow(
		"SELECT id, email, name FROM contacts WHERE email = ?",
		email,
	)

	var contact Contact
	err := row.Scan(&contact.ID, &contact.Email, &contact.Name)
	if err == nil {
		return &contact
	} else if err == sql.ErrNoRows {
		return nil
	} else {
		panic(err)
	}
}
