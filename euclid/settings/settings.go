package settings

import (
	"fmt"
	"strconv"
	"strings"
)

type Settings interface {
	Getter
	Setter
	Copy() Settings
}

type Getter interface {
	Query(string) (StoreItem, error)
	String(string) string
	Bool(string) bool
	Int(string) int
	Float(string) float64
	List(string) []string
}

type Setter interface {
	Add(string, string)
	Insert(StoreItem)
}

type settings struct {
	s store
}

type store map[string]StoreItem

func (s store) add(k, v string) {
	s[k] = NewStoreItem(k, v)
}

func (s store) insert(i StoreItem) {
	s[i.Key()] = i
}

var NoSetting = Xrror("Setting \"%s\" is unavailable").Out

func (s *settings) Query(key string) (StoreItem, error) {
	if v, ok := s.s[key]; ok {
		return v, nil
	}
	return nil, NoSetting(key)
}

func (s *settings) String(key string) string {
	q, _ := s.Query(key)
	if q != nil {
		return q.String()
	}
	return ""
}

func (s *settings) Bool(key string) bool {
	q, _ := s.Query(key)
	if q != nil {
		return q.Bool()
	}
	return false
}

func (s *settings) Int(key string) int {
	q, _ := s.Query(key)
	if q != nil {
		return q.Int()
	}
	return 0
}

func (s *settings) Float(key string) float64 {
	q, _ := s.Query(key)
	if q != nil {
		return q.Float()
	}
	return 0.0
}

func (s *settings) List(key string) []string {
	q, _ := s.Query(key)
	if q != nil {
		return q.List()
	}
	return []string{}
}

func (s *settings) Pad(key string) [4]int {
	q, _ := s.Query(key)
	if q != nil {
		return q.Pad()
	}
	return [4]int{0, 0, 0, 0}
}

func (s *settings) Add(k, v string) {
	s.s.add(k, v)
}

func (s *settings) Insert(i StoreItem) {
	s.s.insert(i)
}

func (s *settings) Copy() Settings {
	nst := make(store)
	for k, v := range s.s {
		nst[k] = v
	}
	return &settings{s: nst}
}

type StoreItem interface {
	Set(string, string, string)
	Key() string
	String() string
	Bool() bool
	Float() float64
	Int() int
	Int64() int64
	List(...string) []string
	Pad() [4]int
}

type storeItem struct {
	key        string
	value      string
	localvalue string
	global     bool
}

func NewStoreItem(k, v string) StoreItem {
	return &storeItem{
		key:    k,
		value:  v,
		global: true,
	}
}

func (i *storeItem) valueScope() string {
	if !i.global {
		return i.localvalue
	}
	return i.value
}

func (i *storeItem) Set(scope, key, value string) {
	switch scope {
	case "global":
		i.value = value
	case "local":
		i.localvalue = value
	}
	i.key = key
}

func (i *storeItem) Key() string {
	return i.key
}

func (i *storeItem) String() string {
	return i.valueScope()
}

var boolString = map[string]bool{
	"t":     true,
	"true":  true,
	"y":     true,
	"yes":   true,
	"on":    true,
	"1":     true,
	"f":     false,
	"false": false,
	"n":     false,
	"no":    false,
	"off":   false,
	"0":     false,
}

func (i *storeItem) Bool() bool {
	if value, ok := boolString[strings.ToLower(i.valueScope())]; ok {
		return value
	}
	return false
}

func (i *storeItem) Float() float64 {
	if value, err := strconv.ParseFloat(i.valueScope(), 64); err == nil {
		return value
	}
	return 0.0
}

func (i *storeItem) Int() int {
	if value, err := strconv.Atoi(i.valueScope()); err == nil {
		return value
	}
	return 0
}

func (i *storeItem) Int64() int64 {
	if value, err := strconv.ParseInt(i.valueScope(), 10, 64); err == nil {
		return value
	}
	return -1
}

func isAppendable(s string, ss []string) bool {
	for _, x := range ss {
		if x == s {
			return false
		}
	}
	return true
}

func doAdd(s string, ss []string) []string {
	if isAppendable(s, ss) {
		ss = append(ss, s)
	}
	return ss
}

func (i *storeItem) List(l ...string) []string {
	if len(l) > 0 {
		list := strings.Split(i.valueScope(), ",")
		for _, item := range l {
			list = doAdd(item, list)
		}
		lst := strings.Join(list, ",")
		if i.global {
			i.value = lst
		} else {
			i.localvalue = lst
		}
	}
	return strings.Split(i.valueScope(), ",")
}

func Pad(key string, right, up, left, down int) StoreItem {
	val := fmt.Sprintf("%d,%d,%d,%d", right, up, left, down)
	return NewStoreItem(key, val)
}

func DefaultPad() StoreItem {
	return Pad("DefaultPad", 0, 0, 0, 0)
}

func (i *storeItem) Pad() [4]int {
	var ret [4]int
	lst := i.List()
	for i, v := range lst {
		if num, err := strconv.Atoi(v); err == nil {
			ret[i] = num
		}
	}
	return ret
}

//WmName                   string
//ExternalRulesCommand     string
//StatusPrefix             string
//SplitRatio               float64
//WindowGap                int
//BorderWidth              uint
//HistoryAwareFocus        bool
//FocusByDistance          bool
//BorderlessMonocle        bool
//GaplessMonocle           bool
//LeafMonocle              bool
//FocusFollowsPointer      bool
//PointerFollowsFocus      bool
//PointerFollowsMonitor    bool
//AutoAlternate            bool
//AutoCancel               bool
//ApplyFloatingAtom        bool
//IgnoreEwmhFocus          bool
//CenterPseudoTiled        bool
//RemoveDisabledMonitors   bool
//RemoveUnpluggedMonitors  bool
//MergeOverlappingMonitors bool
//ewmhSupported            []string
//visible                  bool
//autoRaise                bool
//stickyStill              bool
//RecordHistory            bool

/*
	FocusedBorderColor          "#7e7f89"
	ActiveBorderColor           "#545350"
	NormalBorderColor           "#3f3e3b"
	PreselBorderColor           "#e8e8f4"
	FocusedLockedBorderColor   "#c7b579"
	ActiveLockedBorderColor    "#545350"
	NormalLockedBorderColor    "#3f3e3b"
	FocusedStickyBorderColor   "#e3a5da"
	ActiveStickyBorderColor    "#545350"
	NormalStickyBorderColor    "#3f3e3b"
	FocusedPrivateBorderColor  "#42cad9"
	ActivePrivateBorderColor   "#5c5955"
	NormalPrivateBorderColor   "#34322e"
	UrgentBorderColor           "#efa29a"
*/
