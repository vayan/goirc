//piou piou

function actionmenu(page) {
    $('.content').load('ajx/' + page);
}


$(".sidebar li").click(function() {
    var name = $(this).find("a").attr("href").substring(1);
    var irc = $("#clientirc").html();
    $(".sidebar li").removeClass("active");
    if(name == "irc") {
        $('.content').html(irc);
        // $("tr:even").css("background-color", "#f7f7f9");
        // $("tr:odd").css("background-color", "#fff");
        //$(".bufferchan").scrollbars();
        $(".formirc input").keyup(function(event) {
            if(event.keyCode == 13) {
                $(".formirc button").click();
            }
        });
    } else {
        actionmenu(name);
    }
    $(this).addClass("active");
});