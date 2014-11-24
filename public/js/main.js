function DialogShow(text, time, done){
    if ($("#dialog").length > 0)
        $("#dialog").remove();

    var id = $.mobile.activePage.attr("id");
    if (!id)
        id = "content";
    $("#" + id).append("<div data-role='popup' id='dialog' data-transition='fade' data-position-to='window'><p>" + text + "</p></div>");

    var dialog = $("#dialog");
    dialog.popup({ history: false });
    dialog.css({"padding": "2em", "opacity": "0.85"});
    dialog.popup("open").on({
        popupafteropen: function (){
            if (!time)
                time = 2000;
            setTimeout(function (){
                dialog.popup('close')
				if (done)
					done();
            }, time);
        }
    });
}

var menuShow = false
function ToggleMenu(){
	if (menuShow){
		$("#menu").hide();
		menuShow = false;
	}else{
		$("#menu").show();
		menuShow = true;
	}
}

function Send(url, json, callback)
{
    var request = $.ajax({
        type: 'POST',
        url: url,
        data: JSON.stringify(json),
        contentType: 'application/json',
        dataType: 'json',
        timeout: 30000,
        async: false
    });

    request.done(function (data)
    {
        if (!callback && data)
        {
            if (!data.Message)
                DialogShow("Ok");
            else
                DialogShow(data.Message);
        }
		if (callback)
			callback(data);
    });
    request.fail(function (data)
    {
		var msg = ""
		if (data && data.responseText)
			msg = data.responseText;
        DialogShow("Произошла ошибка, повторите еще раз позже"+msg);
    });
}
