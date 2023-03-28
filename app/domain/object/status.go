package object

type (
	StatusID = int64

	// Account status
	Status struct {
		// The internal ID of the status
		ID StatusID `json:"id" db:"id"`

		// The account of the status
		Account *Account `json:"account" db:"-"`

		// The content of the status
		Content string `json:"content" db:"content"`

		// The time the status was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
