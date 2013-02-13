//piou piou

function actionmenu(page) {
    $('.content').load('ajx/' + page);
}

$(".sidebar li").click(function() {
    var name = $(this).find("a").attr("href").substring(1);
    
    $(".sidebar li").removeClass("active");
    if(name == "irc") {
        load_irc();
    } else {
        actionmenu(name);
    }
    $(this).addClass("active");
});