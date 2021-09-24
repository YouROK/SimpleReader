import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/widgets/blurry_dialog.dart';

class LoginPage extends StatefulWidget {
  @override
  State<LoginPage> createState() => LoginState();
}

class LoginState extends State<LoginPage> {
  late TextEditingController _cntLogin;
  late TextEditingController _cntPass;

  @override
  void initState() {
    super.initState();
    _cntLogin = TextEditingController();
    _cntPass = TextEditingController();
  }

  @override
  void dispose() {
    _cntLogin.dispose();
    _cntPass.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
            heightFactor: 1.5,
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
                  Text("Login"),
                  TextField(
                    controller: _cntLogin,
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
                            await Api.login(_cntLogin.value.text, _cntPass.value.text);
                          } catch (error) {
                            _showDialog(context, "Ошибка", error.toString());
                          }
                        },
                        child: Text("Login")),
                  )
                ],
              ),
            )));
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
