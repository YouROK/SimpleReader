import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/routes.dart';

import 'bookfile.dart';

class UploadPage extends StatefulWidget {
  const UploadPage({Key? key}) : super(key: key);

  @override
  State<UploadPage> createState() => UploadState();
}

class UploadState extends State<UploadPage> {
  @override
  void initState() {
    super.initState();
    checkLogin();
  }

  Future<bool> checkLogin() async {
    final isl = await Api.isLogin();
    if (!isl) Routes.router.navigateTo(context, "/login");
    return isl;
  }

  var bookFiles = List<BookFile>.empty(growable: true);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(centerTitle: true, title: Text("Загрузить книги")),
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
                heightFactor: 1.0,
                child: Container(
                  padding: EdgeInsets.all(20),
                  margin: EdgeInsets.all(10),
                  width: 700,
                  decoration: BoxDecoration(
                    color: Colors.grey[100],
                    borderRadius: BorderRadius.only(topLeft: Radius.circular(10), topRight: Radius.circular(10), bottomLeft: Radius.circular(10), bottomRight: Radius.circular(10)),
                    boxShadow: [BoxShadow(color: Colors.black, spreadRadius: 0.5, blurRadius: 7, offset: Offset(1, 1))],
                  ),
                  child: Column(
                    children: [
                      Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                        SizedBox(
                            height: 40,
                            width: 300,
                            child: ElevatedButton(
                                onPressed: () async {
                                  final isl = await checkLogin();
                                  if (!isl) return;

                                  bookFiles.clear();
                                  setState(() {});
                                  var result = await FilePicker.platform.pickFiles(
                                    allowMultiple: true,
                                    allowedExtensions: ['fb2'],
                                    type: FileType.custom,
                                  );

                                  if (result != null) {
                                    for (var e in result.paths) {
                                      final ctrl = BookFileController();
                                      bookFiles.add(BookFile(e, ctrl));
                                    }
                                    setState(() {});
                                  }
                                },
                                child: Text("Выбрать файлы"))),
                        SizedBox(width: 20),
                        SizedBox(
                            height: 40,
                            width: 300,
                            child: ElevatedButton(
                                onPressed: () async {
                                  final isl = await checkLogin();
                                  if (!isl) return;
                                  for (var el in bookFiles) {
                                    await el.controller.execute();
                                    setState(() {});
                                  }
                                  bookFiles.clear();
                                  setState(() {});
                                },
                                child: Text("Загрузить"))),
                      ]),
                      SizedBox(height: 20),
                      Expanded(
                        child: SingleChildScrollView(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.start,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: bookFiles,
                          ),
                        ),
                      ),
                    ],
                  ),
                ))));
  }
}
