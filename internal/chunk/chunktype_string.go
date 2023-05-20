// Code generated by "stringer -type=ChunkType -trimprefix=C"; DO NOT EDIT.

package chunk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CMessages-0]
	_ = x[CThreadMessages-1]
	_ = x[CFiles-2]
	_ = x[CUsers-3]
	_ = x[CChannels-4]
	_ = x[CChannelInfo-5]
	_ = x[CWorkspaceInfo-6]
	_ = x[CStarredItems-7]
	_ = x[CBookmarks-8]
}

const _ChunkType_name = "MessagesThreadMessagesFilesUsersChannelsChannelInfoWorkspaceInfoStarredItemsBookmarks"

var _ChunkType_index = [...]uint8{0, 8, 22, 27, 32, 40, 51, 64, 76, 85}

func (i ChunkType) String() string {
	if i >= ChunkType(len(_ChunkType_index)-1) {
		return "ChunkType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ChunkType_name[_ChunkType_index[i]:_ChunkType_index[i+1]]
}
