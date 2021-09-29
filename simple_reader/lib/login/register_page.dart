import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/routes.dart';
import 'package:simple_reader/widgets/blurry_dialog.dart';

class RegisterPage extends StatefulWidget {
  final String hash;

  const RegisterPage(this.hash, {Key? key}) : super(key: key);

  @override
  State<RegisterPage> createState() => RegisterState();
}

class RegisterState extends State<RegisterPage> {
  late TextEditingController _cntLogin;
  late TextEditingController _cntPass;
  late TextEditingController _cntEmail;

  var email = "";

  @override
  void initState() {
    super.initState();
    _cntLogin = TextEditingController();
    _cntPass = TextEditingController();
    _cntEmail = TextEditingController();
    loadRegEmail();
  }

  void loadRegEmail() async {
    email = await Api.getRegisterEmail(widget.hash);
    if (email.isEmpty) Routes.router.navigateTo(context, "/notfound");
    _cntEmail.text = email;
    setState(() {});
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
              height: 375,
              width: 500,
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.only(topLeft: Radius.circular(10), topRight: Radius.circular(10), bottomLeft: Radius.circular(10), bottomRight: Radius.circular(10)),
                boxShadow: [BoxShadow(color: Colors.grey.withOpacity(0.5), spreadRadius: 0.5, blurRadius: 7, offset: Offset(1, 1))],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (email.isNotEmpty)
                    Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                      Text("Почта"),
                      TextField(
                        controller: _cntEmail,
                        enabled: false,
                      ),
                      SizedBox(height: 20),
                    ]),
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
                  if (email.isEmpty) SizedBox(height: 40, child: Center(child: CircularProgressIndicator())),
                  if (email.isNotEmpty)
                    SizedBox(
                      height: 40,
                      width: double.infinity,
                      child: ElevatedButton(
                          onPressed: () async {
                            if (email.isEmpty) {
                              _showDialog(context, "", "Подождите загрузке данных...");
                            } else {
                              try {
                                await Api.setRegister(widget.hash, email, _cntLogin.value.text, _cntPass.value.text);
                                Routes.router.navigateTo(context, "/");
                              } catch (error) {
                                _showDialog(context, "Ошибка", error.toString());
                              }
                            }
                          },
                          child: Text("Отправить")),
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
