package clients

import "github.com/thrisp/scpwm/euclid/branch"

func New() *branch.Branch {
	return branch.New("clients")
}

type MatchClient func(Client) bool

func seek(clients *branch.Branch, fn MatchClient) Client {
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

func seekOffset(clients *branch.Branch, fn MatchClient, offset int) Client {
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

func isFocused(c Client) bool {
	if c.Focused() {
		return true
	}
	return false
}

func Focused(clients *branch.Branch) Client {
	return seek(clients, isFocused)
}

func seekAny(clients *branch.Branch, fn MatchClient) []Client {
	var ret []Client
	curr := clients.Front()
	for curr != nil {
		client := curr.Value.(Client)
		if match := fn(client); match {
			ret = append(ret, client)
		}
		curr = curr.Next()
	}
	return ret
}

func All(clients *branch.Branch) []Client {
	fn := func(c Client) bool {
		return true
	}
	return seekAny(clients, fn)
}
