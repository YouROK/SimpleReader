import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/upload/bookfile.dart';
import 'package:simple_reader/widgets/blurry_dialog.dart';

class UploadPage extends StatefulWidget {
  const UploadPage({Key? key}) : super(key: key);

  @override
  State<UploadPage> createState() => UploadState();
}

class UploadState extends State<UploadPage> {
  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

  var files = List<String?>.empty();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
            heightFactor: 1.0,
            child: Container(
              padding: EdgeInsets.all(20),
              height: 270,
              width: 500,
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.only(topLeft: Radius.circular(10), topRight: Radius.circular(10), bottomLeft: Radius.circular(10), bottomRight: Radius.circular(10)),
                boxShadow: [BoxShadow(color: Colors.grey.withOpacity(0.5), spreadRadius: 0.5, blurRadius: 7, offset: Offset(1, 1))],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text("Загрузить книги"),
                  SizedBox(height: 20),
                  ElevatedButton(
                      onPressed: () async {
                        var result = await FilePicker.platform.pickFiles(
                          allowMultiple: true,
                          allowedExtensions: ['fb2'],
                          type: FileType.custom,
                        );

                        if (result != null) {
                          files = result.paths;
                          setState(() {});
                        } else {
                          // User canceled the picker
                        }
                      },
                      child: Text("Выбрать файлы")),
                  SizedBox(height: 20),
                  Column(children: files.map((e) => BookFile(e)).toList()),
                ],
              ),
            )));
  }

  void _showDialog(BuildContext context, String title, String text) {
    BlurryMessage alert = BlurryMessage(title, text);
    showDialog(
      context: context,
      builder: (BuildContext context) {
        Future.delayed(Duration(seconds: 2), () {
          Navigator.of(context).pop(true);
        });
        return alert;
      },
    );
  }
}
