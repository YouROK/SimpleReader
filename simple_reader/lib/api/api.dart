import 'dart:convert';

import 'package:http/http.dart';
import 'package:simple_reader/api/excetion.dart';
import 'package:simple_reader/api/session.dart';
import 'package:simple_reader/api/utils.dart';
import 'package:simple_reader/models/book_info.dart';

class Api {
  static Session session = Session();

  static Uri getServLink(String link) => session.getLink(link);

  static Future<bool> isLogin() async {
    try {
      final resp = await session.get("/api/login");
      return resp.statusCode == 200;
    } on NotLoginException catch (error) {
      return false;
    }
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

  static Future<List<BookInfo>> getUserReads() async {
    final resp = await session.get("/api/user/reads");
    final list = jsonDecode(resp.body);
    return (list as List<Map<String, dynamic>>).map((e) => BookInfo.fromJson(e)).toList();
  }

  static Future<Map<String, dynamic>> getUserStyle() async {
    final resp = await session.get("/api/user/style");
    return jsonDecode(resp.body);
  }

  static Future<Map<String, dynamic>> getBooks() async {
    final resp = await session.get("/api/book/all");
    return jsonDecode(resp.body);
  }

  static Future<void> uploadBook(
    String file,
  ) async {
    final resp = await session.get("/api/book/all");
    return jsonDecode(resp.body);
  }

  static Future<Map<String, dynamic>> getBookDesc(String hash) async {
    final resp = await session.get("/api/book/desc/$hash");
    return jsonDecode(resp.body);
  }

  static Future<List<String>> search(String query) async {
    final resp = await session.get("/api/search?query=$query");
    final json = jsonDecode(resp.body);
    return (json as List<dynamic>).map((e) => e.toString()).toList();
  }
}
