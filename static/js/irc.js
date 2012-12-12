$(document).ready(function() {
  $("tr:even").css("background-color", "#f7f7f9");
  $("tr:odd").css("background-color", "#fff");
  $(".bufferchan").scrollbars();
});


$(".formirc button").click(function () {
	ws.send($(".formirc input").val());
	$(".formirc input").val("");
});

