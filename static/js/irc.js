var ws;
var host = window.location.hostname;
var buffers;
var usersettings = null;

if (ws != null) {
    ws.close();
    ws = null;
}

var get_user_pref = function() {
    $.post("/ajx/getsettings").done(function(data) {
        usersettings = JSON.parse(data);
    });
};

var open_settings = function() {
    $("#settings-notify").bootstrapSwitch("setState", usersettings.Notify);
    $("#settings-savesession").bootstrapSwitch("setState", usersettings.Save_session);

    $('#form-settings').submit(function() {
        set_user_pref();
        return false;
    });
};

var set_user_pref = function() {
    console.log("gere");
    var notify = $("#settings-notify").bootstrapSwitch('status');
    var Save_session = $("#settings-savesession").bootstrapSwitch('status');

    $.post("/ajx/settings", {
        Notify: notify,
        Save_Session: Save_session
    }).done(function(data) {
        get_user_pref();
    });
};


if (usersettings === null) {
    get_user_pref();
}

ws = new WebSocket("ws://" + host + ":1112/ws");


ws.onopen = function() {
    console.log("open ws");
    if ($("#yuid").val() != "") {
        ws.send("co]" + $("#yuid").val());
    }
};

ws.onmessage = function(e) {
    console.log("receive : " + e.data);
    parse_irc(e.data);

};

ws.onclose = function(e) {
    console.log("close ws");
};


var SplitN = function(s, sep, n) {
    arr = s.split(sep),
    result = arr.splice(0, n);
    result.push(arr.join(sep));
    return result;
};


var update_active_sidebar = function(page) {
    $('.switch-userlist').hide();
    $(".sidebar #menu  li").removeClass("active");
    $(".sidebar #menu #" + page).addClass("active");
};



var ChangePage = function(page) {
    switch (page) {
        case "irc":
            update_active_sidebar(page);
            load_irc();
            break;
        case "home":
            update_active_sidebar(page);
            $('.content').load('ajx/home');
            break;
        case "register":
            update_active_sidebar(page);
            $('.content').load('ajx/register');
            break;
        case "login":
            update_active_sidebar(page);
            $('.content').load('ajx/login');
            break;
        case "settings":
            update_active_sidebar(page);
            $('.content').load('ajx/settings', function() {
                open_settings();
            });
            break;
        default:
            break;
    }
};


var load_irc = function() {
    var irc = $("#clientirc").html();

    $('.content').html(irc);
    $(".inputtextirc").keyup(function(event) {
        if (event.keyCode == 13) {
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
    $('#userlisttab').click();
};

var switch_buffer = function(id) {
    aff_user_list(id);

    $(".inputpseudo").val($("#" + id + " .current-nick").val());
};

var update_user_list = function(id) {
    var newcss;

    $.post("/ajx/userslist", {
        channel: id
    }).done(function(data) {
        $("#userlist-buffer" + id).html($("#userlist-buffer0").html());
        jsonres = JSON.parse(data).UserList;
        $("#userlist-buffer" + id + " .userlist-style").html("<style></style>");
        for (var i = 0; i < jsonres.length; i++) {
            newcss += add_user_list(jsonres[i].Nick, jsonres[i].Color, jsonres[i].NickClean, id);
        }
        $("#userlist-buffer" + id + " .userlist-style style").html(newcss);
        $(".userlist-buffer").hide();
        //TODO : rendre visible uniquement la list active
        $("#userlist-buffer" + id).show();
    });
};

var add_user_list = function(name, color, id, buffer) {
    var newrulecss = ".nick-" + id + " {  color : " + color + " ; } ";
    $("#userlist-buffer" + buffer + ' .userlist-user').append("<li class='nick-" + id + "'>" + name + "</li>");
    return newrulecss;
};

var aff_user_list = function(id) {

    if ($("#userlist-buffer" + id).length <= 0) {
        var newcss;
        $.post("/ajx/userslist", {
            channel: id
        }).done(function(data) {
            $("#userlist").append("<span id='userlist-buffer" + id + "' class='userlist-buffer' >" + $("#userlist-buffer0").html() + "</span>");
            jsonres = JSON.parse(data).UserList;
            // $('#userlist-user').html("");
            $("#userlist-buffer" + id + " .userlist-style").html("<style></style>");
            for (var i = 0; i < jsonres.length; i++) {
                newcss += add_user_list(jsonres[i].Nick, jsonres[i].Color, jsonres[i].NickClean, id);
            }
            $("#userlist-buffer" + id + " .userlist-style style").html(newcss);
            $(".userlist-buffer").hide();
            $("#userlist-buffer" + id).show();
        });
    } else {
        $(".userlist-buffer").hide();
        $("#userlist-buffer" + id).show();
    }
};

var send_new_co_serv = function() {
    var msg = "log]/connect " + $("#adressserv").val() + ":" + $("#portserv").val();
    ws.send(msg);
    new_message("log", "log", "Connecting to " + $("#adressserv").val() + "...");
    $(".menu-irc > .selected").click();
};

var send_new_join_chan = function() {
    var msg = $("#idnetwork").val() + "]/join " + $("#adresschan").val();
    ws.send(msg);
    new_message("log", "log", "Joining " + $("#adresschan").val() + "...");
    $(".menu-irc > .selected").click();
};

var send_message = function() {
    if ($(".inputtextirc").val() != '') {
        var buffer_id = $(".main-irc .active a").attr('href').substring(1);
        var txt = $(".inputtextirc").val();
        var msg = buffer_id + "]" + txt;

        console.log(msg);
        ws.send(msg);
        new_message(buffer_id, $(".inputpseudo").val(), txt);
    }
    $(".inputtextirc").val("").focus();
};

var remove_buffer = function(bufferid) {
    //TODO : alert to confirm
    console.log("rm buffer");
    $("#bufferid" + bufferid).hide();
    $("#" + bufferid).hide();
    ws.send(bufferid + "]/close");
};


var add_new_buffer = function(id, name, nick) {
    if (name[0] != '#') {
        new_message("log", "log", "Connected to " + name + "!");
    } else {
        new_message("log", "log", "Joined " + name + "!");
    }
    $('.listbuffer').append('<li id="bufferid' + id + '" onclick="switch_buffer(' + id + ')" ><a href="#' + id + '" data-toggle="tab">' + name + '<span class="remove-buffer" onclick="remove_buffer(' + id + ')">X</span></a></li>');
    $('.contentbuffer').append('<div class="tab-pane bufferchan" id="' + id + '"><table class="table table-striped allmsg"></table> <input type="hidden" class="current-nick" value="log" /></div>');
    $("#" + id + " .current-nick").val(nick);
    $.post("/ajx/backlog", {
        idbuffer: id
    }).done(function(data) {
        $('.contentbuffer #' + id + ' .allmsg').append(data);
        // TODO : JSON this stuff
        check_all_inline_element();
    });
};

var nick_changed = function(oldnick, newnick, buffer) {
    new_message(buffer, "----", oldnick + " is now known as " + newnick);

    if ($("#" + buffer + " .current-nick").val() == oldnick) {
        $("#" + buffer + " .current-nick").val(newnick);
    }
    update_user_list(buffer);
};

var new_message = function(id_buffer, nick, msg) {
    msg = escape_html(msg);
    console.log(msg);
    msg = check_inline_element(msg);
    if (msg.charAt(0) == '/') return;
    $('.contentbuffer #' + id_buffer + ' .allmsg').append('<tr class="msg"><td  class="pseudo nick-' + nick + '">' + nick + '</td><td class="message">' + msg + '</td><td class="time">' + get_timestamp_now() + '</td></tr>');
    $('#' + id_buffer).scrollTop($('#' + id_buffer)[0].scrollHeight);
};

var get_timestamp_now = function() {
    var d = new Date();
    var hour = '0' + d.getHours();
    var min = '0' + d.getMinutes();
    var timestamp = hour.slice(-2) + ":" + min.slice(-2);
    return timestamp;
};

var parse_irc = function(msg) {
    var buff = SplitN(msg, ']', 2);
    switch (buff[0]) {
        case "buffer":
            var buff_nick = buff[2].split(' ');
            console.log("new buffer " + buff[1]);
            add_new_buffer(buff[1], buff_nick[0], buff_nick[1]);
            break;
        case "nick":
            var nicks = buff[2].split(' ');
            nick_changed(nicks[0], nicks[1], buff[1]);
            break;
        case "upul":
            aff_user_list(buff[1]);
            break;
        case "join":
            new_message(buff[1], "<----", buff[2] + " has joined");
            update_user_list(buff[1]);
            break;
        case "part":
            new_message(buff[1], "---->", buff[2] + " has left");
            update_user_list(buff[1]);
            break;
        default:
            new_message(buff[0], buff[1], buff[2]);
            check_mention(buff[0], buff[2]);
            break;
    }
};

$(document).ready(function() {
    var hash = window.location.hash.substring(1);
    ChangePage(hash);
});

var resetnick = function(nick) {
    $(".formirc .add-on").html("<i onclick='changenick()' class='icon-edit'></i><input type='text' disabled='disabled' value='" + nick + "' class='inputpseudo'>");
};

var send_change_nick = function() {
    var nick = $(".inputpseudo").val();
    var buffer_id = $(".main-irc .active a").attr('href').substring(1);
    var msg = buffer_id + "]/nick " + nick;
    ws.send(msg);
    resetnick(nick);
};

var changenick = function() {
    var pseudo = $(".inputpseudo").val();
    $(".inputpseudo").removeAttr("disabled");
    $(".inputpseudo").addClass("activate");
    $(".add-on i").removeClass("icon-edit");
    $(".add-on i").addClass("icon-remove");
    $(".add-on .icon-remove").attr("onclick", "resetnick(\"" + pseudo + "\")");
    $(".add-on").append("<i onclick='send_change_nick()' class='icon-ok'></i>");

    $(".inputpseudo").focus();
};

var check_all_inline_element = function() {
    $('.messagediv').each(function() {
        $(this).html(check_inline_element($(this).html()));
    });
};

var check_inline_element = function(string) {
    // TODO: embeded youtube
    var exp = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
    return string.replace(exp,

    function(url) {
        if (string[string.indexOf(url) - 2] != "=") {
            var clean_url = url.split("?")[0];
            if (clean_url.match(/\.(jpeg|jpg|gif|png)$/) !== null) {
                return '<a target="_blank" href="' + url + '"><img src="' + url + '" width="150" height="150" alt="bou" /></a>';
            }
            return '<a target="_blank" href="' + url + '">' + url + '</a>';
        } else {
            return url;
        }
    });
};

var check_mention = function(id, string) {
    var mynick = $("#" + id + " .current-nick").val();
    if (string.indexOf(mynick) != -1) {
        console.log("mention!");
        notify("test", "hey !");
    }
};

var notify = function(title, body) {
    var icon = "../img/icon_notif.jpg";
    var perm = window.webkitNotifications.checkPermission();
    if (perm === 0 && title !== '' && body !== '' && usersettings.Notify) { //si perm
        var notification = window.webkitNotifications.createNotification(
        icon,
        title,
        body);
        notification.show();
    } else { //sinn request perm
        window.webkitNotifications.requestPermission();
    }
};

var escape_html = function(str) {
    var tagsToReplace = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;'
    };

    function replaceTag(tag) {
        return tagsToReplace[tag] || tag;
    }
    return str.replace(/[&<>]/g, replaceTag);
};


$(".sidebar #menu li").click(function() {
    var name = $(this).find("a").attr("href").substring(1);

    if (name == "irc") {
        load_irc();
    } else {
        ChangePage(name);
    }
});