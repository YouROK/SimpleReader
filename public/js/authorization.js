function GetKey(){
	var key = null;
    $.ajax({
        type: 'GET',
        url: '/getkey',
        timeout: 30000,
        async: false,
        dataType: 'json'
    }).always(function (data){
            if (data) key = data;
        });
    return key;
}

function Login(backPage){
	const username = document.querySelector('#username').value;
	const password = document.querySelector('#password').value;
	const captcha = document.querySelector('#captchaInp').value;

	if (!username || !password) {
		ons.notification.alert('Заполните все поля',{title:""});
		return;
	}

	var reg = new RegExp("[^a-zA-Z0-9@#*_+]+");
	if (reg.test(username) || username.length > 30) {
		ons.notification.alert('Указанное имя недопустимо, содержит недопустимые символы или имеет длину более 30 знаков',{title:""});
		return;
	}

	var key = GetKey();
	if (key){
		var pass = b64_sha1(password);
       	var name = username;

		var rsa = new RSAKey();
		rsa.setPublic(key.Mod, key.Exp);

		var cryptPass = hex2b64(rsa.encrypt(pass));
		var cryptName = hex2b64(rsa.encrypt(name));

		var json = {Login: cryptName, Password: cryptPass, Captcha: captcha};
		var request = $.ajax({
           		type: 'POST',
           		url: '/login',
           		data: JSON.stringify(json),
           		contentType: 'application/json',
           		dataType: 'json',
           		timeout: 30000,
           		async: false
       		});
		request.always(function (data,status){
            if (data){
				if (data.Code !== 0){
					if (data.Captcha){
						$('#captchaDiv').show();
						$('#captchaImg').attr('src', data.Captcha);
					}
					if(data.Msg)
						ons.notification.alert(data.Msg,{title:""});
				}else{
					if (backPage)
						window.location = backPage;
					else
						window.location = "/";
				}
            }
        });
	}
}

function Register(name,pass,pass2,email,captcha){
		if (!email || !captcha || !name || !pass || !pass2){
			DialogShow("Заполните все поля", 2000);
			return;
		}

		if (pass != pass2){
			DialogShow("Пароли не совпадает", 2000);
			return;
		}

		var reg = new RegExp("[^a-zA-Z0-9@#*_+]+");
		if (reg.test(name) || name.length > 30){
			DialogShow("Указанное имя недопустимо, содержит недопустимые символы или имеет длину более 30 знаков", 2000);
			return;
		}

		$.mobile.loading("show");

		var key = GetKey();
		if (key){
			pass = b64_sha1(pass);

			var rsa = new RSAKey();
			rsa.setPublic(key.Mod, key.Exp);

			var cryptPass = hex2b64(rsa.encrypt(pass));
			var cryptName = hex2b64(rsa.encrypt(name));
			var cryptEmail = hex2b64(rsa.encrypt(email));

			var json = {Login: cryptName, Password: cryptPass, Email: cryptEmail, Captcha: captcha};

			var request = $.ajax({
				type: 'POST',
				url: '/registration',
				data: JSON.stringify(json),
				contentType: 'application/json',
				dataType: 'json',
				timeout: 30000,
				async: false
			});

			request.done(function (data){
				if (data){
					if (data.Code == 0)
						window.location = "/";
					else{
						if (data.Msg)
							DialogShow(data.Msg, 2000, location.reload);
						else
							DialogShow("Произошла ошибка, повторите еще раз позже");
						if (data.Captcha)
							$("#captchaImage").attr("src", data.Captcha);
					}
				}
			});
			request.fail(function (){
				DialogShow("Произошла ошибка, повторите еще раз позже", 2000, location.reload);
			});
		}
		$.mobile.loading("hide");
	}
