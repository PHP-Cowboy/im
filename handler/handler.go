package handler

import (
	"im/ws"
)

func RegisterHandlers(s *ws.Server) {
	s.AddRoutes([]ws.Route{
		{
			Method:  "login", //登录
			Handler: Login(),
		},
		{
			Method:  "enterRoom", //进入直播间
			Handler: EnterRoom(),
		},
		{
			Method:  "leaveRoom", //离开直播间
			Handler: LeaveRoom(),
		},
		{
			Method:  "applyVoiceChat", //申请加入语聊
			Handler: ApplyVoiceChat(),
		},
		{
			Method:  "privateChat", //私聊
			Handler: PrivateChat(),
		},
		{
			Method:  "liveRoomChat", //直播间聊天
			Handler: LiveRoomChat(),
		},
		{
			Method:  "broadcast", //广播
			Handler: Broadcast(),
		},
		{
			Method:  "sendGift", //送礼
			Handler: SendGift(),
		},
		{
			Method:  "otherToUser", //其他消息，只做转发
			Handler: OtherToUser(),
		},
		{
			Method:  "otherToAll", //其他消息，只做转发
			Handler: OtherToAll(),
		},
	})
}
