package video

import "douyin/app/service/video/internal/sys"

const (
	RdbKey = sys.SysName + ":"

	RdbKeyFavorite    = RdbKey + "favorite:"
	RdbKeyFavoriteCnt = RdbKey + "favorite_cnt:"
	RdbKeyCommentCnt  = RdbKey + "comment_cnt:"
)
