package handler

import (
	"im/global"
	"im/model/user"
	"im/ws"
)

type ApplyVoiceChatReq struct {
	RoomId string `json:"roomId"`
}

type ApplyVoiceChatRsp struct {
	Uid      int    `json:"uid"`
	NickName string `json:"nickName"`
}

func ApplyVoiceChat() ws.HandleFunc {
	return func(s *ws.Server, conn *ws.Conn, msg *ws.Message) {
		var (
			req ApplyVoiceChatReq
			rsp ApplyVoiceChatRsp
		)

		userId := conn.GetUserId()

		u, err := GetUserInfoById(userId)

		if err != nil {
			global.Logger["err"].Errorf("GetUserInfoById failed,err:%v", err.Error())
			return
		}

		rsp.Uid = u.ID
		rsp.NickName = u.NickName

		//查询主播id 发送消息给主播 anchorId
		anchorId, err := GetAnchorIdByShowId(req.RoomId)

		if err != nil {
			global.Logger["err"].Errorf("GetAnchorIdByShowId failed,err:%v", err.Error())
			return
		}

		anchorConn, ok := s.UserToConnMp[anchorId]

		if !ok {
			//主播未找到
			global.Logger["err"].Errorf("主播未找到,roomId:%v,申请userId:%v", req.RoomId, userId)
			return
		}

		err = s.Send(ws.NewMessage(0, 0, rsp), anchorConn)

		if err != nil {
			global.Logger["err"].Error(err.Error())
			return
		}
	}
}

func GetAnchorIdByShowId(showId string) (int, error) {
	obj := new(user.UserLive)

	db := global.DB

	data, err := obj.GetOneByShowId(db, showId)

	if err != nil {
		global.Logger["err"].Errorf("obj.GetOneByShowId failed,err:%v", err.Error())
		return 0, err
	}

	return data.Uid, nil
}
