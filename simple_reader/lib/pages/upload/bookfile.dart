import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/api/excetion.dart';

import 'file_service.dart';

class BookFile extends StatefulWidget {
  BookFile(this.filePath, this.controller, {Key? key}) : super(key: key);
  final filePath;
  final BookFileController controller;
  var complete = false;

  @override
  createState() => BookFileState();
}

class BookFileState extends State<BookFile> {
  var progress = 0.0;
  var err = "";

  @override
  void initState() {
    super.initState();
    widget.controller.addListener(upload);
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
        height: 50,
        width: double.infinity,
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(widget.filePath),
          if (progress > 0 && progress < 1) LinearProgressIndicator(value: progress),
          if (progress == 1.0) Text("Готово", style: TextStyle(fontSize: 13)),
          if (err.isNotEmpty) Text(err, style: TextStyle(color: Colors.red)),
        ]));
  }

  Future<void> upload() async {
    if (progress == 1.0 && err.isEmpty) return;

    final link = Api.session.getLink("/api/book/upload");
    try {
      await FileService.fileUploadMultipart(
        url: link.toString(),
        file: widget.filePath,
        onUploadProgress: (sentBytes, totalBytes) {
          progress = sentBytes / totalBytes;
          if (mounted) setState(() {});
        },
      ).catchError((error) {
        if (error is NotLoginException)
          err = "Ошибка пользователь не зашел на сайт";
        else
          err = error.toString();
        widget.complete = false;
        if (mounted) setState(() {});
        return "";
      });
      widget.complete = true;
      if (mounted) setState(() {});
    } catch (error) {
      widget.complete = false;
      if (mounted) setState(() {});
      err = error.toString();
    }
  }
}

class BookFileController {
  late Function listener;
  addListener(Function fn) {
    listener = fn;
  }

  Future<void> execute() async {
    await listener();
  }
}
