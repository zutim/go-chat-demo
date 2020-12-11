var webServer = {
    config:{
        url : "localhost:8080/ws",
    },
    ws:'',
    msgSender:function (x) {
        this.ws.send(JSON.stringify(x));
    },
    encode:function (x) {
        return JSON.stringify(x);
    },
    decode:function (x) {
        return JSON.parse(x);
    }
};
var netManager = {
    msg:{
        op: "",
        args: "",
        msg: "",
        msgType: "",
        flagId: 0,
    },
    msgHandler:function(x){
        let m = webServer.decode(x);
        switch (m.op) {
            case "connErr":
                netManager.msg.op = "close";
                webServer.msgSender(netManager.msg)
                break;
            case "connSuccess":

                var userList = webServer.decode(m.msg);

                var html="";
                for(var i in userList){
                    if(userList[i] != null){
                        var user = userList[i];
                        html+="<li data-userid='"+user.id+"' data-username='"+user.username+"' onclick='app.chatTo(this)'>\n" +
                            "       <img src='"+user.headimg+"'>" +
                            "       <span>"+user.username+"</span>\n"+
                            "                    <div class='right'>" ;
                        if(user.unread !=0 ){
                            html+="<div class='unreadnum'>"+user.unread+"</div>";
                        }

                        if(user.online == 1){
                            html+="<div class='line blue'>" +
                                "<div class=\"status-point\" style=\" background-color:#67C23A\" /> 在线</div></div>";
                        }else{
                            html+="<div class='line block'>" +
                                "<div class=\"status-point\" style=\" background-color:#b6b8a9\" /> 离线</div></div>";
                        }
                        html+=  "                    </div>\n" +
                                "                </li>";
                    }
                }
                $(".chat ul").empty()
                $(".chat ul").html(html)
                break;
        }
        console.log(m)
    }
}