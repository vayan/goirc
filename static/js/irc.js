var load_irc = function () {
    var irc = $("#clientirc").html();

    $('.content').html(irc);
    $(".formirc input").keyup(function(event) {
        if(event.keyCode == 13) {
            send_message();
        }
    });

    $(".item-menu-irc").click(function() {
            $(".item-menu-irc").removeClass("selected");
            if ($(".menu-settings").css("display") == "block") {
                console.log("hide");
                $(".list").css("top", "60px");
                $(".bufferchan").css("top", "60px");
                $(".menu-settings").hide();
            } else {
                console.log("show");
                $(".list").css("top", "130px");
                $(".bufferchan").css("top", "130px");
                $('.menu-settings').load('ajx/set-' + $(this).html().toLowerCase());
                $(".menu-settings").show();
                $(this).addClass("selected");
            }
        });
    $('.switch-userlist').show();
    $('#userlisttab').click ();
};

//TODO : Optimize this call less
var aff_user_list = function (id) {
    $.post("/ajx/userslist", {
        channel: id
    }).done(function(data) {
        $('#userlist').html(data);
    });
};

var send_new_co_serv = function () {
    var msg = "log]/connect " + $("#adressserv").val() + ":" + $("#portserv").val();
    ws.send(msg);
    new_message("log", "log", "Connecting to "+$("#adressserv").val()+"...");
    $(".menu-irc > .selected").click();
};

var send_new_join_chan = function () {
    var msg = $("#idnetwork").val()+"]/join "+$("#adresschan").val();
    ws.send(msg);
    new_message("log", "log", "Joining "+$("#adresschan").val()+"...");
    $(".menu-irc > .selected").click();
};

var send_message = function () {
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
};

var add_new_buffer = function (buffer) {
    var id = buffer[1];

    if (buffer[2][0] != '#') {
        new_message("log", "log", "Connected to "+buffer[2]+"!");
    } else {
        new_message("log", "log", "Joined "+buffer[2]+"!");
    }

    $('.listbuffer').append('<li onclick="aff_user_list('+id+')"><a href="#' + id + '" data-toggle="tab">' + buffer[2] + '</a></li>');
    $('.contentbuffer').append('<div class="tab-pane bufferchan" id="' + id + '"><table class="table table-striped allmsg"></table></div>');
    $.post("/ajx/backlog", {
        idbuffer: id
    }).done(function(data) {
        $('.contentbuffer #' + id + ' .allmsg').append(data);
    });
};

var new_message = function (id_buffer, nick, msg) {
    if(msg.charAt(0) == '/') return;
    $('.contentbuffer #' + id_buffer + ' .allmsg').append('<tr class="msg"><td class="pseudo">' + nick + '</td><td class="message">' + msg + '</td><td class="time">' + get_timestamp_now() + '</td></tr>');
    $('#' + id_buffer).scrollTop($('#' + id_buffer)[0].scrollHeight);
};

var get_timestamp_now = function () {
    var d = new Date();
    var timestamp = d.getHours() + ":" + d.getMinutes();
    return timestamp;
};

var parse_irc = function (msg) {
// TODO : check le ] dans le message pour eviter split useless
var buff = msg.split(']');
switch(buff[0]) {
case "buffer":
console.log("new buffer " + buff[1]);
add_new_buffer(buff);
break;
default:
new_message(buff[0], buff[1], buff[2]);
break;
}
};