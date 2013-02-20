//piou piou

function update_active_sidebar(page) {
    $('.switch-userlist').hide();
    $(".sidebar #menu  li").removeClass("active");
    $(".sidebar #menu #"+page).addClass("active");
}

function ChangePage(page) {
    switch(page) {
        case "irc":
            update_active_sidebar(page);
            load_irc();
        break;
        case "home" :
            update_active_sidebar(page);
            $('.content').load('ajx/home');
        break;
         case "register" :
            update_active_sidebar(page);
            $('.content').load('ajx/register');
        break;
         case "login" :
            update_active_sidebar(page);
            $('.content').load('ajx/login');
        break;
        default:
        break;
    }
}

$(".sidebar #menu li").click(function() {
    var name = $(this).find("a").attr("href").substring(1);
    
    if(name == "irc") {
        load_irc();
    } else {
        ChangePage(name);
    }
});