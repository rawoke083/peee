
var PCanvas = {
 
 Graphics:{
 imgP1:0, 
 imgP2:0,
 imgBall:0,
 canvas:0,
 context:0,
},

 GameData:{
	GameId:0,
	LastSendTick:0,
	GameWorld:0,
 },
 
 
getWorldData:function(){
	
		var cmdObj = {
			GameId:parseInt(this.GameData.GameId),
			Cmd:20,
			Cmd2:0,
			Rkey:"",
			timestamp: Date.now()
		};
	
	
	PeeeWS.send_cmd(cmdObj);
},

drawWorld:function(){

	if(this.GameData.GameWorld == 0 ){
	
		return;
	}
	
	
	
	// buffer canvas
var canvas2 = document.createElement('canvas');
canvas2.width = 1300;
canvas2.height = 600;
var context2 = canvas2.getContext('2d');
	
	context2.save();
//	context2.clearRect(0, 0, canvas.width, canvas.height);
	context2.fillStyle = "#009933";
	context2.fillRect(0,0,1300,600);
	
	//context2.fillStyle = "#fff";
	context2.fillStyle = "#9FF781";
	
	context2.fillRect(660,0,20,600);
	
GWData = this.GameData.GameWorld;
	
	
	context2.drawImage(this.Graphics.imgP1,GWData.Players[0].XPos, GWData.Players[0].YPos);
	context2.drawImage(this.Graphics.imgP2,GWData.Players[1].XPos, GWData.Players[1].YPos);
	context2.drawImage(this.Graphics.imgBall,GWData.Ball.XPos,GWData.Ball.YPos);
	
	if(GWData.TextMsg) {
		
		context2.fillStyle = "#222222";
	
		context2.fillRect(360,235,600,90);
		
		context2.fillStyle = "#efefef";
		context2.font="40px Sans";
		context2.fillText(GWData.TextMsg,400,300);
	}
	
	 context2.restore();
	
	
	
	this.Graphics.context.save();
	this.Graphics.context.clearRect(0, 0, this.Graphics.canvas.width, this.Graphics.canvas.height);
	this.Graphics.context.drawImage(canvas2, 0, 0);

    
	 this.Graphics.context.restore();
	 $("#p1score").text(GWData.Players[0].Score);
	 $("#p2score").text(GWData.Players[1].Score);
	 


},//end drawworld

GameLoop:function(){

	this.getWorldData();
	this.drawWorld();
	
},

 
  // a very basic method
  init: function () {
   
   PCanvas.Graphics.canvas =  document.getElementById('myCanvas');
		
    
	
	PCanvas.Graphics.context = PCanvas.Graphics.canvas.getContext('2d');
	
	PCanvas.Graphics.imgP1 =  new Image();
	PCanvas.Graphics.imgP2 = new Image();
	PCanvas.Graphics.imgBall = new Image();
	
	PCanvas.Graphics.imgP1.src = 'img/p1.png';
	PCanvas.Graphics.imgP2.src = 'img/p2.png';
	PCanvas.Graphics.imgBall.src= 'img/ball.png';
	
	PCanvas.Graphics.imgP1.onload = function() {
		PCanvas.Graphics.context.drawImage(PCanvas.Graphics.imgP1,0, 0);
	};

	PCanvas.Graphics.imgP2.onload = function() {
		PCanvas.Graphics.context.drawImage(PCanvas.Graphics.imgP2,0, 0);
	};
	
	PCanvas.Graphics.imgBall.onload = function() {
		PCanvas.Graphics.context.drawImage(PCanvas.Graphics.imgBall,0, 0);
		PCanvas.Graphics.context.drawImage(PCanvas.Graphics.imgBall,0, 0);
	};

   
	PCanvas.GameData.GameId = AppConfig.getParam('gameid');
	
//	PeeeWS.wsocket.onmessage = function(e) {this.onWSMessage(e);}
	
	
 
  },//end init
 
  onWSMessage:function(e){
	this.GameData.GameWorld = JSON.parse(e.data);
	console.log("RX"+e.data);
  }
  
  
};

/*
define(["peeeconfig", "jQuery"], function(peeeconfig, $) {

	var returnedModule = function() {
	
		this.getGameList = function(htmlElem) {
				$.ajax({
					//crossDomain: true,
					url: "http://" + peeeconfig.HOST_NAME + "/rest/game",
					dataType: 'json',
					async: false,
					success: function(data) {
							htmlElem.html("");
							
							$.each(data, function(index) {
								var items = [];
								if (data[index].Id > 0) {
									items.push('<li> <h4><a href="/canvas.html?gameid=' + data[index].Id + '"> As TV (' + data[index].Id + ')| </a> | <a href="/pcontrol.html?gameid=' + data[index].Id + '"> As RemoteControl -  Players(' + data[index].PCount + ')</a></h4></li>');
								}
								htmlElem.append(items.join(''));
							}); //end foreach
						} //end success
				}); //end ajax-call
				
			} //end get game list
	};
	return returnedModule;
});
*/
