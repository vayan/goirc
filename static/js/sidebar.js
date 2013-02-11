//piou piou

function actionmenu(page) {
    $('.content').load('ajx/' + page);
}

function load_irc() {
    var irc = $("#clientirc").html();

    $('.content').html(irc);
        $(".formirc input").keyup(function(event) {
            if(event.keyCode == 13) {
                $(".formirc button").click();
            }
        });
        $(".item-menu-irc").click(function() {
            console.log($(this).html());
            $(".list").css("top", "100px");
            $(".bufferchan").css("top", "100px");
            $(".menu-settings").show();
        });
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