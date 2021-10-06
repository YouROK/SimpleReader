import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:simple_reader/api/api.dart';
import 'package:simple_reader/api/excetion.dart';
import 'package:simple_reader/routes.dart';
import 'package:simple_reader/widgets/book_card.dart';

class SearchPage extends StatefulWidget {
  const SearchPage({Key? key}) : super(key: key);

  @override
  State<SearchPage> createState() => SearchState();
}

class SearchState extends State<SearchPage> {
  late TextEditingController _cntSearch;
  List<BookCard> books = List.empty();
  String query = "";

  @override
  void initState() {
    super.initState();
    checkLogin();
    _cntSearch = TextEditingController();
  }

  @override
  void dispose() {
    _cntSearch.dispose();
    super.dispose();
  }

  Future<bool> checkLogin() async {
    final isl = await Api.isLogin();
    if (!isl) Routes.router.navigateTo(context, "/login");
    return isl;
  }

  Future<Widget> search() async {
    try {
      final list = await Api.search(query);
      books = list.map((e) => BookCard(e)).toList();

      return Expanded(
          child: GridView.builder(
        itemCount: books.length,
        gridDelegate: SliverGridDelegateWithMaxCrossAxisExtent(
          maxCrossAxisExtent: 150,
          childAspectRatio: 185 / 278,
        ),
        itemBuilder: (context, index) => books[index],
        padding: const EdgeInsets.all(2.0),
      ));
    } on NotLoginException catch (e) {
      Routes.router.navigateTo(context, "/login");
      return Container();
    } catch (err) {
      return Text(err.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(centerTitle: true, title: Text("Поиск")),
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
              margin: EdgeInsets.all(10),
              width: 700,
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.only(topLeft: Radius.circular(10), topRight: Radius.circular(10), bottomLeft: Radius.circular(10), bottomRight: Radius.circular(10)),
                boxShadow: [BoxShadow(color: Colors.black, spreadRadius: 0.5, blurRadius: 7, offset: Offset(1, 1))],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(children: [
                    Expanded(
                        child: TextField(
                      controller: _cntSearch,
                      decoration: const InputDecoration(
                        contentPadding: EdgeInsets.fromLTRB(10, 0, 10, 0),
                        hintText: "Введите название книги или имя автора",
                      ),
                      onSubmitted: (value) {
                        query = value;
                        setState(() {});
                      },
                    )),
                    FloatingActionButton(
                      onPressed: () {
                        query = _cntSearch.text;
                        setState(() {});
                      },
                      child: Icon(Icons.search),
                    )
                  ]),
                  SizedBox(height: 20),
                  if (query.isNotEmpty)
                    FutureBuilder<Widget>(
                      future: search(),
                      builder: (context, snapshot) {
                        if (snapshot.hasData) {
                          return snapshot.data!;
                        } else {
                          return Center(child: CircularProgressIndicator());
                        }
                      },
                    ),
                ],
              ),
            ))));
  }
}
