package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	*log.Logger
	visible      bool
	stickyStill  bool
	autoRaise    bool
	randr        bool
	xinerama     bool
	desktopCount int
	s            store
}

var NoSetting = Xrror("Setting \"%s\" is unavailable").Out

func (s *Settings) query(key string) (*storeItem, error) {
	if v, ok := s.s[key]; ok {
		return v, nil
	}
	return nil, NoSetting(key)
}

func (e Settings) String(key string) string {
	s, _ := e.query(key)
	if s != nil {
		return s.String()
	}
	return ""
}

func (e Settings) Bool(key string) bool {
	s, _ := e.query(key)
	if s != nil {
		return s.Bool()
	}
	return false
}

func (e Settings) Int(key string) int {
	s, _ := e.query(key)
	if s != nil {
		return s.Int()
	}
	return 0
}

func (e Settings) Float(key string) float64 {
	s, _ := e.query(key)
	if s != nil {
		return s.Float()
	}
	return 0.0
}

func (e Settings) List(key string) []string {
	s, _ := e.query(key)
	if s != nil {
		return s.List()
	}
	return []string{}
}

type store map[string]*storeItem

func (s store) Add(k, v string) {
	s[k] = NewStoreItem(k, v)
}

func (s store) Insert(i *storeItem) {
	s[i.key] = i
}

type storeItem struct {
	key   string
	value string
}

func NewStoreItem(k, v string) *storeItem {
	return &storeItem{
		key:   k,
		value: v,
	}
}

func (i *storeItem) String() string {
	return i.value
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
	if value, ok := boolString[strings.ToLower(i.value)]; ok {
		return value
	}
	return false
}

func (i *storeItem) Float() float64 {
	if value, err := strconv.ParseFloat(i.value, 64); err == nil {
		return value
	}
	return 0.0
}

func (i *storeItem) Int() int {
	if value, err := strconv.Atoi(i.value); err == nil {
		return value
	}
	return 0
}

func (i *storeItem) Int64() int64 {
	if value, err := strconv.ParseInt(i.value, 10, 64); err == nil {
		return value
	}
	return -1
}

func (i *storeItem) List(l ...string) []string {
	if len(l) > 0 {
		list := strings.Split(i.value, ",")
		for _, item := range l {
			list = doAdd(item, list)
		}
		i.value = strings.Join(list, ",")
	}
	return strings.Split(i.value, ",")
}

var (
	verbose       bool
	ConfigHomeEnv string = "XDG_CONFIG_HOME"
	ConfigFile    string = "euclid/euclidrc"
	ConfigPath    string
	socketEnv     string = "SCPWM_SOCKET"
	socketPathTpl string = "/tmp/scpwm%s_%d_%d-socket"
	socketPath    string
	ewmhSupported []string = []string{
		"_NET_SUPPORTED",
		"_NET_SUPPORTING_WM_CHECK",
		"_NET_DESKTOP_NAMES",
		"_NET_NUMBER_OF_DESKTOPS",
		"_NET_CURRENT_DESKTOP",
		"_NET_CLIENT_LIST",
		"_NET_ACTIVE_WINDOW",
		"_NET_CLOSE_WINDOW",
		"_NET_WM_DESKTOP",
		"_NET_WM_STATE",
		"_NET_WM_STATE_FULLSCREEN",
		"_NET_WM_STATE_STICKY",
		"_NET_WM_STATE_DEMANDS_ATTENTION",
		"_NET_WM_WINDOW_TYPE",
		"_NET_WM_WINDOW_TYPE_DOCK",
		"_NET_WM_WINDOW_TYPE_DESKTOP",
		"_NET_WM_WINDOW_TYPE_NOTIFICATION",
		"_NET_WM_WINDOW_TYPE_DIALOG",
		"_NET_WM_WINDOW_TYPE_UTILITY",
		"_NET_WM_WINDOW_TYPE_TOOLBAR",
	}
)

func DefaultSettings() *Settings {
	s := make(store)
	s.Add("WindowManagerName", "scpwm")
	s.Add("ExternalRulesCommand", "TBD")
	s.Add("StatusPrefix", "W")
	s.Add("SplitRatio", "0.5")
	s.Add("WindowGap", "6")
	s.Add("BorderWidth", "1")
	s.Add("CenterPseudoTiled", "true")
	s.Add("RecordHistory", "true")
	s.Add("DefaultMonitorName", "MONITOR")
	s.Add("DefaultDesktopName", "DESKTOP")
	s.Add("InitialPolarity", "right")

	n := &storeItem{key: "ewmhSupported"}
	n.List(ewmhSupported...)
	s.Insert(n)

	return &Settings{
		Logger:      log.New(os.Stderr, "[EUCLID] ", log.Ldate|log.Lmicroseconds),
		visible:     true,
		stickyStill: true,
		autoRaise:   true,
		s:           s,
	}
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

//func init() {
//numMonitors, numDesktops, numClients := 0, 0, 0
//monitorUid, desktopUid := 0, 0
//mon, monHead, monTail, priMon := nil, nil, nil, nil
//historyHead, historyTail, historyNeedle := nil, nil, nil
//ruleHead, ruleTail := nil, nil
//stackHead, stackTail := nil, nil
//subscribeHead, subscribetail := nil, nil
//pendingRuleHead, pendingRuleTail := nil, nil
//lastMotionTime, lastMotionX, lastMotionY := 0, 0, 0
//visible, autoRaise, stickyStill, RecordHistory = true, true, true, true
//randrBase := 0
//exit_status = 0;
//}

//child_polarity_t initial_polarity;
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
