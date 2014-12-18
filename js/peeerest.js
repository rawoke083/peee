var PeeeREST = {
	gameListElemt: 0,
	init: function(htmlElemGameList) {
		this.gameListElemt = htmlElemGameList;
	},
	
	
	// a very basic method
	GetGameList: function() {
		$.ajax({
			//crossDomain: true,
			url: "http://" + AppConfig.HOST_NAME + "/rest/game",
			dataType: 'json',
			async: false,
			success: function(data) {
					PeeeREST.gameListElemt.html("");
					$.each(data, function(index) {
						var items = [];
						if (data[index].Id > 0) {
							items.push('<li> <h4><a href="/canvas.html?gameid=' + data[index].Id + '"> As TV (' + data[index].Id + ')| </a> | <a href="/pcontrol.html?gameid=' + data[index].Id + '"> As RemoteControl -  Players(' + data[index].PCount + ')</a></h4></li>');
						}
						PeeeREST.gameListElemt.append(items.join(''));
					}); //end foreach
				} //end success
		}); //end ajax-call
	},
	// output a value based on the current configuration
	
	
	NewGame: function() {
		$.ajax({
			type: "POST",
			url: "http://" + AppConfig.HOST_NAME + "/rest/game",
			data: "koos=1",
			success: function(data) {
					PeeeREST.GetGameList();
				} //end succues
		}); //end ajax;
	}, //end new game
	
	
	JoinGame: function(GameId, RKey) {
			$.ajax({
				type: "POST",
				url: "http://" + AppConfig.HOST_NAME + "/rest/gamejoin?gameid=" + GameId + "&rkey=" + RKey,
				async: false,
				dataType: 'json',
				success: function(data) {
						if (data == 0 || data == undefined) {
							return;
						}
						if (data.Players[0].RKey == PControl.GameData.RKey) {
							PControl.GameData.PlayerIndex = 0;
						} else if (data.Players[1].RKey == PControl.GameData.RKey) {
							PControl.GameData.PlayerIndex = 1;
						}
						$('#pstatus').text("Connected - Game (" + GameId + ")");
						$('#pindex').text("Player (" + (PControl.GameData.PlayerIndex + 1) + ")");
						if (PControl.GameData.PlayerIndex == 0) {
							$('#pindex2').html("<img width='250px;' src='img/p1.png'/>");
						} else {
							$('#pindex2').html("<img  width='250px;' src='img/p2.png'/>");
						}
						navigator.vibrate([200, 300, 200]);
						PControl.GameData.Ready = 1;
					} //end succues
			}); //end ajax;
		} //end join
};
