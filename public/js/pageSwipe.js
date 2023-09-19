var xDown,yDown,tDown,isDown;
var sLeft,sRight;

function initSwipe(swipeLeft, swipeRight){
	sLeft = swipeLeft
	sRight = swipeRight
	$(window).scroll(function() {scrolling = true});

	$("#content").on("swipeleft", sLeft);
	$("#content").on("swiperight", sRight);
}

function mouseDown(e){
	isDown = true;
	xDown = e.screenX;
	yDown = e.screenY;
	tDown = e.timeStamp;
	//$("#header").html("D "+xDown+" "+yDown+" "+tDown);
}

function mouseUp(e){
	if (isDown){
		xDown = e.screenX-xDown;
		yDown = e.screenY-yDown;
		tDown = e.timeStamp-tDown;

		if( tDown>40 && tDown<500 && Math.abs(yDown/xDown)<1.5 && Math.abs(xDown)>50 ) {
			if (xDown>0 && sRight)
				sRight(e);
			if (xDown<0 && sLeft)
				sLeft(e);
		}
		// $("#header").html("U "+xDown+" "+yDown+" "+tDown);
	}
	isDown = false;
}

function slideNext(done){
	var width = $(window).width();
    var time = 400;
	var curr = $('#currentPage');
	var next = $('#nextPage');
    curr.css({position: "relative"}).animate({left: -width, opacity: 0.25}, time,function(){
		scrollToTop();
		next.css({position: "relative", left: width}).show();
		curr.hide();
		next.animate({left: 0, opacity: 1}, time, function() {
			next.attr('id',"currentPage");
			curr.attr('id',"nextPage").empty();
			if (done)
				done();
		});
	});
}

function slidePrev(done){
	var width = $(window).width();
    var time = 400;
	var curr = $('#currentPage');
	var next = $('#nextPage');
    curr.css({position: "relative"}).animate({left: width, opacity: 0.25}, time,function(){
		scrollToTop();
		next.css({position: "relative", left: -width}).show();
		curr.hide();
		next.animate({left: 0, opacity: 1}, time, function() {
			next.attr('id',"currentPage");
			curr.attr('id',"nextPage").empty();
			if (done)
				done();
		});
	});
}


function fadeNext(done) {
	var time = 400;
	var curr = $('#currentPage');
	var next = $('#nextPage');
	$('#header').fadeOut(time);
	curr.fadeOut(time, function () {
		$('#header').fadeIn(time);
		scrollToTop();
		next.fadeIn(time / 2, function () {
			curr.hide();
			next.attr('id',"currentPage");
			curr.attr('id',"nextPage").empty();
			if (done)
				done();
		});
	});
}

function scrollToTop(){
	$.mobile.silentScroll($('#header').offset().top+$('#header').height());
}
