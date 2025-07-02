package repo

func NewClient() *Client {
	return &Client{}
}

type Client struct{}

// Reconcile reads .alveus.yml within a branch in a repo
// it then (re)generates various files and commits those
func (c *Client) Reconcile(branch string) {}
