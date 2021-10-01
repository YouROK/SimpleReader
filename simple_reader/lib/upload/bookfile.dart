import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:simple_reader/api/api.dart';

import 'file_service.dart';

class BookFile extends StatefulWidget {
  const BookFile(this.filePath, {Key? key}) : super(key: key);
  final filePath;

  @override
  createState() => BookFileState();
}

class BookFileState extends State<BookFile> {
  var progress = 0.0;
  var err = "";
  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [Text(widget.filePath), LinearProgressIndicator(value: progress), Text(err)]);
  }

  @override
  void initState() {
    final link = Api.session.getLink("/api/book/upload");
    try {
      FileService.fileUploadMultipart(
        url: link.toString(),
        file: widget.filePath,
        onUploadProgress: (sentBytes, totalBytes) {
          progress = sentBytes / totalBytes;
          setState(() {});
        },
      );
    } catch (error) {
      err = error.toString();
      setState(() {});
    }
  }
}
