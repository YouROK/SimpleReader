import 'dart:convert';

import 'package:http/http.dart';
import 'package:simple_reader/api/session.dart';
import 'package:simple_reader/api/utils.dart';

class Api {
  static Session session = Session();

  static Uri getServLink(String link) => session.getLink(link);

  static Future<bool> isLogin() async {
    final resp = await session.get("login");
    return resp.statusCode == 200;
  }

  static Future<Response> login(String email, String pass) async {
    if (email.isEmpty || pass.isEmpty) throw Exception("Заполните все поля");
    var passHash = getPassHash(email, pass);

    return await session.post(
        "/api/login",
        jsonEncode(<String, String>{
          'email': email,
          'pass': passHash,
        }));
  }

  static Future<String> getRegisterEmail(String hash) async {
    if (hash.isEmpty) throw Exception("Произошла ошибка обратитесь к создателю");
    final resp = await session.get("/api/register/$hash");
    return resp.body;
  }

  static Future<void> setRegister(String hash, String email, String login, String pass) async {
    if (hash.isEmpty) throw Exception("Произошла ошибка обратитесь к создателю");
    if (login.isEmpty || email.isEmpty || pass.isEmpty) throw Exception("Заполните все поля");
    var passHash = getPassHash(email, pass);

    await session.post(
        "/api/register/$hash",
        jsonEncode(<String, String>{
          'login': login,
          'pass': passHash,
        }));
  }

  static Future<Map<String, dynamic>> getUserReads() async {
    final resp = await session.get("/api/user/reads");
    return jsonDecode(resp.body);
  }

  static Future<Map<String, dynamic>> getUserStyle() async {
    final resp = await session.get("/api/user/style");
    return jsonDecode(resp.body);
  }

  static Future<Map<String, dynamic>> getBooks() async {
    final resp = await session.get("/api/books");
    return jsonDecode(resp.body);
  }
}
