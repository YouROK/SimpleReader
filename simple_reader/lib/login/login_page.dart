import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/routes.dart';
import 'package:simple_reader/widgets/blurry_dialog.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => LoginState();
}

class LoginState extends State<LoginPage> {
  late TextEditingController _cntEMail;
  late TextEditingController _cntPass;

  @override
  void initState() {
    super.initState();
    _cntEMail = TextEditingController();
    _cntPass = TextEditingController();
    _cntEMail.text = "8yourok8@mail.ru";
    _cntPass.text = "19851985";
  }

  @override
  void dispose() {
    _cntEMail.dispose();
    _cntPass.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Container(
            height: double.infinity,
            width: double.infinity,
            decoration: BoxDecoration(
              image: DecorationImage(
                image: NetworkImage(Api.session.getLink("/img/back.png").toString()),
                fit: BoxFit.cover,
                isAntiAlias: true,
              ),
            ),
            child: Center(
                child: Container(
              padding: EdgeInsets.all(20),
              height: 270,
              width: 500,
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.only(topLeft: Radius.circular(10), topRight: Radius.circular(10), bottomLeft: Radius.circular(10), bottomRight: Radius.circular(10)),
                boxShadow: [BoxShadow(color: Colors.grey.withOpacity(0.5), spreadRadius: 0.5, blurRadius: 7, offset: Offset(1, 1))],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text("Почта"),
                  TextField(
                    controller: _cntEMail,
                    onSubmitted: (value) {},
                  ),
                  SizedBox(height: 20),
                  Text("Пароль"),
                  TextField(
                    controller: _cntPass,
                    obscureText: true,
                    enableSuggestions: false,
                    autocorrect: false,
                  ),
                  SizedBox(height: 30),
                  SizedBox(
                    height: 40,
                    width: double.infinity,
                    child: ElevatedButton(
                        onPressed: () async {
                          try {
                            final resp = await Api.login(_cntEMail.value.text, _cntPass.value.text);
                            if (resp.statusCode == 200) {
                              Routes.router.navigateTo(context, "/");
                            } else {
                              _showDialog(context, "Ошибка", resp.body);
                            }
                          } catch (error) {
                            _showDialog(context, "Ошибка", error.toString());
                          }
                        },
                        child: Text("Войти")),
                  )
                ],
              ),
            ))));
  }

  void _showDialog(BuildContext context, String title, String text) {
    BlurryMessage alert = BlurryMessage(title, text);
    showDialog(
      context: context,
      builder: (BuildContext context) {
        Future.delayed(Duration(seconds: 2), () {
          Navigator.of(context).pop(true);
        });
        return alert;
      },
    );
  }
}
