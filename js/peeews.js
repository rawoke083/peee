var PeeeWS = {
	wsocket: 0,
	init: function(obj) {
		this.wsocket = new WebSocket("ws://" + AppConfig.HOST_NAME + ":8080/rest/cmd");
		
			this.wsocket.onmessage=function(e){
				return obj.onWSMessage(e);

			}
			/*
		this.wsocket.onmessage=function(e){
			return PCanvas.onWSMessage(e);

			}*/
       /*
        g_wsocket.onmessage = function(e) {
			g_WData = JSON.parse(e.data);
        }*/
       // alert(AppConfig.HOST_NAME);
	},
	button_a:function(){
		alert("fff");
		var cmdObj = {
			GameId: parseInt(PControl.GameData.GameId),
			Cmd: 10,
			Cmd2: 0,
			Rkey:PControl.GameData.RKey,
			timestamp: Date.now()
		};
		
		this.send_cmd(cmdObj);
	},
	
	button_b:function(){
		var cmdObj = {
			GameId: parseInt(PControl.GameData.GameId),
			Cmd: 12,
			Cmd2: 0,
			Rkey:PControl.GameData.RKey,
			timestamp: Date.now()
		};
		
		this.send_cmd(cmdObj);
	},
	
	send_cmd: function(cmdObj) {
			if (this.wsocket.readyState != 1) {
				console.log("NOT SENDING");
				return;
			}
			console.log("PWSENDING");	
			PeeeWS.wsocket.send(JSON.stringify(cmdObj));
		} //end send_cmd
		
	
};
