var PControl = {
	Graphics: {
		imgP1: 0,
		imgP2: 0,
		imgBall: 0,
		canvas: 0,
		context: 0,
	},
	GameData: {
		GameId: 0,
		RKey: 0,
		LastSendTick: 0,
		GameWorld: 0,
		PlayerIndex: -1,
		Ready: 0,
		LastMsgId:0,
	},
	
	getWorldData: function() {
		var cmdObj = {
			GameId: parseInt(this.GameData.GameId),
			Cmd: 20,
			Cmd2: 0,
			Rkey: "",
			timestamp: Date.now()
		};
		console.log("sending" + this.GameData.GameId);
		PeeeWS.send_cmd(cmdObj);
	},
	
	// a very basic method
	init: function() {
			this.GameData.GameId = AppConfig.getParam("gameid");
			this.GameData.RKey = AppConfig.getParam("rkey");
				
				
			if (this.GameData.RKey == 0 || this.GameData.RKey == undefined ) {
				this.GameData.RKey = Math.floor((Math.random() * 10000000) + 1) + Math.floor((Math.random() * 10000000) + 1);
				 window.location = window.location + "&rkey=" + this.GameData.RKey;
				
			}else{
				PeeeREST.JoinGame(PControl.GameData.GameId,PControl.GameData.RKey);
				PeeeWS.wsocket.onmessage = function(e) {this.onWSMessage(e)}
			}
			
		
			
			
		
		}, //end init
		
	 onWSMessage:function(e){
	
		this.GameData.GameWorld = JSON.parse(e.data);
	

		pindex = this.GameData.PlayerIndex;

		pcmd = this.GameData.GameWorld.Players[pindex].PMsg;
		
		
		if(pcmd > 0 && this.GameData.GameWorld.Players[pindex].PMsgId != this.GameData.LastMsgId ){
			
			if(pcmd==100){
				  navigator.vibrate([200,100,200]);
				  
			}else if(pcmd==101){
				  navigator.vibrate(100);
			}
		
			
			this.GameData.LastMsgId = this.GameData.GameWorld.Players[pindex].PMsgId;
		}
		
  }
};
