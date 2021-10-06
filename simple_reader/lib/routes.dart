import 'package:fluro/fluro.dart';
import 'package:simple_reader/pages/home/home_page.dart';
import 'package:simple_reader/pages/home/not_found_page.dart';
import 'package:simple_reader/pages/login/login_page.dart';
import 'package:simple_reader/pages/login/register_page.dart';
import 'package:simple_reader/pages/search/search_page.dart';
import 'package:simple_reader/pages/upload/upload_page.dart';

class Routes {
  static late FluroRouter router;

  static FluroRouter getRouter() {
    FluroRouter router = FluroRouter();
    Routes.configureRoutes(router);
    Routes.router = router;
    return router;
  }

  static void configureRoutes(FluroRouter router) {
    router.notFoundHandler = Handler(handlerFunc: (context, params) => NotFoundPage());

    router.define("/", handler: Handler(handlerFunc: (context, params) => HomePage()), transitionType: TransitionType.fadeIn);

    router.define("/upload", handler: Handler(handlerFunc: (context, params) => UploadPage()), transitionType: TransitionType.fadeIn);

    router.define("/login", handler: Handler(handlerFunc: (context, params) => LoginPage()), transitionType: TransitionType.fadeIn);

    router.define("/register/:hash", handler: Handler(handlerFunc: (context, params) => RegisterPage(params["hash"]?.first ?? "")), transitionType: TransitionType.fadeIn);

    router.define("/search", handler: Handler(handlerFunc: (context, params) => SearchPage()), transitionType: TransitionType.fadeIn);
  }
}
