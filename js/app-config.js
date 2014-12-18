/******************************
 * Config
 ******************************/
require.config({
	paths: {
		'jQuery': 'jquery.x.min',
		'peeerest': 'peeerest'
	},
	shim: {
		'jQuery': {
			exports: '$'
		}
	}
});

define({
   HOST_NAME:"10.0.0.16"
});
