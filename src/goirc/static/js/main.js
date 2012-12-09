$(document).ready(function() {
	var hash = window.location.hash.substring(1);
	if (hash != "")
		$('#content').load('ajx/'+hash);
});