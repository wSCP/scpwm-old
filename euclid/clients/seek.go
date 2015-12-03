package clients

import "github.com/thrisp/scpwm/euclid/branch"

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

func isFocused(c Client) bool {
	if c.Focused() {
		return true
	}
	return false
}

func Focused(clients *branch.Branch) Client {
	return seek(clients, isFocused)
}

func TiledCount(clients *branch.Branch) int {
	var count int
	fn := func(c Client) bool {
		if c.Tiled() || c.Pseudotiled() {
			count++
			return true
		}
		return false
	}
	seek(clients, fn)
	return count
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

func Prev(clients *branch.Branch) Client {
	return nil //seekOffset(clients, isFocused, -1)
}

func Next(clients *branch.Branch) Client {
	return nil //seekOffset(clients, isFocused, 1)
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
