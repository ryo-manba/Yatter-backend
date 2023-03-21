package object

type (
	StatusID = int64

	// Account status
	Status struct {
		// The internal ID of the status
		ID StatusID `json:"-" db:"id"`

		// The accountId of the status
		//		AccountID AccountID `json:"account_id" db:"account_id"`

		// The content of the status
		Content string `json:"content" db:"content"`

		// The time the status was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`

		// The account of the status
		Account *Account `db:"-"`
	}
)