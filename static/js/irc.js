function add_new_buffer(buffer) {
	var id = buffer.replace(/\s/g, "").toLowerCase();

	$('.listbuffer').append('<li><a href="#'+id+'" data-toggle="tab">'+buffer+'</a></li>');
	$('.contentbuffer').append('<div class="tab-pane bufferchan" id="'+id+'"></div>');
}

function new_message(id_buffer, msg) {

	$('.contentbuffer #'+id_buffer+' .allmsg').append('<tr class="msg"><td class="pseudo">xxx</td><td class="message">'+msg+'</td><td class="time">13h30</td></tr>');

}


function parse_irc(msg) {
	var buff = msg.split(']');
	switch(buff[0])
		{
		case "successserv":
		  console.log("Connect to server "+buff[1]);
		  add_new_buffer(buff[1])
		  break;
		case "successchan":
		  console.log("Connect to channel "+buff[1]);
		  var servchan = buff[1].split('?');
		  add_new_buffer(servchan[0]+" "+servchan[1]);
		  break;
		  default :
		  	var servchan = buff[0].split('?');
		  	var id_buffer = servchan[0]+servchan[1];
		 	new_message(id_buffer.toLowerCase(), buff[1]);
		  break
		}
}


$(document).ready(function() {
  $("tr:even").css("background-color", "#f7f7f9");
  $("tr:odd").css("background-color", "#fff");
  $(".bufferchan").scrollbars();
});


$(".formirc button").click(function () {
	if ($(".formirc input").val() != '') {
		//console.log($(".active a").html());
		 var info = $("#main-irc .active a").html().split(' ');
		 var serv = info[0];
		 var channel = info[1];
		//console.log(serv+channel);
		var msg = serv+"?"+channel+"]"+$(".formirc input").val();
		console.log(msg);
		ws.send(msg);
	}
	$(".formirc input").val("");
});

