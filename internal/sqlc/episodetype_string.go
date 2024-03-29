// Code generated by "stringer -type=EpisodeType"; DO NOT EDIT.

package sqlc

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EpRegular-1]
	_ = x[EpSpecial-2]
	_ = x[EpCredit-3]
	_ = x[EpTrailer-4]
	_ = x[EpParody-5]
	_ = x[EpOther-6]
}

const _EpisodeType_name = "EpRegularEpSpecialEpCreditEpTrailerEpParodyEpOther"

var _EpisodeType_index = [...]uint8{0, 9, 18, 26, 35, 43, 50}

func (i EpisodeType) String() string {
	i -= 1
	if i < 0 || i >= EpisodeType(len(_EpisodeType_index)-1) {
		return "EpisodeType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _EpisodeType_name[_EpisodeType_index[i]:_EpisodeType_index[i+1]]
}
