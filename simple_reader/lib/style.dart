import 'package:flutter/material.dart';

ThemeData appTheme() {
  return ThemeData(
    primarySwatch: Colors.blue,
    splashColor: Colors.green.withAlpha(200),
    textTheme: const TextTheme(
      bodyText1: TextStyle(fontSize: 20.0),
      bodyText2: TextStyle(fontSize: 18.0),
      button: TextStyle(fontSize: 18.0),
    ),
  );
}
