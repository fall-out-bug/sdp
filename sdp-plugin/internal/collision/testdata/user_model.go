package testdata

// User represents a user in the system.
type User struct {
	Email     string
	Name      string
	Role      string
	CreatedAt int64
}

// Profile represents user profile information.
type Profile struct {
	UserID    string
	Bio       string
	AvatarURL string
}
