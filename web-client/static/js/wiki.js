$(function(){
	//将首页自我介绍的图片居中，且留出足够的margin-bottom
	//$("img[alt=ME]").wrap("<a class='intro' data-title='Ollie' data-description='This is a test!'></a>");
	$("img[alt=ME]").wrap("<a></a>");
	$("img[alt=ME]").closest("a").addClass("intro")
		//.attr("data-title","Ollie Liu")
		//.attr("data-description","这里有我的点点滴滴!");
	$('img[alt=ME]').closest("p").attr("align","center").css("marginBottom",90);
	//令图片居中显示
	$('img').closest("p").attr("align","center");
	//自动检body大小，如果body小于屏幕高度，则将body高度设置为屏幕高度
	//减去20是因为body在上下，分别有10px的margin，这样不至于出现滚动条
	//if($('body').height()<$(window).height()){
		//$('body').height($(window).height()-20);
	//}
});
