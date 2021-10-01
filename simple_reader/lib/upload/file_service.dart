import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:http/http.dart' as http;
import 'package:path/path.dart' as fileUtil;
import 'package:simple_reader/api/api.dart';

typedef void OnUploadProgressCallback(int sentBytes, int totalBytes);

class FileService {
  static HttpClient getHttpClient() {
    HttpClient httpClient = HttpClient()
      ..connectionTimeout = const Duration(seconds: 10)
      ..badCertificateCallback = ((X509Certificate cert, String host, int port) => true);
    return httpClient;
  }

  static Future<String> fileUploadMultipart({required String url, required String file, required OnUploadProgressCallback onUploadProgress}) async {
    final httpClient = getHttpClient();
    final request = await httpClient.postUrl(Uri.parse(url));
    int byteCount = 0;
    var multipart = await http.MultipartFile.fromPath(fileUtil.basename(file), file);
    var requestMultipart = http.MultipartRequest("POST", Uri.parse(url));
    requestMultipart.files.add(multipart);
    var msStream = requestMultipart.finalize();
    var totalByteLength = requestMultipart.contentLength;
    request.contentLength = totalByteLength;
    final contentType = (requestMultipart.headers[HttpHeaders.contentTypeHeader]) ?? "";
    Api.session.headers.forEach((key, value) {
      request.headers.set(key, value);
    });
    request.headers.set(HttpHeaders.contentTypeHeader, contentType);
    Stream<List<int>> streamUpload = msStream.transform(
      StreamTransformer.fromHandlers(
        handleData: (data, sink) {
          sink.add(data);
          byteCount += data.length;
          onUploadProgress(byteCount, totalByteLength);
        },
        handleError: (error, stack, sink) {
          throw error;
        },
        handleDone: (sink) {
          sink.close();
        },
      ),
    );

    await request.addStream(streamUpload);

    final httpResponse = await request.close();
    var statusCode = httpResponse.statusCode;

    if (statusCode ~/ 100 != 2) {
      throw Exception('Error uploading file, Status code: ${httpResponse.statusCode}');
    } else {
      return await readResponseAsString(httpResponse);
    }
  }

  static Future<String> readResponseAsString(HttpClientResponse response) {
    var completer = new Completer<String>();
    var contents = new StringBuffer();
    response.transform(utf8.decoder).listen((String data) {
      contents.write(data);
    }, onDone: () => completer.complete(contents.toString()));
    return completer.future;
  }
}
