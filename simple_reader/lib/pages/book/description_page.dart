import 'package:flutter/material.dart';
import 'package:simple_reader/api/api.dart';

class DescPage extends StatefulWidget {
  const DescPage(this.hash, {Key? key}) : super(key: key);

  final String hash;

  @override
  createState() => DescState();
}

class DescState extends State<DescPage> {
  Map<String, dynamic> info = {};
  @override
  void initState() {
    super.initState();
    loadInfo();
  }

  Future<void> loadInfo() async {
    info = await Api.getBookDesc(widget.hash);
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
    );
  }
}
