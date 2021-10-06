import 'package:auto_size_text/auto_size_text.dart';
import 'package:flutter/material.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/widgets/cache_image.dart';

class BookCard extends StatefulWidget {
  BookCard(this.hash);

  final String hash;

  @override
  createState() => BookCardState();
}

class BookCardState extends State<BookCard> {
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
    return Stack(
      children: [
        Positioned.fill(
          child: Card(
              clipBehavior: Clip.antiAlias,
              elevation: 10,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(3.0),
              ),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(3.0),
                child: CacheImage(getCoverLink()),
              )),
        ),
        Positioned.directional(
          textDirection: Directionality.of(context),
          end: 5,
          start: 5,
          bottom: 10,
          child: Container(
            color: Colors.black.withAlpha(190),
            padding: const EdgeInsets.all(5.0),
            child: AutoSizeText(getBookTitle(), style: TextStyle(color: Colors.white), maxLines: 4, textAlign: TextAlign.center),
          ),
        ),
        Positioned.fill(
          child: Material(
            color: Colors.transparent,
            child: InkWell(onTap: () {
              // Routes.navToDetail(context, entity.media_type, entity.id, entity);
            }),
          ),
        )
      ],
    );
  }

  String getCoverLink() {
    try {
      var coverName = info["TitleInfo"]["Coverpage"]["Image"]["Href"] as String;
      if (coverName.startsWith("#")) coverName = coverName.substring(1);
      return Api.session.getLink("/api/book/bin/${widget.hash}/$coverName").toString();
    } catch (e) {
      return Api.session.getLink("/img/cover.jpg").toString();
    }
  }

  String getBookTitle() {
    try {
      return info["TitleInfo"]["BookTitle"] as String;
    } catch (e) {
      return "";
    }
  }
}
