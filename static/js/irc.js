function add_new_buffer(buffer) {
    var id = buffer[1];

    $('.listbuffer').append('<li><a href="#' + id + '" data-toggle="tab">' + buffer[2] + '</a></li>');
    $('.contentbuffer').append('<div class="tab-pane bufferchan" id="' + id + '"><table class="table table-striped allmsg"></table></div>');
    $.post("/ajx/backlog", { idbuffer: id })
        .done (function(data) {
            $('.contentbuffer #' + id + ' .allmsg').append(data);
        });
}

function new_message(id_buffer, nick, msg) {
    if(msg.charAt(0) == '/') return;
    $('.contentbuffer #' + id_buffer + ' .allmsg').append('<tr class="msg"><td class="pseudo">' + nick + '</td><td class="message">' + msg + '</td><td class="time">' + get_timestamp_now() + '</td></tr>');
    $('#' + id_buffer).scrollTop($('#' + id_buffer)[0].scrollHeight);
}

function get_timestamp_now() {
    var d = new Date();
    var timestamp = d.getHours() + ":" + d.getMinutes();
    return timestamp;
}

function parse_irc(msg) {
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
}

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