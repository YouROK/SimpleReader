class BookInfo {
  String BookHash = "";
  int LastReadSection = 0;
  String LastRead = "";

  BookInfo(this.BookHash, this.LastReadSection, this.LastRead);

  BookInfo.fromJson(Map<String, dynamic> json)
      : BookHash = json['BookHash'],
        LastReadSection = json['LastReadSection'],
        LastRead = json['LastRead'];
}
