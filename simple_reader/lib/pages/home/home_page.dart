import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/api/excetion.dart';
import 'package:simple_reader/models/book_info.dart';
import 'package:simple_reader/routes.dart';

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  createState() => HomeState();
}

class HomeState extends State<HomePage> {
  List<BookInfo> userReads = List.empty();

  @override
  void initState() {
    super.initState();
    inits();
  }

  Future<void> inits() async {
    try {
      userReads = await Api.getUserReads();
      userReads.sort((a, b) => a.LastRead.compareTo(b.LastRead));
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
            height: double.infinity,
            width: double.infinity,
            decoration: BoxDecoration(
              image: DecorationImage(
                image: NetworkImage(Api.session.getLink("/img/back.jpg").toString()),
                fit: BoxFit.cover,
                isAntiAlias: true,
              ),
            ),
            child: Center(
              child: Container(
                padding: EdgeInsets.all(20),
                height: 285,
                width: 500,
                decoration: BoxDecoration(
                  color: Colors.grey[100],
                  borderRadius: BorderRadius.only(topLeft: Radius.circular(10), topRight: Radius.circular(10), bottomLeft: Radius.circular(10), bottomRight: Radius.circular(10)),
                  boxShadow: [BoxShadow(color: Colors.black, spreadRadius: 0.5, blurRadius: 7, offset: Offset(1, 1))],
                ),
                child: Container(
                  child: Column(
                    children: [
                      SizedBox(
                          height: 40,
                          width: double.infinity,
                          child: ElevatedButton(
                              onPressed: () {
                                Routes.router.navigateTo(context, "/upload");
                              },
                              child: Text("Загрузть книги"))),
                      SizedBox(height: 10),
                      SizedBox(
                          height: 40,
                          width: double.infinity,
                          child: ElevatedButton(
                              onPressed: () {
                                Routes.router.navigateTo(context, "/mybooks");
                              },
                              child: Text("Мои книги"))),
                      SizedBox(height: 10),
                      SizedBox(
                          height: 40,
                          width: double.infinity,
                          child: ElevatedButton(
                              onPressed: () {
                                Routes.router.navigateTo(context, "/search");
                              },
                              child: Text("Поиск"))),
                      SizedBox(height: 10),
                      SizedBox(
                          height: 40,
                          width: double.infinity,
                          child: ElevatedButton(
                              onPressed: () {
                                // Routes.router.navigateTo(context, "/upload");
                              },
                              child: Text("Настройки"))),
                      SizedBox(height: 15),
                      SizedBox(
                          height: 40,
                          width: double.infinity,
                          child: ElevatedButton(
                              onPressed: () {
                                Api.session.exit();
                                Routes.router.navigateTo(context, "/login");
                              },
                              child: Text("Выход"))),
                    ],
                  ),
                ),
              ),
            )));
  }
}
