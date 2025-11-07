package ldap

// Client is the interface to abstract LDAP operations for tests.
type Client interface {
	Ping() error
	SearchUsers(filter string, limit int) ([]User, error)
	CreateUser(u User) error
}

type User struct {
	DN   string
	UID  string
	Name string
	Mail string
}
