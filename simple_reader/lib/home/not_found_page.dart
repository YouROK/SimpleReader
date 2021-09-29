import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/routes.dart';

class NotFoundPage extends StatelessWidget {
  const NotFoundPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
            child: Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text("Страница не найдена"),
        SizedBox(height: 20),
        SizedBox(height: 40, child: ElevatedButton(onPressed: () => Routes.router.navigateTo(context, "/"), child: Text("Перейти на главную"))),
      ],
    )));
  }
}
