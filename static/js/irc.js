function add_new_buffer(buffer) {
	var id = buffer[1]

	$('.listbuffer').append('<li><a href="#'+id+'" data-toggle="tab">'+buffer[2]+'</a></li>');
	$('.contentbuffer').append('<div class="tab-pane bufferchan" id="'+id+'"><table class="table table-striped allmsg"></table></div>');
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
		  add_new_buffer(buff)
		  break;
		case "successchan":
		  console.log("Connect to channel "+buff[1]);
		  add_new_buffer(buff);
		  break;
		  default :
		 	new_message(buff[0], buff[1]);
		  break
		}
}


$(document).ready(function() {
  $("tr:even").css("background-color", "#f7f7f9");
  $("tr:odd").css("background-color", "#fff");
  $(".bufferchan").scrollbars();
});

$(".formirc input").keyup(function(event){
    if(event.keyCode == 13){
        $(".formirc button").click();
    }
});


$(".formirc button").click(function () {
	if ($(".formirc input").val() != '') {
		//console.log($(".active a").html());
		 var buffer_id = $("#main-irc .active a").attr('href').substring(1);
		var msg = buffer_id+"]"+$(".formirc input").val();
		console.log(msg);
		ws.send(msg);
	}
	$(".formirc input").val("").focus();
});

