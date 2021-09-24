import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

class BlurryMessage extends StatelessWidget {
  String title;
  String content;

  BlurryMessage(this.title, this.content);
  TextStyle textStyle = TextStyle(color: Colors.black);

  @override
  Widget build(BuildContext context) {
    return BackdropFilter(
      filter: ImageFilter.blur(sigmaX: 1, sigmaY: 1),
      child: AlertDialog(
        title: Text(title, style: textStyle),
        content: Text(content, style: textStyle),
      ),
    );
  }
}
