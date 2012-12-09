//piou piou

function actionmenu(page) {
		$('#content').load('ajx/'+page);
}


$("#sidebar li").click(function () {
	$("#sidebar li").removeClass("active");
    $(this).addClass("active");
    actionmenu($(this).find("a").attr("href").substring(1));
});

