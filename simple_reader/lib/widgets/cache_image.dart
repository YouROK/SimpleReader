import 'package:auto_size_text/auto_size_text.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:simple_reader/api/api.dart';

class CacheImage extends StatelessWidget {
  CacheImage(String this.url);

  final String url;
  Widget placeholder = Center(child: CircularProgressIndicator());

  @override
  Widget build(BuildContext context) {
    try {
      return CachedNetworkImage(
          fit: BoxFit.fill,
          placeholder: (context, url) => placeholder,
          httpHeaders: Api.session.headers,
          imageUrl: url,
          errorWidget: (context, url, error) {
            print(error);
            return FittedBox(child: CachedNetworkImage(imageUrl: Api.session.getLink("/img/cover.jpg").toString()), fit: BoxFit.fill);
          });
    } catch (e) {
      return Expanded(
          child: Center(
              child: AutoSizeText(
        "Не удалось загрузить",
        maxLines: 2,
      )));
    }
  }
}
