基于Monibuca的直播推流服务#

需要同时启动monibuca/main.go与live_server.go

接口：

/createLive

    POST
    Form:
        name: 直播名称
        poster: 直播封面

    用name和post创建一个直播，并分配一个id，每个直播的名字不能重复
    如果名字重复，返回406

/getAllLive

    GET
    
    获取当前存在的所有直播信息

/getLiveName
    
    GET
    Paramas:
        stream_id: 流标识

    根据流标识查询直播信息
    如果不存在，返回404

/pushVideoToStream

    POST
    Form:
        stream_id: 流标识
        path: 视频地址
    
    根据视频地址和流标识，将视频推流到Monibuca

/pushStreamToRtmp
    
    POST
    Form:
        stream_id: 流标识
        rtmp_addr: rtmp地址
    
    根据将Monibuca上的视频流推送到rtmp地址
    如果不存在，返回404