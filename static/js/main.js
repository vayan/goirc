$(document).ready(function() {
    var hash = window.location.hash.substring(1);
    if(hash != "")  $('.content').load('ajx/' + hash);
    else $('.content').load('ajx/home');

    //FIXME : click not working
     $(".item-menu-irc").click(function() {
        console.log("test");
        $(".list").css("top", "160px");
        $(".bufferchan").css("top", "160px");
    });

});