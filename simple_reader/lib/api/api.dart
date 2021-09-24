import 'dart:convert';

import 'package:crypto/crypto.dart';
import 'package:simple_reader/api/session.dart';

class Api {
  static Session session = Session();
  static String ID = "";

  static Future<void> login(String login, String pass) async {
    if (login.isEmpty || pass.isEmpty) throw Exception("Заполните все поля");
    var passHash = sha256.convert(utf8.encode(login + pass + "1632387031")).toString();

    await session.post(
        "login",
        jsonEncode(<String, String>{
          'login': login,
          'pass': passHash,
        }));
  }
}
