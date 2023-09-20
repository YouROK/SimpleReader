package utils

import "strings"

func LineGenres(Genre []string) string {
	return strings.Join(GetGenres(Genre), ", ")
}

func GetGenres(Genre []string) []string {
	gen := make([]string, len(Genre))
	for i, k := range Genre {
		switch k {
		case "sf_history":
			gen[i] = "Альтернативная история"
		case "sf_action":
			gen[i] = "Боевая фантастика"
		case "sf_epic":
			gen[i] = "Эпическая фантастика"
		case "sf_heroic":
			gen[i] = "Героическая фантастика"
		case "sf_detective":
			gen[i] = "Детективная фантастика"
		case "sf_cyberpunk":
			gen[i] = "Киберпанк"
		case "sf_space":
			gen[i] = "Космическая фантастика"
		case "sf_social":
			gen[i] = "Социально-психологическая фантастика"
		case "sf_horror":
			gen[i] = "Ужасы и Мистика"
		case "sf_humor":
			gen[i] = "Юмористическая фантастика"
		case "sf_fantasy":
			gen[i] = "Фэнтези"
		case "sf":
			gen[i] = "Научная Фантастика "
		case "det_classic":
			gen[i] = "Классический детектив"
		case "det_police":
			gen[i] = "Полицейский детектив"
		case "det_action":
			gen[i] = "Боевик"
		case "det_irony":
			gen[i] = "Иронический детектив"
		case "det_history":
			gen[i] = "Исторический детектив"
		case "det_espionage":
			gen[i] = "Шпионский детектив"
		case "det_crime":
			gen[i] = "Криминальный детектив"
		case "det_political":
			gen[i] = "Политический детектив"
		case "det_maniac":
			gen[i] = "Маньяки"
		case "det_hard":
			gen[i] = "Крутой детектив"
		case "thriller":
			gen[i] = "Триллер"
		case "detective":
			gen[i] = "Детектив (не относящийся в прочие категории)."
		case "prose_classic":
			gen[i] = "Классическая проза"
		case "prose_history":
			gen[i] = "Историческая проза"
		case "prose_contemporary":
			gen[i] = "Современная проза"
		case "prose_counter":
			gen[i] = "Контркультура"
		case "prose_rus_classic":
			gen[i] = "Русская классическая проза"
		case "prose_su_classics":
			gen[i] = "Советская классическая проза"
		case "love_contemporary":
			gen[i] = "Современные любовные романы"
		case "love_history":
			gen[i] = "Исторические любовные романы"
		case "love_detective":
			gen[i] = "Остросюжетные любовные романы"
		case "love_short":
			gen[i] = "Короткие любовные романы"
		case "love_erotica":
			gen[i] = "Эротика"
		case "adv_western":
			gen[i] = "Вестерн"
		case "adv_history":
			gen[i] = "Исторические приключения"
		case "adv_indian":
			gen[i] = "Приключения про индейцев"
		case "adv_maritime":
			gen[i] = "Морские приключения"
		case "adv_geo":
			gen[i] = "Путешествия и география"
		case "adv_animal":
			gen[i] = "Природа и животные"
		case "adventure":
			gen[i] = "Прочие приключения (то, что не вошло в другие категории)"
		case "child_tale":
			gen[i] = "Сказка"
		case "child_verse":
			gen[i] = "Детские стихи"
		case "child_prose":
			gen[i] = "Детскиая проза"
		case "child_sf":
			gen[i] = "Детская фантастика"
		case "child_det":
			gen[i] = "Детские остросюжетные"
		case "child_adv":
			gen[i] = "Детские приключения"
		case "child_education":
			gen[i] = "Детская образовательная литература"
		case "children":
			gen[i] = "Прочая детская литература (то, что не вошло в другие категории)"
		case "poetry":
			gen[i] = "Поэзия"
		case "dramaturgy":
			gen[i] = "Драматургия"
		case "antique_ant":
			gen[i] = "Античная литература"
		case "antique_european":
			gen[i] = "Европейская старинная литература"
		case "antique_russian":
			gen[i] = "Древнерусская литература"
		case "antique_east":
			gen[i] = "Древневосточная литература"
		case "antique_myths":
			gen[i] = "Мифы. Легенды. Эпос"
		case "antique":
			gen[i] = "Прочая старинная литература (то, что не вошло в другие категории)"
		case "sci_history":
			gen[i] = "История"
		case "sci_psychology":
			gen[i] = "Психология"
		case "sci_culture":
			gen[i] = "Культурология"
		case "sci_religion":
			gen[i] = "Религиоведение"
		case "sci_philosophy":
			gen[i] = "Философия"
		case "sci_politics":
			gen[i] = "Политика"
		case "sci_business":
			gen[i] = "Деловая литература"
		case "sci_juris":
			gen[i] = "Юриспруденция"
		case "sci_linguistic":
			gen[i] = "Языкознание"
		case "sci_medicine":
			gen[i] = "Медицина"
		case "sci_phys":
			gen[i] = "Физика"
		case "sci_math":
			gen[i] = "Математика"
		case "sci_chem":
			gen[i] = "Химия"
		case "sci_biology":
			gen[i] = "Биология"
		case "sci_tech":
			gen[i] = "Технические науки"
		case "science":
			gen[i] = "Прочая научная литература (то, что не вошло в другие категории)"
		case "comp_www":
			gen[i] = "Интернет"
		case "comp_programming":
			gen[i] = "Программирование"
		case "comp_hard":
			gen[i] = "Компьютерное \" железо\" (аппаратное обеспечение)"
		case "comp_soft":
			gen[i] = "Программы"
		case "comp_db":
			gen[i] = "Базы данных"
		case "comp_osnet":
			gen[i] = "ОС и Сети"
		case "computers":
			gen[i] = "Прочая околокомпьтерная литература (то, что не вошло в другие категории)"
		case "ref_encyc":
			gen[i] = "Энциклопедии"
		case "ref_dict":
			gen[i] = "Словари"
		case "ref_ref":
			gen[i] = "Справочники"
		case "ref_guide":
			gen[i] = "Руководства"
		case "reference":
			gen[i] = "Прочая справочная литература (то, что не вошло в другие категории)"
		case "nonf_biography":
			gen[i] = "Биографии и Мемуары"
		case "nonf_publicism":
			gen[i] = "Публицистика"
		case "nonf_criticism":
			gen[i] = "Критика"
		case "design":
			gen[i] = "Искусство и Дизайн"
		case "nonfiction":
			gen[i] = "Прочая документальная литература (то, что не вошло в другие категории)"
		case "religion_rel":
			gen[i] = "Религия"
		case "religion_esoterics":
			gen[i] = "Эзотерика"
		case "religion_self":
			gen[i] = "Самосовершенствование"
		case "religion":
			gen[i] = "Прочая религионая литература (то, что не вошло в другие категории)"
		case "humor_anecdote":
			gen[i] = "Анекдоты"
		case "humor_prose":
			gen[i] = "Юмористическая проза"
		case "humor_verse":
			gen[i] = "Юмористические стихи"
		case "humor":
			gen[i] = "Прочий юмор (то, что не вошло в другие категории)"
		case "home_cooking":
			gen[i] = "Кулинария"
		case "home_pets":
			gen[i] = "Домашние животные"
		case "home_crafts":
			gen[i] = "Хобби и ремесла"
		case "home_entertain":
			gen[i] = "Развлечения"
		case "home_health":
			gen[i] = "Здоровье"
		case "home_garden":
			gen[i] = "Сад и огород"
		case "home_diy":
			gen[i] = "Сделай сам"
		case "home_sport":
			gen[i] = "Спорт"
		case "home_sex":
			gen[i] = "Эротика, Секс"
		case "home":
			gen[i] = "Прочиее домоводство (то, что не вошло в другие категории)"
		default:
			gen[i] = k
		}
	}
	return gen
}
