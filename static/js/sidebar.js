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
            $(".item-menu-irc").removeClass("selected");
            if ($(".menu-settings").css("display") == "block") {
                $(".list").css("top", "60px");
                $(".bufferchan").css("top", "60px");
                $(".menu-settings").hide();
            } else {
                $(".list").css("top", "130px");
                $(".bufferchan").css("top", "130px");
                $('.menu-settings').load('ajx/set-' + $(this).html().toLowerCase());
                $(".menu-settings").show();
                $(this).addClass("selected");
            }   
        });

        $(".formirc button").click(function() {
            if($(".formirc input").val() != '') {
                //console.log($(".active a").html());
                var buffer_id = $(".main-irc .active a").attr('href').substring(1);
                var txt = $(".formirc input").val();
                var msg = buffer_id + "]" + txt;

                console.log(msg);
                ws.send(msg);
                new_message(buffer_id, "me", txt);
            }
            $(".formirc input").val("").focus();
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