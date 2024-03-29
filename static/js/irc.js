var ws;
var host = window.location.hostname;
var buffers;
var usersettings = null;
var all_buffers = [];

if (ws != null) {
    ws.close();
    ws = null;
}

var get_user_pref = function() {
    $.post("/ajx/getsettings").done(function(data) {
        if (data !== '') {
            usersettings = JSON.parse(data);
        }
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
    var notify = $("#settings-notify").bootstrapSwitch('status');
    var Save_session = $("#settings-savesession").bootstrapSwitch('status');
    $("#form-settings .btn").attr("disabled", true).html("Saving...");
    $.post("/ajx/settings", {
        Notify: notify,
        Save_Session: Save_session
    }).done(function(data) {
        get_user_pref();
        $("#form-settings .btn").attr("disabled", false).html("Save");
        notify_alert($("#message-alert"), "Saved !", "success");
    });
};


if (usersettings === null) {
    get_user_pref();
}

ws = new WebSocket("ws://" + host + ":1112/ws");


ws.onopen = function() {
    if ($("#yuid").val() !== "") {
        ws.send("co]" + $("#yuid").val());
    }
    $("#status-connexion").hide();
};

ws.onmessage = function(e) {
    console.log("receive : " + e.data);
    parse_irc(e.data);

};

ws.onclose = function(e) {
    $("#status-connexion").show();
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
        case "profile":
        update_active_sidebar(page);
        $('.content').load('ajx/profile');
        break;
        default:
        update_active_sidebar("home");
        $('.content').load('ajx/home');
        break;
    }
};

var send_new_co_serv = function() {
    var msg = "log]/connect " + $("#selectnetwork").val();
    ws.send(msg);
    new_message("log", "log", "Connecting to " + $("#selectnetwork").val() + "...");
    $(".menu-irc > .selected").click();
};

var send_new_join_chan = function() {
    var msg = $("#idnetwork").val() + "]/join " + $("#adresschan").val();
    ws.send(msg);
    new_message("log", "log", "Joining " + $("#adresschan").val() + "...");
    $(".menu-irc > .selected").click();
    $("body").click();
};


var load_irc = function() {
    var irc = $("#clientirc").html();
    $('.content').html(irc);
    $(".inputtextirc").keyup(function(event) {
        if (event.keyCode == 13) {
            send_message();
        }
    });

    $('.listbuffer').on('shown', 'a[data-toggle="tab"]', function (e) {
        var id = $(this).attr("href").substring(1);
        switch_buffer(id);
    });

    $("#join-channel").click(function() {
        $("#idnetwork").html("");
        var select = "selected=\"selected\"";
        $.ajax({url: "ajx/set-channels"}).done(function(data){
            if (data !== '') {
                jsonres = JSON.parse(data);
                for (var key in jsonres) {
                    $("#idnetwork").append("<option "+select+" value='" + jsonres[key] + "'>" + key + "</option>");
                    select = "";
                }
            }
        });
    });

    $("#connectnetwork").click(function() {
        $("#selectnetwork").html("");
        var select = "selected=\"selected\"";
        $.ajax({url: "ajx/set-server"}).done(function(data){
            if (data !== '') {
                jsonres = JSON.parse(data);
                for (var key in jsonres) {
                    $("#selectnetwork").append("<option "+select+" value='" + jsonres[key] + "'>" + key + "</option>");
                    select = "";
                }
            }
        });
    });

    $('.switch-userlist').show();
    $('#userlisttab').click();

    $(".menu-irc").on("click", ".dropdown-menu", function(event){
        event.stopPropagation();
    });

    $("#formsetserv").on("submit", "form", function(event){
        event.preventDefault();
        send_new_co_serv();
    });

    $("#formsetchan").on("submit", "form", function(event){
        event.preventDefault();
        send_new_join_chan();
    });
};

var switch_buffer = function(id) {
    if (id == "log") return;
    aff_user_list(id);
    $(".inputpseudo").val(all_buffers[id].nick);
    $(".message iframe").height($(".message iframe").width() / 1.77);
    $('#' + id).scrollTop($('#' + id)[0].scrollHeight);
};

var update_user_list = function(id) {
    var newhtml;

    $.post("/ajx/userslist", {
        channel: id
    }).done(function(data) {
        if (data !== '') {
            $("#userlist-buffer" + id).html($("#userlist-buffer0").html());
            jsonres = JSON.parse(data).UserList;
            $("#userlist-buffer" + id + " .userlist-style").html("<style></style>");
            for (var i = 0; i < jsonres.length; i++) {
                newhtml += add_user_list(jsonres[i].Nick, jsonres[i].Color, jsonres[i].NickClean, id);
            }
            $("#userlist-buffer" + id + " .userlist-style style").html(newhtml);
            $(".userlist-buffer").hide();
            //TODO : rendre visible uniquement la list active
            $("#userlist-buffer" + id).show();
        }
    });
};

var open_set_user = function(name) {

};

var add_friend = function(id, nick) {
    var msg = id + "]/addfriend " + nick;
    ws.send(msg);
};

var get_friends = function(id) {
    $.post("/ajx/getfriends", {
        channel: id
    }).done(function(data) {
        if (data !== '') {
            jsonres = JSON.parse(data).FriendList;
            //
        }
    });
};

var add_user_list = function(name, color, id, buffer) {
    var html = ".nick-" + id + " {  color : " + color + " ; } ";
    $("#userlist-buffer" + buffer + ' .userlist-user').append("<div class=\"btn-group\"> <a class=\"btn dropdown-toggle nick-" + id + "\" data-toggle=\"dropdown\" href=\"#\"> " + name + " <span class=\"caret\"></span> </a><ul class=\"dropdown-menu\"> <li><a href=\"#\">Block</a></li> <li><a onclick=\"add_friend("+buffer+",\'" + name + "\')\" href=\"#\">Add as friend</a></li> <li><a href=\"#\">Private Message</a></li> </ul></div>");
    return html;
};

var aff_user_list = function(id) {

    if ($("#userlist-buffer" + id).length <= 0) {
        var newcss;
        $.post("/ajx/userslist", {
            channel: id
        }).done(function(data) {
            if (data !== '') {
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
            }
        });
    } else {
        $(".userlist-buffer").hide();
        $("#userlist-buffer" + id).show();
    }
};

var send_message = function() {
    if ($(".inputtextirc").val() !== '') {
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
    all_buffers.push(id);
    all_buffers[id] = {nick: nick};
    $('.listbuffer').append('<li id="bufferid' + id + '"    ><a href="#' + id + '" data-toggle="tab">' + name + '<span class="remove-buffer" onclick="remove_buffer(' + id + ')">X</span></a></li>');
    $('.contentbuffer').append('<div class="tab-pane bufferchan" id="' + id + '"><table class="table table-striped allmsg"></table> <input type="hidden" class="current-nick" value="log" /></div>');
    $("#" + id + " .current-nick").val(nick);
    $.post("/ajx/backlog", {
        idbuffer: id
    }).done(function(data) {
        if (data !== '') {
            jsonres = JSON.parse(data);
            var html;
            for (var key in jsonres) {
                html += gen_html_new_message(jsonres[key].Nick, jsonres[key].Message, jsonres[key].Time);
            }
            html += "<tr class'sephist'><td colspan='3'></td></tr> "
            $('.contentbuffer #' + id + ' .allmsg').append(html);
        }
        //check_all_inline_element();
    });
};

var nick_changed = function(oldnick, newnick, buffer) {
    new_message(buffer, "----", oldnick + " is now known as " + newnick);

    if (all_buffers[buffer].nick == oldnick) {
        all_buffers[buffer].nick = newnick;
    }
    update_user_list(buffer);
};

var gen_html_new_message = function(nick, msg, time) {
    nick = nick.substr(0, 15);
    nick = escape_html(nick);
    msg = escape_html(msg);
    msg = check_inline_element(msg);
    if (msg.charAt(0) == '/') return;
    return '<tr class="msg"><td  class="pseudo nick-' + nick + '">' + nick + '</td><td class="message">' + msg + '</td><td class="time">' + time + '</td></tr>';
};

var new_message = function(id_buffer, nick, msg, time) {
    time = typeof time !== 'undefined' ? time : '';
    if (time.length < 1 ) {
        time = get_timestamp_now();
    }
    nick = nick.substr(0, 15);
    nick = escape_html(nick);
    msg = escape_html(msg);
    msg = check_inline_element(msg);
    if (msg.charAt(0) == '/') return;
    $('.contentbuffer #' + id_buffer + ' .allmsg').append('<tr class="msg"><td  class="pseudo nick-' + nick + '">' + nick + '</td><td class="message">' + msg + '</td><td class="time">' + time + '</td></tr>');
    $(".message iframe").height($(".message iframe").width() / 1.77);
    $('#' + id_buffer).scrollTop($('#' + id_buffer)[0].scrollHeight);
};

var get_timestamp_now = function() {
    var d = new Date();
    var hour = '0' + d.getHours();
    var min = '0' + d.getMinutes();
    var timestamp = hour.slice(-2) + ":" + min.slice(-2);
    return timestamp;
};

var disable_buffer = function(id) {

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
        case "leave" :
        disable_buffer(buff[1]);
        break;
        default:
        new_message(buff[0], buff[1], buff[2]);
        check_mention(buff[0], buff[2]);
        break;
    }
};

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
    //TODO handle pseudo in map
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

var youtube_valid = function(url) {
  var p = /^(?:https?:\/\/)?(?:www\.)?(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?$/;
  return (url.match(p)) ? RegExp.$1 : false;
};

var check_inline_element = function(string) {
    var exp = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
    return string.replace(exp,

        function(url) {
            if (string[string.indexOf(url) - 2] != "=") {
                var clean_url = url.split("?")[0];
                var yt = youtube_valid(url);
                if (yt !== false) {
                    return '<iframe width="100%" src="http://www.youtube.com/embed/'+yt+'" frameborder="0" allowfullscreen></iframe>';
                }
                if (clean_url.match(/\.(jpeg|jpg|gif|png)$/) !== null) {
                    var img = new Image();
                    img.src = url;
                    if ( (typeof img.width === 'number') && (typeof img.height === 'number') && img.width <= 1500 && img.height <= 1500) {
                        return '<a rel="lightbox" target="_blank" href="' + url + '"><img src="' + url + '" alt="bou" /></a>';
                    }
                }
                return '<a target="_blank" href="' + url + '">' + url + '</a>';
            } else {
                return url;
            }
        });
};

var check_mention = function(id, string) {
    if (string.indexOf(all_buffers[id].nick) != -1) {
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

var notify_alert = function(div, message, type) {
    div.append("<div class=\"alert alert-"+type+"\"><button type=\"button\" class=\"close\" data-dismiss=\"alert\">&times;</button>"+ message +"</div>");
};

var process_pool_error = function(data) {
    no_err = true;
    if (data !== '') {
       json = JSON.parse(data);
       for (var i in json["errors"]) {
            if (json["errors"][i].length > 0) {
                no_err = false;
                notify_alert($("#message-alert"), json["errors"][i], "error");
            }
        }
    }
    return no_err;
};

$(document).ready(function() {
    var hash = window.location.hash.substring(1);
    ChangePage(hash);

    $("#logo-header").click(function(){
        $('#mainmenutab').click();
        ChangePage("home");
    });

    $(".sidebar #menu li").click(function() {
        var name = $(this).find("a").attr("href").substring(1);
        ChangePage(name);
    });

    $(".content").on("submit", "#login-form", function(event){
      event.preventDefault();
      var no_err = true;
      var $form = $(this),
      mail = $form.find('input[name="InputMail"]').val(),
      pass = $form.find('input[name="InputPass"]').val(),
      button = $form.find('button[type="submit"]'),
      url = $form.attr('action');
      button.attr("disabled", "disabled");
      button.html("Connecting...");
      $.post( url, {
        InputMail: mail,
        InputPass: pass
    }).done(function(data) {
        $("#message-alert").html("");
        if (process_pool_error(data)) {
            window.location.href = "/";
        } else {
            button.removeAttr("disabled");
            button.html("Submit");
        }
    });
});

    $(".content").on("submit", "#register-form", function(event){
      event.preventDefault();
      var no_err = true;
      var $form = $(this),
      mail = $form.find('input[name="InputMail"]').val(),
      pseudo = $form.find('input[name="InputPseudo"]').val(),
      pass = $form.find('input[name="InputPass"]').val(),
      pass2 = $form.find('input[name="InputPassVerif"]').val(),
      button = $form.find('button[type="submit"]'),
      url = $form.attr('action');
      button.attr("disabled", "disabled");
      button.html("Registering...");
      $.post( url, {
        InputMail: mail,
        InputPseudo: pseudo,
        InputPass: pass,
        InputPassVerif: pass2
        }).done(function(data) {
        $("#message-alert").html("");
            button.removeAttr("disabled");
            button.html("Submit");
            if (process_pool_error(data)) {
                $("#register-form").html("");
                notify_alert($("#register-form"), "All good ! Check your mails at " + mail, "success");
            }
        });
    });
});

//JS for handled
if( /Android|webOS|iPhone|iPad|iPod|BlackBerry/i.test(navigator.userAgent) ) {
    if (("#clientirc").length === 0) {
        ChangePage("login");
    }
    ChangePage("irc");
}