package clients

import "github.com/thrisp/scpwm/euclid/branch"

func New() *branch.Branch {
	return branch.New("clients")
}

func All(clients *branch.Branch) []Client {
	var ret []Client
	curr := clients.Front()
	for curr != nil {
		c := curr.Value.(Client)
		ret = append(ret, c)
		curr = curr.Next()
	}
	return ret
}

type SelectClient func(Client) bool

func seek(clients *branch.Branch, fn SelectClient) Client {
	curr := clients.Front()
	for curr != nil {
		client := curr.Value.(Client)
		if found := fn(client); found {
			return client
		}
		curr = curr.Next()
	}
	return nil
}

func seekOffset(clients *branch.Branch, fn SelectClient, offset int) Client {
	curr := clients.Front()
	for curr != nil {
		client := curr.Value.(Client)
		if found := fn(client); found {
			switch offset {
			case -1:
				client = curr.PrevContinuous().Value.(Client)
			case 1:
				client = curr.NextContinuous().Value.(Client)
			}
			return client
		}
		curr = curr.Next()
	}
	return nil
}
