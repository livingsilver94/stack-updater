// Code generated by "enumer -transform=lower -type=SupportedStack"; DO NOT EDIT.

package stack

import (
	"fmt"
)

const _SupportedStackName = "kdemate"

var _SupportedStackIndex = [...]uint8{0, 3, 7}

const _SupportedStackLowerName = "kdemate"

func (i SupportedStack) String() string {
	if i < 0 || i >= SupportedStack(len(_SupportedStackIndex)-1) {
		return fmt.Sprintf("SupportedStack(%d)", i)
	}
	return _SupportedStackName[_SupportedStackIndex[i]:_SupportedStackIndex[i+1]]
}

var _SupportedStackValues = []SupportedStack{0, 1}

var _SupportedStackNameToValueMap = map[string]SupportedStack{
	_SupportedStackName[0:3]:      0,
	_SupportedStackLowerName[0:3]: 0,
	_SupportedStackName[3:7]:      1,
	_SupportedStackLowerName[3:7]: 1,
}

var _SupportedStackNames = []string{
	_SupportedStackName[0:3],
	_SupportedStackName[3:7],
}

// SupportedStackString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SupportedStackString(s string) (SupportedStack, error) {
	if val, ok := _SupportedStackNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to SupportedStack values", s)
}

// SupportedStackValues returns all values of the enum
func SupportedStackValues() []SupportedStack {
	return _SupportedStackValues
}

// SupportedStackStrings returns a slice of all String values of the enum
func SupportedStackStrings() []string {
	strs := make([]string, len(_SupportedStackNames))
	copy(strs, _SupportedStackNames)
	return strs
}

// IsASupportedStack returns "true" if the value is listed in the enum definition. "false" otherwise
func (i SupportedStack) IsASupportedStack() bool {
	for _, v := range _SupportedStackValues {
		if i == v {
			return true
		}
	}
	return false
}
