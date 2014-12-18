/******************************
 * Config
 ******************************/
var AppConfig = {
	//HOST_NAME: "152.111.241.186",
	HOST_NAME: "10.0.0.16",
	
	getParam: function(pindex) {
		var params = [];
		window.location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(str, key, value) {
			params[key] = value;
		});
		return params[pindex];
	}
	
};

var AppGlobal = {
	GameId:0,
	RKey:0,
	CType:0
	
};
