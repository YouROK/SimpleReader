import 'package:flutter/material.dart';
import 'package:flutter_settings_screens/flutter_settings_screens.dart';
import 'package:simple_reader/routes.dart';
import 'package:simple_reader/style.dart';

Future<void> main() async {
  await Settings.init();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Simple Reader',
      onGenerateRoute: Routes.getRouter().generator,
      theme: appTheme(),
      initialRoute: '/upload',
    );
  }
}

/*class Homepage extends StatefulWidget {
  @override
  _HomepageState createState() => _HomepageState();
}

class _HomepageState extends State<Homepage> {
  final PageController pageController = PageController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: GestureDetector(
          onPanUpdate: (details) {
            _onScroll(details.delta.dx * -1);
          },
          child: Listener(
            onPointerSignal: (pointerSignal) {
              if (pointerSignal is PointerScrollEvent) _onScroll(pointerSignal.scrollDelta.dy);
            },
            child: PageView.builder(
              scrollDirection: Axis.horizontal,
              controller: pageController,
              physics: const NeverScrollableScrollPhysics(),
              itemBuilder: (context, index) => Center(child: Text(index.toString())),
            ),
          ),
        ),
      ),
    );
  }

  bool pageIsScrolling = false;
  void _onScroll(double offset) {
    if (pageIsScrolling == false) {
      pageIsScrolling = true;
      if (offset > 0) {
        pageController.nextPage(duration: const Duration(milliseconds: 300), curve: Curves.fastOutSlowIn).then((value) => pageIsScrolling = false);
      } else {
        pageController.previousPage(duration: const Duration(milliseconds: 300), curve: Curves.fastOutSlowIn).then((value) => pageIsScrolling = false);
      }
    }
  }
}
*/
