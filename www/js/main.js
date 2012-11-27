$(document).ready(function() {
	var hash = window.location.hash.substring(1);
	$('#content').load('js/ajx/'+hash+'.html');
});