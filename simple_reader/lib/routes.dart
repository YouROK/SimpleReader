import 'package:fluro/fluro.dart';
import 'package:simple_reader/home/home_page.dart';

import 'home/not_found_page.dart';
import 'login/login_page.dart';
import 'login/register_page.dart';

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

    router.define("/", handler: Handler(handlerFunc: (context, params) => HomePage()));
    router.define("/login", handler: Handler(handlerFunc: (context, params) => LoginPage()));
    router.define("/register/:hash", handler: Handler(handlerFunc: (context, params) => RegisterPage(params["hash"]?.first ?? "")));
  }
}
