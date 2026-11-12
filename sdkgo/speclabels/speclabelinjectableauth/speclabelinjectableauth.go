package speclabelinjectableauth

import "strings"

const Key = "injectable_auth"

const (
	User    = "user"
	Dataset = "dataset"
	Never   = "never"
)

const Separator = "|"

const (
	DatasetOrUser = Dataset + Separator + User
	UserOrDataset = User + Separator + Dataset
)

func Join(kinds ...string) string {
	return strings.Join(kinds, Separator)
}
