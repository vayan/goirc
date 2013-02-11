$(document).ready(function() {
    var hash = window.location.hash.substring(1);
    if (hash != "") {
    	if (hash == "irc") load_irc();
    	else $('.content').load('ajx/' + hash);
    }  
    else $('.content').load('ajx/home');
});