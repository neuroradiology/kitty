// License: GPLv3 Copyright: 2022, Kovid Goyal, <kovid at kovidgoyal.net>

package loop

import (
	"fmt"
	"strconv"
	"strings"

	"kitty"
)

// key encoding mappings {{{
// start csi mapping (auto generated by gen-key-constants.py do not edit)
var functional_key_number_to_name_map = map[int]string{57344: "ESCAPE", 57345: "ENTER", 57346: "TAB", 57347: "BACKSPACE", 57348: "INSERT", 57349: "DELETE", 57350: "LEFT", 57351: "RIGHT", 57352: "UP", 57353: "DOWN", 57354: "PAGE_UP", 57355: "PAGE_DOWN", 57356: "HOME", 57357: "END", 57358: "CAPS_LOCK", 57359: "SCROLL_LOCK", 57360: "NUM_LOCK", 57361: "PRINT_SCREEN", 57362: "PAUSE", 57363: "MENU", 57364: "F1", 57365: "F2", 57366: "F3", 57367: "F4", 57368: "F5", 57369: "F6", 57370: "F7", 57371: "F8", 57372: "F9", 57373: "F10", 57374: "F11", 57375: "F12", 57376: "F13", 57377: "F14", 57378: "F15", 57379: "F16", 57380: "F17", 57381: "F18", 57382: "F19", 57383: "F20", 57384: "F21", 57385: "F22", 57386: "F23", 57387: "F24", 57388: "F25", 57389: "F26", 57390: "F27", 57391: "F28", 57392: "F29", 57393: "F30", 57394: "F31", 57395: "F32", 57396: "F33", 57397: "F34", 57398: "F35", 57399: "KP_0", 57400: "KP_1", 57401: "KP_2", 57402: "KP_3", 57403: "KP_4", 57404: "KP_5", 57405: "KP_6", 57406: "KP_7", 57407: "KP_8", 57408: "KP_9", 57409: "KP_DECIMAL", 57410: "KP_DIVIDE", 57411: "KP_MULTIPLY", 57412: "KP_SUBTRACT", 57413: "KP_ADD", 57414: "KP_ENTER", 57415: "KP_EQUAL", 57416: "KP_SEPARATOR", 57417: "KP_LEFT", 57418: "KP_RIGHT", 57419: "KP_UP", 57420: "KP_DOWN", 57421: "KP_PAGE_UP", 57422: "KP_PAGE_DOWN", 57423: "KP_HOME", 57424: "KP_END", 57425: "KP_INSERT", 57426: "KP_DELETE", 57427: "KP_BEGIN", 57428: "MEDIA_PLAY", 57429: "MEDIA_PAUSE", 57430: "MEDIA_PLAY_PAUSE", 57431: "MEDIA_REVERSE", 57432: "MEDIA_STOP", 57433: "MEDIA_FAST_FORWARD", 57434: "MEDIA_REWIND", 57435: "MEDIA_TRACK_NEXT", 57436: "MEDIA_TRACK_PREVIOUS", 57437: "MEDIA_RECORD", 57438: "LOWER_VOLUME", 57439: "RAISE_VOLUME", 57440: "MUTE_VOLUME", 57441: "LEFT_SHIFT", 57442: "LEFT_CONTROL", 57443: "LEFT_ALT", 57444: "LEFT_SUPER", 57445: "LEFT_HYPER", 57446: "LEFT_META", 57447: "RIGHT_SHIFT", 57448: "RIGHT_CONTROL", 57449: "RIGHT_ALT", 57450: "RIGHT_SUPER", 57451: "RIGHT_HYPER", 57452: "RIGHT_META", 57453: "ISO_LEVEL3_SHIFT", 57454: "ISO_LEVEL5_SHIFT"}

var csi_number_to_functional_number_map = map[int]int{2: 57348, 3: 57349, 5: 57354, 6: 57355, 7: 57356, 8: 57357, 9: 57346, 11: 57364, 12: 57365, 13: 57345, 14: 57367, 15: 57368, 17: 57369, 18: 57370, 19: 57371, 20: 57372, 21: 57373, 23: 57374, 24: 57375, 27: 57344, 127: 57347}

var letter_trailer_to_csi_number_map = map[string]int{"A": 57352, "B": 57353, "C": 57351, "D": 57350, "E": 57427, "F": 8, "H": 7, "P": 11, "Q": 12, "S": 14}

var tilde_trailers = map[int]bool{57348: true, 57349: true, 57354: true, 57355: true, 57366: true, 57368: true, 57369: true, 57370: true, 57371: true, 57372: true, 57373: true, 57374: true, 57375: true}

// end csi mapping
// }}}

var name_to_functional_number_map map[string]int
var functional_to_csi_number_map map[int]int
var csi_number_to_letter_trailer_map map[int]string

type KeyEventType uint8
type KeyModifiers uint16

const (
	PRESS   KeyEventType = 1
	REPEAT  KeyEventType = 2
	RELEASE KeyEventType = 4
)

const (
	SHIFT     KeyModifiers = 1
	ALT       KeyModifiers = 2
	CTRL      KeyModifiers = 4
	SUPER     KeyModifiers = 8
	HYPER     KeyModifiers = 16
	META      KeyModifiers = 32
	CAPS_LOCK KeyModifiers = 64
	NUM_LOCK  KeyModifiers = 128
)

func (self KeyModifiers) WithoutLocks() KeyModifiers {
	return self & ^(CAPS_LOCK | NUM_LOCK)
}

func (self KeyEventType) String() string {
	switch self {
	case PRESS:
		return "PRESS"
	case REPEAT:
		return "REPEAT"
	case RELEASE:
		return "RELEASE"
	default:
		return fmt.Sprintf("KeyEventType:%d", int(self))
	}
}

func (self KeyModifiers) String() string {
	ans := make([]string, 0)
	if self&SHIFT != 0 {
		ans = append(ans, "shift")
	}
	if self&ALT != 0 {
		ans = append(ans, "alt")
	}
	if self&CTRL != 0 {
		ans = append(ans, "ctrl")
	}
	if self&SUPER != 0 {
		ans = append(ans, "super")
	}
	if self&HYPER != 0 {
		ans = append(ans, "hyper")
	}
	if self&META != 0 {
		ans = append(ans, "meta")
	}
	if self&CAPS_LOCK != 0 {
		ans = append(ans, "caps_lock")
	}
	if self&NUM_LOCK != 0 {
		ans = append(ans, "num_lock")
	}
	return strings.Join(ans, "+")
}

func (self KeyModifiers) HasCapsLock() bool {
	return self&CAPS_LOCK != 0
}

type KeyEvent struct {
	Type         KeyEventType
	Mods         KeyModifiers
	Key          string
	ShiftedKey   string
	AlternateKey string
	Text         string
	Handled      bool

	// The CSI string this key event was decoded from. Empty if not decoded from CSI.
	CSI string
}

func (self *KeyEvent) String() string {
	key := self.Key
	if self.Mods > 0 {
		key = self.Mods.String() + "+" + key
	}
	ans := fmt.Sprint(self.Type, "{ ", key, " ")
	if self.Text != "" {
		ans += "Text: " + self.Text + " "
	}
	if self.ShiftedKey != "" {
		ans += "ShiftedKey: " + self.ShiftedKey + " "
	}
	if self.AlternateKey != "" {
		ans += "AlternateKey: " + self.AlternateKey + " "
	}
	return ans + "}"
}

func (self *KeyEvent) HasCapsLock() bool {
	return self.Mods.HasCapsLock()
}

func KeyEventFromCSI(csi string) *KeyEvent {
	if len(csi) == 0 {
		return nil
	}
	orig_csi := csi
	last_char := csi[len(csi)-1:]
	if !strings.Contains("u~ABCDEHFPQRS", last_char) || (last_char == "~" && (csi == "200~" || csi == "201~")) {
		return nil
	}
	csi = csi[:len(csi)-1]
	sections := strings.Split(csi, ";")

	get_sub_sections := func(section string, missing int) []int {
		p := strings.Split(section, ":")
		ans := make([]int, len(p))
		for i, x := range p {
			if x == "" {
				ans[i] = missing
			} else {
				q, err := strconv.Atoi(x)
				if err != nil {
					return nil
				}
				ans[i] = q
			}
		}
		return ans
	}
	first_section := get_sub_sections(sections[0], 0)
	second_section := []int{}
	third_section := []int{}
	if len(sections) > 1 {
		second_section = get_sub_sections(sections[1], 1)
	}
	if len(sections) > 2 {
		third_section = get_sub_sections(sections[2], 0)
	}
	var ans = KeyEvent{Type: PRESS, CSI: orig_csi}
	var keynum int
	if val, ok := letter_trailer_to_csi_number_map[last_char]; ok {
		keynum = val
	} else {
		if len(first_section) == 0 {
			return nil
		}
		keynum = first_section[0]
	}

	key_name := func(keynum int) string {
		switch keynum {
		case 0:
			return ""
		case 13:
			if last_char == "u" {
				return "ENTER"
			}
			return "F3"
		default:
			if val, ok := csi_number_to_functional_number_map[keynum]; ok {
				keynum = val
			}
			ans := ""
			if val, ok := functional_key_number_to_name_map[keynum]; ok {
				ans = val
			} else {
				ans = string(rune(keynum))
			}
			return ans
		}
	}

	ans.Key = key_name(keynum)
	if len(first_section) > 1 {
		ans.ShiftedKey = key_name(first_section[1])
	}
	if len(first_section) > 2 {
		ans.AlternateKey = key_name(first_section[2])
	}
	if len(second_section) > 0 {
		ans.Mods = KeyModifiers(second_section[0] - 1)
	}
	if len(second_section) > 1 {
		switch second_section[1] {
		case 2:
			ans.Type = REPEAT
		case 3:
			ans.Type = RELEASE
		}
	}
	if len(third_section) > 0 {
		runes := make([]rune, len(third_section))
		for i, ch := range third_section {
			runes[i] = rune(ch)
		}
		ans.Text = string(runes)
	}
	return &ans
}

type ParsedShortcut struct {
	Mods    KeyModifiers
	KeyName string
}

func (self *ParsedShortcut) String() string {
	ans := self.KeyName
	if self.Mods > 0 {
		ans = self.Mods.String() + "+" + ans
	}
	return ans
}

var parsed_shortcut_cache map[string]*ParsedShortcut

func ParseShortcut(spec string) *ParsedShortcut {
	if parsed_shortcut_cache == nil {
		parsed_shortcut_cache = make(map[string]*ParsedShortcut, 128)
	}
	if val, ok := parsed_shortcut_cache[spec]; ok {
		return val
	}
	ospec := spec
	if strings.HasSuffix(spec, "+") {
		ospec = spec[:len(spec)-1] + "plus"
	}
	parts := strings.Split(ospec, "+")
	key_name := parts[len(parts)-1]
	if val, ok := kitty.FunctionalKeyNameAliases[strings.ToUpper(key_name)]; ok {
		key_name = val
	}
	if _, is_functional_key := name_to_functional_number_map[strings.ToUpper(key_name)]; is_functional_key {
		key_name = strings.ToUpper(key_name)
	} else {
		if val, ok := kitty.CharacterKeyNameAliases[strings.ToUpper(key_name)]; ok {
			key_name = val
		}
	}
	ans := ParsedShortcut{KeyName: key_name}
	if len(parts) > 1 {
		for _, q := range parts[:len(parts)-1] {
			val, ok := kitty.ConfigModMap[strings.ToUpper(q)]
			if ok {
				ans.Mods |= KeyModifiers(val)
			} else {
				ans.Mods |= META << 8
			}
		}
	}
	parsed_shortcut_cache[spec] = &ans
	return &ans
}

func (self *KeyEvent) MatchesParsedShortcut(ps *ParsedShortcut, event_type KeyEventType) bool {
	if self.Type&event_type == 0 {
		return false
	}
	mods := self.Mods.WithoutLocks()
	if mods == ps.Mods && self.Key == ps.KeyName {
		return true
	}
	if self.ShiftedKey != "" && mods&SHIFT != 0 && (mods & ^SHIFT) == ps.Mods && self.ShiftedKey == ps.KeyName {
		return true
	}
	return false
}

func (self *KeyEvent) Matches(spec string, event_type KeyEventType) bool {
	return self.MatchesParsedShortcut(ParseShortcut(spec), event_type)
}

func (self *KeyEvent) MatchesPressOrRepeat(spec string) bool {
	return self.MatchesParsedShortcut(ParseShortcut(spec), PRESS|REPEAT)
}

func (self *KeyEvent) MatchesCaseSensitiveTextOrKey(spec string) bool {
	if self.MatchesParsedShortcut(ParseShortcut(spec), PRESS|REPEAT) {
		return true
	}
	return self.Text == spec
}

func (self *KeyEvent) MatchesCaseInsensitiveTextOrKey(spec string) bool {
	if self.MatchesParsedShortcut(ParseShortcut(spec), PRESS|REPEAT) {
		return true
	}
	return strings.ToLower(self.Text) == strings.ToLower(spec)
}

func (self *KeyEvent) MatchesRelease(spec string) bool {
	return self.MatchesParsedShortcut(ParseShortcut(spec), RELEASE)
}

func (self *KeyEvent) AsCSI() string {
	key := csi_number_for_name(self.Key)
	shifted_key := csi_number_for_name(self.ShiftedKey)
	alternate_key := csi_number_for_name(self.AlternateKey)
	trailer, found := csi_number_to_letter_trailer_map[key]
	if !found {
		trailer = "u"
	}
	if self.Key == "ENTER" {
		trailer = "u"
	}
	if trailer != "u" {
		key = 1
	}
	ans := strings.Builder{}
	ans.Grow(32)
	ans.WriteString("\033[")
	if key != 1 || self.Mods != 0 || shifted_key != 0 || alternate_key != 0 || self.Text != "" {
		ans.WriteString(fmt.Sprint(key))
	}
	if shifted_key != 0 || alternate_key != 0 {
		ans.WriteString(":")
		if shifted_key != 0 {
			ans.WriteString(fmt.Sprint(shifted_key))
		}
		if alternate_key != 0 {
			ans.WriteString(fmt.Sprint(":", alternate_key))
		}
	}
	action := 1
	switch self.Type {
	case REPEAT:
		action = 2
	case RELEASE:
		action = 3
	}
	if self.Mods != 0 || action > 1 || self.Text != "" {
		m := uint(self.Mods)
		if action > 1 || m != 0 {
			ans.WriteString(fmt.Sprintf(";%d", m+1))
			if action > 1 {
				ans.WriteString(fmt.Sprintf(":%d", action))
			}
		} else if self.Text != "" {
			ans.WriteString(";")
		}
	}
	if self.Text != "" {
		runes := []rune(self.Text)
		codes := make([]string, len(runes))
		for i, r := range runes {
			codes[i] = strconv.Itoa(int(r))
		}
		ans.WriteString(";")
		ans.WriteString(strings.Join(codes, ":"))
	}
	fn, found := name_to_functional_number_map[self.Key]
	if found && tilde_trailers[fn] {
		trailer = "~"
	}
	ans.WriteString(trailer)
	return ans.String()
}

func csi_number_for_name(key_name string) int {
	if key_name == "" {
		return 0
	}
	if key_name == "F3" || key_name == "ENTER" {
		return 13
	}
	fn, ok := name_to_functional_number_map[key_name]
	if !ok {
		return int(rune(key_name[0]))
	}
	ans, ok := functional_to_csi_number_map[fn]
	if ok {
		return ans
	}
	return fn
}

func init() {
	name_to_functional_number_map = make(map[string]int, len(functional_key_number_to_name_map))
	for k, v := range functional_key_number_to_name_map {
		name_to_functional_number_map[v] = k
	}
	functional_to_csi_number_map = make(map[int]int, len(csi_number_to_functional_number_map))
	for k, v := range csi_number_to_functional_number_map {
		functional_to_csi_number_map[v] = k
	}
	csi_number_to_letter_trailer_map = make(map[int]string, len(letter_trailer_to_csi_number_map))
	for k, v := range letter_trailer_to_csi_number_map {
		csi_number_to_letter_trailer_map[v] = k
	}
}