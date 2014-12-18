/******************************
 * Config
 ******************************/


var gHOST_NAME = "152.111.241.186";


var WS_MSG_GET_GAME_DATA=20;


//Generic
function ws_cmd_send(cmd,cmd2){
	
	
	if( g_wsocket.readyState ==1  && (g_GLAST_CMD_TIME == 0  || (Date.now()- GLAST_CMD_TIME) > 1))
	{

		var cmdObj = {
			GameId:parseInt(g_GameId),
			Cmd:cmd,
			Cmd2:cmd2,
			Rkey:"",
			timestamp: Date.now()
		};
  
		g_wsocket.send(JSON.stringify(cmdObj));
		GLAST_CMD_TIME =   Date.now();
		
	}//end timecheck 
	
}//end ws_cmd


//Main

function main_new_game() {
    $.ajax({
        type: "POST",
        url: "http://"+gHOST_NAME+"/rest/game",
        data: "koos=1",
        success: function(data) {
                main_list_games()
            } //end succues

    }); //end ajax;

}


function main_list_games() {

    $("#glist").html("");

    $.ajax({
        crossDomain: true,

        url: "http://"+gHOST_NAME+"/rest/game",
        dataType: 'json',
        success: function(data) {
			
            $.each(data, function(index) {


                var items = [];
                if (data[index].Id > 0) {

                    items.push('<li> <h4><a href="/canvas.html?gameid=' + data[index].Id + '"> As TV (' + data[index].Id + ')| </a> | <a href="/pcontrol.html?gameid=' + data[index].Id + '"> As RemoteControl -  Players(' + data[index].PCount + ')</a></h4></li>');

                }
                $('#glist').append(items.join(''));

            });//end foreach

        }//end success
    });//end ajax-call


}//end main_list_games




