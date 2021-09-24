import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

class Session {
  Map<String, String> headers = {};

  Uri getLink(String link) {
    final host = kReleaseMode ? "book.yourok.ru" : "127.0.0.1:9000";
    if (link.startsWith('/')) return Uri.parse("http://$host$link");
    return Uri.parse("http://$host/$link");
  }

  Future<String> get(String link) async {
    http.Response response = await http.get(getLink(link), headers: headers);
    updateCookie(response);
    return response.body;
  }

  Future<String> post(String link, dynamic data) async {
    http.Response response = await http.post(getLink(link), body: data, headers: headers);
    updateCookie(response);
    return response.body;
  }

  void updateCookie(http.Response response) {
    String rawCookie = response.headers['set-cookie'] ?? "";
    if (rawCookie.isNotEmpty) {
      int index = rawCookie.indexOf(';');
      headers['cookie'] = (index == -1) ? rawCookie : rawCookie.substring(0, index);
    }
  }
}
