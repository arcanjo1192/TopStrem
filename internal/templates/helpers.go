package templates

// Idiomas suportados (códigos ISO 639-1)
var supportedLangs = []string{"pt", "en", "es", "fr", "de", "it", "ja", "zh", "ru", "ar", "hi", "ko"}

// Tabelas de tradução
var translations = map[string]map[string]string{
    // Catálogo – Filmes
    "catalog_title_movie": {
        "pt": "Filmes mais assistidos",
        "en": "Most Watched Movies",
        "es": "Películas más vistas",
        "fr": "Films les plus regardés",
        "de": "Meistgesehene Filme",
        "it": "Film più visti",
        "ja": "最も視聴された映画",
        "zh": "观看最多的电影",
        "ru": "Самые просматриваемые фильмы",
        "ar": "الأفلام الأكثر مشاهدة",
        "hi": "सबसे ज्यादा देखी जाने वाली फिल्में",
        "ko": "가장 많이 시청된 영화",
    },
    "catalog_title_series": {
        "pt": "Séries mais assistidas",
        "en": "Most Watched Series",
        "es": "Series más vistas",
        "fr": "Séries les plus regardées",
        "de": "Meistgesehene Serien",
        "it": "Serie più viste",
        "ja": "最も視聴されたシリーズ",
        "zh": "观看最多的剧集",
        "ru": "Самые просматриваемые сериалы",
        "ar": "المسلسلات الأكثر مشاهدة",
        "hi": "सबसे ज्यादा देखी जाने वाली सीरीज",
        "ko": "가장 많이 시청된 시리즈",
    },
    "catalog_desc_movie": {
        "pt": "Explore os filmes mais assistidos, top 10 filmes e filmes em alta no momento. Sinopses, elenco, trailers e muito mais no TopStrem.",
        "en": "Explore the most watched movies, top 10 movies and trending movies right now. Synopses, cast, trailers and more on TopStrem.",
        "es": "Explora las películas más vistas, las 10 mejores películas y las tendencias del momento. Sinopsis, reparto, tráilers y más en TopStrem.",
        "fr": "Découvrez les films les plus regardés, le top 10 des films et les films tendances du moment. Synopsis, casting, bandes-annonces et plus sur TopStrem.",
        "de": "Entdecken Sie die meistgesehenen Filme, die Top-10-Filme und die angesagtesten Filme. Zusammenfassungen, Besetzung, Trailer und mehr auf TopStrem.",
        "it": "Scopri i film più visti, la top 10 dei film e i film di tendenza. Sinossi, cast, trailer e altro su TopStrem.",
        "ja": "最も視聴された映画、トップ10映画、話題の映画をチェック。あらすじ、キャスト、予告編など。TopStremで。",
        "zh": "探索观看最多的电影、十大电影和热门电影。剧情简介、演员表、预告片等尽在TopStrem。",
        "ru": "Откройте для себя самые просматриваемые фильмы, топ-10 фильмов и актуальные фильмы. Синопсис, актёры, трейлеры и многое другое на TopStrem.",
        "ar": "استكشف الأفلام الأكثر مشاهدة، أفضل 10 أفلام، والأفلام الرائجة حالياً. ملخصات، طاقم التمثيل، إعلانات تشويقية والمزيد على TopStrem.",
        "hi": "सबसे ज्यादा देखी जाने वाली फिल्में, टॉप 10 फिल्में और ट्रेंडिंग फिल्में देखें। सारांश, कलाकार, ट्रेलर और भी बहुत कुछ TopStrem पर।",
        "ko": "가장 많이 시청된 영화, 상위 10개 영화 및 트렌딩 영화를 찾아보세요. 시놉시스, 출연진, 예고편 등을 TopStrem에서 확인하세요.",
    },
    "catalog_desc_series": {
        "pt": "Confira as séries mais assistidas, top 10 séries e séries em alta no momento. Descubra o que está bombando na Netflix, Max, Prime Video e outras plataformas.",
        "en": "Check out the most watched series, top 10 series, and trending series right now. Discover what's hot on Netflix, Max, Prime Video and more.",
        "es": "Consulta las series más vistas, las 10 mejores series y las series de moda. Descubre lo que está arrasando en Netflix, Max, Prime Video y más.",
        "fr": "Découvrez les séries les plus regardées, le top 10 des séries et les séries tendances. Découvrez ce qui fait fureur sur Netflix, Max, Prime Video et plus.",
        "de": "Erkunden Sie die meistgesehenen Serien, die Top-10-Serien und die angesagtesten Serien. Erfahren Sie, was auf Netflix, Max, Prime Video und anderen Plattformen angesagt ist.",
        "it": "Scopri le serie più viste, le top 10 serie e le serie di tendenza. Scopri cosa sta spopolando su Netflix, Max, Prime Video e altri.",
        "ja": "最も視聴されたシリーズ、トップ10シリーズ、話題のシリーズをチェック。Netflix、Max、Prime Videoなどで何が人気か見つけましょう。",
        "zh": "查看观看最多的剧集、十大剧集和热门剧集。发现Netflix、Max、Prime Video等平台上的热门内容。",
        "ru": "Откройте для себя самые просматриваемые сериалы, топ-10 сериалов и актуальные сериалы. Узнайте, что сейчас популярно на Netflix, Max, Prime Video и других платформах.",
        "ar": "تحقق من المسلسلات الأكثر مشاهدة، أفضل 10 مسلسلات، والمسلسلات الرائجة حالياً. اكتشف ما هو رائج على Netflix وMax وPrime Video وخدمات أخرى.",
        "hi": "सबसे ज्यादा देखी जाने वाली सीरीज, टॉप 10 सीरीज और ट्रेंडिंग सीरीज देखें। Netflix, Max, Prime Video और अन्य प्लेटफार्मों पर क्या चल रहा है, जानें।",
        "ko": "가장 많이 시청된 시리즈, 상위 10개 시리즈 및 트렌딩 시리즈를 확인하세요. Netflix, Max, Prime Video 등에서 무엇이 인기 있는지 알아보세요.",
    },
    "nav_movie": {
        "pt": "Filmes Populares",
        "en": "Popular Movies",
        "es": "Películas Populares",
        "fr": "Films Populaires",
        "de": "Beliebte Filme",
        "it": "Film Popolari",
        "ja": "人気映画",
        "zh": "热门电影",
        "ru": "Популярные фильмы",
        "ar": "أفلام شائعة",
        "hi": "लोकप्रिय फिल्में",
        "ko": "인기 영화",
    },
    "nav_series": {
        "pt": "Séries Populares",
        "en": "Popular Series",
        "es": "Series Populares",
        "fr": "Séries Populaires",
        "de": "Beliebte Serien",
        "it": "Serie Popolari",
        "ja": "人気シリーズ",
        "zh": "热门剧集",
        "ru": "Популярные сериалы",
        "ar": "مسلسلات شائعة",
        "hi": "लोकप्रिय सीरीज",
        "ko": "인기 시리즈",
    },
    "cast_label": {
        "pt": "Elenco:",
        "en": "Cast:",
        "es": "Reparto:",
        "fr": "Distribution:",
        "de": "Besetzung:",
        "it": "Cast:",
        "ja": "キャスト:",
        "zh": "演员阵容:",
        "ru": "В ролях:",
        "ar": "طاقم التمثيل:",
        "hi": "कलाकार:",
        "ko": "출연:",
    },
    "director_label": {
        "pt": "Diretor:",
        "en": "Director:",
        "es": "Director:",
        "fr": "Réalisateur:",
        "de": "Regisseur:",
        "it": "Regista:",
        "ja": "監督:",
        "zh": "导演:",
        "ru": "Режиссёр:",
        "ar": "المخرج:",
        "hi": "निर्देशक:",
        "ko": "감독:",
    },
    "country_label": {
        "pt": "País:",
        "en": "Country:",
        "es": "País:",
        "fr": "Pays:",
        "de": "Land:",
        "it": "Paese:",
        "ja": "国:",
        "zh": "国家:",
        "ru": "Страна:",
        "ar": "الدولة:",
        "hi": "देश:",
        "ko": "국가:",
    },
    "awards_label": {
        "pt": "Prêmios:",
        "en": "Awards:",
        "es": "Premios:",
        "fr": "Récompenses:",
        "de": "Auszeichnungen:",
        "it": "Premi:",
        "ja": "受賞歴:",
        "zh": "奖项:",
        "ru": "Награды:",
        "ar": "الجوائز:",
        "hi": "पुरस्कार:",
        "ko": "수상:",
    },
    "back_button": {
        "pt": "Voltar",
        "en": "Back",
        "es": "Volver",
        "fr": "Retour",
        "de": "Zurück",
        "it": "Indietro",
        "ja": "戻る",
        "zh": "返回",
        "ru": "Назад",
        "ar": "رجوع",
        "hi": "पीछे",
        "ko": "뒤로",
    },
    "trailer_button": {
        "pt": "Trailer",
        "en": "Trailer",
        "es": "Tráiler",
        "fr": "Bande-annonce",
        "de": "Trailer",
        "it": "Trailer",
        "ja": "予告編",
        "zh": "预告片",
        "ru": "Трейлер",
        "ar": "مقطع دعائي",
        "hi": "ट्रेलर",
        "ko": "예고편",
    },
    "watch_button": {
        "pt": "Assistir agora (Stremio)",
        "en": "Watch now (Stremio)",
        "es": "Ver ahora (Stremio)",
        "fr": "Regarder maintenant (Stremio)",
        "de": "Jetzt ansehen (Stremio)",
        "it": "Guarda ora (Stremio)",
        "ja": "今すぐ見る (Stremio)",
        "zh": "立即观看 (Stremio)",
        "ru": "Смотреть сейчас (Stremio)",
        "ar": "شاهد الآن (Stremio)",
        "hi": "अभी देखें (Stremio)",
        "ko": "지금 시청 (Stremio)",
    },
    "episodes_button": {
        "pt": "Episódios",
        "en": "Episodes",
        "es": "Episodios",
        "fr": "Épisodes",
        "de": "Episoden",
        "it": "Episodi",
        "ja": "エピソード",
        "zh": "剧集",
        "ru": "Эпизоды",
        "ar": "حلقات",
        "hi": "एपिसोड",
        "ko": "에피소드",
    },
    "sidebar_title": {
        "pt": "Temporadas e Episódios",
        "en": "Seasons and Episodes",
        "es": "Temporadas y Episodios",
        "fr": "Saisons et épisodes",
        "de": "Staffeln und Episoden",
        "it": "Stagioni ed episodi",
        "ja": "シーズンとエピソード",
        "zh": "季与剧集",
        "ru": "Сезоны и эпизоды",
        "ar": "مواسم وحلقات",
        "hi": "सीज़न और एपिसोड",
        "ko": "시즌 및 에피소드",
    },
    "loading": {
        "pt": "Carregando...",
        "en": "Loading...",
        "es": "Cargando...",
        "fr": "Chargement...",
        "de": "Laden...",
        "it": "Caricamento...",
        "ja": "読み込み中...",
        "zh": "加载中...",
        "ru": "Загрузка...",
        "ar": "جار التحميل...",
        "hi": "लोड हो रहा है...",
        "ko": "로딩 중...",
    },
    "notification": {
        "pt": "Notificações",
        "en": "Notifications",
        "es": "Notificaciones",
        "fr": "Notifications",
        "de": "Benachrichtigungen",
        "it": "Notifiche",
        "ja": "通知",
        "zh": "通知",
        "ru": "Уведомления",
        "ar": "الإشعارات",
        "hi": "सूचनाएं",
        "ko": "알림",
    },
    "other_platforms": {
        "pt": "Outras Plataformas",
        "en": "Other Platforms",
        "es": "Otras Plataformas",
        "fr": "Autres plateformes",
        "de": "Andere Plattformen",
        "it": "Altre piattaforme",
        "ja": "他のプラットフォーム",
        "zh": "其他平台",
        "ru": "Другие платформы",
        "ar": "منصات أخرى",
        "hi": "अन्य प्लेटफार्म",
        "ko": "기타 플랫폼",
    },
    "no_streams": {
        "pt": "Nenhum serviço de streaming encontrado para este dispositivo.",
        "en": "No streaming services found for this device.",
        "es": "No se encontraron servicios de streaming para este dispositivo.",
        "fr": "Aucun service de streaming trouvé pour cet appareil.",
        "de": "Keine Streaming-Dienste für dieses Gerät gefunden.",
        "it": "Nessun servizio di streaming trovato per questo dispositivo.",
        "ja": "このデバイスではストリーミングサービスが見つかりませんでした。",
        "zh": "未找到适用于此设备的流媒体服务。",
        "ru": "Стриминговые сервисы для этого устройства не найдены.",
        "ar": "لم يتم العثور على خدمات بث لهذا الجهاز.",
        "hi": "इस डिवाइस के लिए कोई स्ट्रीमिंग सेवा नहीं मिली।",
        "ko": "이 장치에 대한 스트리밍 서비스를 찾을 수 없습니다.",
    },
    "error_streams": {
        "pt": "Erro ao carregar opções de streaming.",
        "en": "Error loading streaming options.",
        "es": "Error al cargar las opciones de streaming.",
        "fr": "Erreur lors du chargement des options de streaming.",
        "de": "Fehler beim Laden der Streaming-Optionen.",
        "it": "Errore durante il caricamento delle opzioni di streaming.",
        "ja": "ストリーミングオプションの読み込み中にエラーが発生しました。",
        "zh": "加载流媒体选项时出错。",
        "ru": "Ошибка загрузки опций потокового вещания.",
        "ar": "خطأ في تحميل خيارات البث.",
        "hi": "स्ट्रीमिंग विकल्प लोड करने में त्रुटि।",
        "ko": "스트리밍 옵션을 로드하는 중 오류가 발생했습니다.",
    },
    "available_on": {
        "pt": "Disponível em:",
        "en": "Available on:",
        "es": "Disponible en:",
        "fr": "Disponible sur:",
        "de": "Verfügbar auf:",
        "it": "Disponibile su:",
        "ja": "対応プラットフォーム:",
        "zh": "可在以下平台观看:",
        "ru": "Доступно на:",
        "ar": "متاح على:",
        "hi": "उपलब्ध है:",
        "ko": "시청 가능:",
    },
    "watch_web": {
        "pt": "Assistir na Web",
        "en": "Watch on web",
        "es": "Ver en web",
        "fr": "Regarder sur le web",
        "de": "Im Web ansehen",
        "it": "Guarda sul web",
        "ja": "ウェブで視聴",
        "zh": "在网页上观看",
        "ru": "Смотреть в вебе",
        "ar": "مشاهدة على الويب",
        "hi": "वेब पर देखें",
        "ko": "웹에서 시청",
    },
    "watch_android": {
        "pt": "Assistir no Android",
        "en": "Watch on Android",
        "es": "Ver en Android",
        "fr": "Regarder sur Android",
        "de": "Auf Android ansehen",
        "it": "Guarda su Android",
        "ja": "Androidで視聴",
        "zh": "在 Android 上观看",
        "ru": "Смотреть на Android",
        "ar": "مشاهدة على أندرويد",
        "hi": "Android पर देखें",
        "ko": "Android에서 시청",
    },
    "watch_ios": {
        "pt": "Assistir no iOS",
        "en": "Watch on iOS",
        "es": "Ver en iOS",
        "fr": "Regarder sur iOS",
        "de": "Auf iOS ansehen",
        "it": "Guarda su iOS",
        "ja": "iOSで視聴",
        "zh": "在 iOS 上观看",
        "ru": "Смотреть на iOS",
        "ar": "مشاهدة على iOS",
        "hi": "iOS पर देखें",
        "ko": "iOS에서 시청",
    },
    "login_with_google": {
        "pt": "Login com Google",
        "en": "Login with Google",
        "es": "Iniciar sesión con Google",
        "fr": "Se connecter avec Google",
        "de": "Mit Google anmelden",
        "it": "Accedi con Google",
        "ja": "Googleでログイン",
        "zh": "使用 Google 登录",
        "ru": "Войти через Google",
        "ar": "تسجيل الدخول باستخدام Google",
        "hi": "Google से लॉगिन करें",
        "ko": "Google로 로그인",
    },
    "favorites_series": {
        "pt": "Séries favoritas",
        "en": "Favorite Series",
        "es": "Series favoritas",
        "fr": "Séries favorites",
        "de": "Lieblingsserien",
        "it": "Serie preferite",
        "ja": "お気に入りのシリーズ",
        "zh": "最喜欢的剧集",
        "ru": "Любимые сериалы",
        "ar": "المسلسلات المفضلة",
        "hi": "पसंदीदा सीरीज",
        "ko": "즐겨찾는 시리즈",
    },
    "favorites_movies": {
        "pt": "Filmes favoritos",
        "en": "Favorite Movies",
        "es": "Películas favoritas",
        "fr": "Films favoris",
        "de": "Lieblingsfilme",
        "it": "Film preferiti",
        "ja": "お気に入りの映画",
        "zh": "最喜欢的电影",
        "ru": "Любимые фильмы",
        "ar": "الأفلام المفضلة",
        "hi": "पसंदीदा फिल्में",
        "ko": "즐겨찾는 영화",
    },
    "logout": {
        "pt": "Sair",
        "en": "Logout",
        "es": "Cerrar sesión",
        "fr": "Déconnexion",
        "de": "Abmelden",
        "it": "Esci",
        "ja": "ログアウト",
        "zh": "退出登录",
        "ru": "Выйти",
        "ar": "تسجيل الخروج",
        "hi": "लॉग आउट",
        "ko": "로그아웃",
    },
    // NOVA CHAVE: placeholder da busca
    "search_placeholder": {
        "pt": "Pesquisar filmes ou séries...",
        "en": "Search movies or series...",
        "es": "Buscar películas o series...",
        "fr": "Rechercher des films ou séries...",
        "de": "Filme oder Serien suchen...",
        "it": "Cerca film o serie...",
        "ja": "映画またはシリーズを検索...",
        "zh": "搜索电影或剧集...",
        "ru": "Искать фильмы или сериалы...",
        "ar": "ابحث عن أفلام أو مسلسلات...",
        "hi": "फिल्में या सीरीज खोजें...",
        "ko": "영화 또는 시리즈 검색...",
    },
}

func getText(key, lang string) string {
    if translations[key] == nil {
        return key
    }
    if text, ok := translations[key][lang]; ok {
        return text
    }
    // fallback para inglês (ou português se preferir)
    if text, ok := translations[key]["en"]; ok {
        return text
    }
    return translations[key]["pt"]
}

func GetCatalogTitle(lang, catalogType string) string {
    if catalogType == "series" {
        return getText("catalog_title_series", lang)
    }
    return getText("catalog_title_movie", lang)
}

func GetCatalogDescription(lang, catalogType string) string {
    if catalogType == "series" {
        return getText("catalog_desc_series", lang)
    }
    return getText("catalog_desc_movie", lang)
}

func GetNavText(lang, catalogType string) string {
    if catalogType == "series" {
        return getText("nav_series", lang)
    }
    return getText("nav_movie", lang)
}

func GetBrandName(lang string) string {
    return "TopStrem"
}

func GetCastLabel(lang string) string {
    return getText("cast_label", lang)
}

func GetDirectorLabel(lang string) string {
    return getText("director_label", lang)
}

func GetCountryLabel(lang string) string {
    return getText("country_label", lang)
}

func GetAwardsLabel(lang string) string {
    return getText("awards_label", lang)
}

func GetBackButtonLabel(lang string) string {
    return getText("back_button", lang)
}

func GetTrailerButtonText(lang string) string {
    return getText("trailer_button", lang)
}

func GetWatchButtonText(lang string) string {
    return getText("watch_button", lang)
}

func GetEpisodesButtonText(lang string) string {
    return getText("episodes_button", lang)
}

func GetSidebarTitle(lang string) string {
    return getText("sidebar_title", lang)
}

func GetLoadingText(lang string) string {
    return getText("loading", lang)
}

func GetNotificationText(lang string) string {
    return getText("notification", lang)
}

func GetOtherPlatformsButtonText(lang string) string {
    return getText("other_platforms", lang)
}

func GetNoStreamsText(lang string) string {
    return getText("no_streams", lang)
}

func GetErrorText(lang string) string {
    return getText("error_streams", lang)
}

func GetAvailableInText(lang string) string {
    return getText("available_on", lang)
}

func GetWatchOnWebText(lang string) string {
    return getText("watch_web", lang)
}

func GetWatchOnAndroidText(lang string) string {
    return getText("watch_android", lang)
}

func GetWatchOnIOSText(lang string) string {
    return getText("watch_ios", lang)
}

func GetLoginButtonText(lang string) string {
    return getText("login_with_google", lang)
}

func GetFavoritesSeriesText(lang string) string {
    return getText("favorites_series", lang)
}

func GetFavoritesMoviesText(lang string) string {
    return getText("favorites_movies", lang)
}

func GetOGType(mediaType string) string {
    if mediaType == "series" {
        return "video.tv_show"
    }
    return "video.movie"
}

func GetLogoutText(lang string) string {
    return getText("logout", lang)
}

func GetSearchPlaceholder(lang string) string {
    return getText("search_placeholder", lang)
}