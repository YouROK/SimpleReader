import 'dart:convert';

import 'package:crypto/crypto.dart';

String getPassHash(String email, String pass) {
  return sha256.convert(utf8.encode(email.toLowerCase() + pass + "1632387031")).toString();
}
