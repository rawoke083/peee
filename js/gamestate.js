define(["peeerest"], function (peeerest) {
	
	prest = new peeerest();
	
	
    var returnedModule = function () {
        var name = peeeconfig.HOST_NAME;
        this.getName = function () {
            
            
            return name;
            
        }
    };

    return returnedModule; 
});
