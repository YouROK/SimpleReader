function setBook(book){
    if (book){
        if ($.jStorage.storageSize() >= 4194304){
            var min = null;
            var books = $.jStorage.index();
            for (var index = 0; index < books.length; ++index){
                var book = $.jStorage.get(books[index]);
                if (book && book.lastread){
                    if (!min || book.lastread < min[1])
                        min = [book.id, book.lastread];
                }
            }
            setPage(getBook(min[0]));
            clearBook(min[0]);
        }
        $.jStorage.set(book.id, book, {TTL: 1209600000});//2 week
    }
}

function getBook(id){
    var cont = $.jStorage.get(id);
    if (!cont)
        cont = {id: id, buffer: {}, count: 0, readpage: 0, nextpage: 0, lastline: 0, readedpages: [], lastread: new Date().getMilliseconds()};
    return cont;
}

function getLastReadPage(book){
	if(book.readedpages.length > 0){
	    return book.readedpages[book.readedpages.length-1];
	}else{
		if (book.readpage == 0)
			return -1;
	    if(book.readpage - 10 > 0)
	        return book.readpage - 10;
	    else
	        return 0;
	}
}

function loadReadPages(line, done){
	if (line>book.count)
		return;
	var wndHeight = $(window).height() - ($('#currentPage').offset().top + $('#footer').outerHeight()+60);
	getPage(book, line, wndHeight, function(currentContent, lastLine){
		currPage = $("#currentPage");
		currPage.empty();
		currPage.html(currentContent);
		if (done)
			done(lastLine);
	})
}

function loadNextPage(line, done){
	if (line>book.count)
		return;
	var wndHeight = $(window).height() - ($('#currentPage').offset().top + $('#footer').outerHeight()+60);
	getPage(book, line, wndHeight, function(nextContent,lastLine){
		nextPage = $("#nextPage");
		nextPage.empty();
		nextPage.html(nextContent);
		if (done)
			done(lastLine);
	})
}

function getPage(book, start, height, done){
	$("#canvas").append("<div id='tmp' class='canvas' style='display:none;'></div>");
	var content = $("#tmp");
	var html = "";
    var notes = "";
	var curLine = start;

	content.empty();

	while(content.height() < height && curLine < book.count){
		var line = getLine(book, curLine);
		if (line){
			curLine ++;
			content.append(line[0]);
            html += line[0];
            if (line[1]){
				notes += line[1];
			}
		}else {
			done = null;
			break;
		}
	}
	if (notes)
        html += "<hr align='center' width='90%' size='1'/><div align=left>" + notes + "</div>";
	content.remove();
	if (done)
		done(html, curLine);
}

function getLine(book, index){
	if (book == null || index>=book.count)
		return;

	if (!book.buffer.hasOwnProperty(index) || !book.buffer[index]){
        var json = loadContent(book.id, index);
        if (json){
            for ( var i=0, l=json.length; i<l; i++) {
        			book.buffer[i+index] = json[i];
			}
        }
    }
	if(book.buffer.hasOwnProperty(index))
		return [book.buffer[index].Text, book.buffer[index].Note];
	else
		return null;
}

function loadContent(id, start){
	var json = null;
	try{
	    $.ajax({
	        type: 'POST',
	        url: '/getcontent',
	        data: JSON.stringify({BookId: id, Start: start, Count: 600}),
	        contentType: 'application/json',
	        dataType: 'json',
			cache: false,
	        timeout: 15000,
	        async: false
	    })
		.done(function (data){
			json = data;
		})
		.fail(function (data){
			ons.notification.toast("Ошибка при загрузке страницы", {timeout: 2000});
			json = null;
		});
	} catch(e) {
		ons.notification.toast("Ошибка при загрузке страницы", {timeout: 2000});
		json = null;
	}
    return json;
}

function setPage(book){
    $.ajax({
        type: 'POST',
        url: '/book/' + book.id + '/' + book.readpage
    }).fail(function (result)
        {
			showProcent(book.nextpage * 100 / book.count, "red")
        }).done(function ()
        {
			showProcent(book.nextpage * 100 / book.count, "#436f8c")
        });
}

function showProcent(prc, clr){
	if (prc != null)
		$('#progress_bar').css('width', prc + '%');
	if (clr != null)
		$('#progress_bar').css({backgroundColor: clr});

}
