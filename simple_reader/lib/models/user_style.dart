class UserStyle {
  int PageSlideEffect = 0;
  double TextSize = 1.0;
  double TextIndent = 1.0;
  double ParagraphIndent = 1.5;
  int TextBright = 100;
  bool DayTheme = true;
  String FontName = "";

  UserStyle();

  UserStyle.fromJson(Map<String, dynamic> json)
      : PageSlideEffect = json['page_slide_effect'],
        TextSize = json['text_size'],
        TextIndent = json['text_indent'],
        ParagraphIndent = json['paragraph_indent'],
        TextBright = json['text_bright'],
        DayTheme = json['day_theme'],
        FontName = json['font_name'];
}
