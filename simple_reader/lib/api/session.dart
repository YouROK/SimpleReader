import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:simple_reader/api/excetion.dart';
import 'package:simple_reader/models/user_style.dart';

class Session {
  Map<String, String> headers = {};
  UserStyle userStyle = UserStyle();

  Uri getLink(String link) {
    final host = kReleaseMode ? "book.yourok.ru" : "127.0.0.1:9000";
    if (link.startsWith('/')) return Uri.parse("http://$host$link");
    return Uri.parse("http://$host/$link");
  }

  void exit() {
    headers.clear();
  }

  Future<http.Response> get(String link) async {
    http.Response response = await http.get(getLink(link), headers: headers);
    updateCookie(response);
    if (response.statusCode == 401) throw NotLoginException();
    return response;
  }

  Future<http.Response> post(String link, dynamic data) async {
    http.Response response = await http.post(getLink(link), body: data, headers: headers);
    updateCookie(response);
    return response;
  }

  void updateCookie(http.Response response) {
    String rawCookie = response.headers['set-cookie'] ?? "";
    if (rawCookie.isNotEmpty) {
      int index = rawCookie.indexOf(';');
      headers['cookie'] = (index == -1) ? rawCookie : rawCookie.substring(0, index);
    }
  }
}
