import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/api/excetion.dart';
import 'package:simple_reader/routes.dart';

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  createState() => HomeState();
}

class HomeState extends State<HomePage> {
  @override
  void initState() {
    super.initState();
    books();
  }

  Future<void> books() async {
    try {
      await Api.getBooks();
    } on NotLoginException catch (e) {
      Routes.router.navigateTo(context, "/login");
    } catch (error) {
      print(error);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        child: ElevatedButton(
            child: Text("Home"),
            onPressed: () {
              Routes.router.navigateTo(context, "/upload");
            }),
      ),
    );
  }
}
