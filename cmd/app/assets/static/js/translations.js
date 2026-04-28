const lang = document.documentElement.lang || 'pt';
const translations = {
    special: lang === 'en' ? 'Special' :
             lang === 'es' ? 'Especial' :
             lang === 'fr' ? 'Spécial' :
             lang === 'de' ? 'Special' :
             lang === 'it' ? 'Speciale' :
             lang === 'ja' ? '特別版' :
             lang === 'zh' ? '特别篇' :
             lang === 'ru' ? 'Специальный' :
             lang === 'ar' ? 'خاص' :
             lang === 'hi' ? 'विशेष' :
             lang === 'ko' ? '스페셜' :
             'Especial',
    season: lang === 'en' ? 'Season' :
            lang === 'es' ? 'Temporada' :
            lang === 'fr' ? 'Saison' :
            lang === 'de' ? 'Staffel' :
            lang === 'it' ? 'Stagione' :
            lang === 'ja' ? 'シーズン' :
            lang === 'zh' ? '季' :
            lang === 'ru' ? 'Сезон' :
            lang === 'ar' ? 'الموسم' :
            lang === 'hi' ? 'सीज़न' :
            lang === 'ko' ? '시즌' :
            'Temporada',
    soon: lang === 'en' ? 'Soon' :
          lang === 'es' ? 'Próximamente' :
          lang === 'fr' ? 'Bientôt' :
          lang === 'de' ? 'Demnächst' :
          lang === 'it' ? 'Prossimamente' :
          lang === 'ja' ? '近日公開' :
          lang === 'zh' ? '即将推出' :
          lang === 'ru' ? 'Скоро' :
          lang === 'ar' ? 'قريباً' :
          lang === 'hi' ? 'जल्द ही' :
          lang === 'ko' ? '곧' :
          'Em breve',
    noEpisodes: lang === 'en' ? 'No episodes found.' :
                lang === 'es' ? 'No se encontraron episodios.' :
                lang === 'fr' ? 'Aucun épisode trouvé.' :
                lang === 'de' ? 'Keine Episoden gefunden.' :
                lang === 'it' ? 'Nessun episodio trovato.' :
                lang === 'ja' ? 'エピソードが見つかりません。' :
                lang === 'zh' ? '未找到剧集。' :
                lang === 'ru' ? 'Эпизоды не найдены.' :
                lang === 'ar' ? 'لم يتم العثور على حلقات.' :
                lang === 'hi' ? 'कोई एपिसोड नहीं मिला।' :
                lang === 'ko' ? '에피소드를 찾을 수 없습니다.' :
                'Nenhum episódio encontrado.',
    error: lang === 'en' ? 'Error loading episodes.' :
           lang === 'es' ? 'Error al cargar episodios.' :
           lang === 'fr' ? 'Erreur lors du chargement des épisodes.' :
           lang === 'de' ? 'Fehler beim Laden der Episoden.' :
           lang === 'it' ? 'Errore nel caricamento degli episodi.' :
           lang === 'ja' ? 'エピソードの読み込みエラー。' :
           lang === 'zh' ? '加载剧集时出错。' :
           lang === 'ru' ? 'Ошибка загрузки эпизодов.' :
           lang === 'ar' ? 'خطأ في تحميل الحلقات.' :
           lang === 'hi' ? 'एपिसोड लोड करने में त्रुटि।' :
           lang === 'ko' ? '에피소드를 로드하는 중 오류가 발생했습니다.' :
           'Erro ao carregar episódios.',
    loading: lang === 'en' ? 'Loading...' :
             lang === 'es' ? 'Cargando...' :
             lang === 'fr' ? 'Chargement...' :
             lang === 'de' ? 'Laden...' :
             lang === 'it' ? 'Caricamento...' :
             lang === 'ja' ? '読み込み中...' :
             lang === 'zh' ? '加载中...' :
             lang === 'ru' ? 'Загрузка...' :
             lang === 'ar' ? 'جار التحميل...' :
             lang === 'hi' ? 'लोड हो रहा है...' :
             lang === 'ko' ? '로딩 중...' :
             'Carregando...',
    noEpisodesInSeason: lang === 'en' ? 'No episodes in this season.' :
                        lang === 'es' ? 'No hay episodios en esta temporada.' :
                        lang === 'fr' ? 'Aucun épisode dans cette saison.' :
                        lang === 'de' ? 'Keine Episoden in dieser Staffel.' :
                        lang === 'it' ? 'Nessun episodio in questa stagione.' :
                        lang === 'ja' ? 'このシーズンにはエピソードがありません。' :
                        lang === 'zh' ? '本季暂无剧集。' :
                        lang === 'ru' ? 'В этом сезоне нет эпизодов.' :
                        lang === 'ar' ? 'لا توجد حلقات في هذا الموسم.' :
                        lang === 'hi' ? 'इस सीज़न में कोई एपिसोड नहीं है।' :
                        lang === 'ko' ? '이 시즌에는 에피소드가 없습니다.' :
                        'Nenhum episódio nesta temporada.',
    watchedEpisodes: lang === 'en' ? 'episodes watched' :
                     lang === 'es' ? 'episodios vistos' :
                     lang === 'fr' ? 'épisodes vus' :
                     lang === 'de' ? 'angesehene Episoden' :
                     lang === 'it' ? 'episodi visti' :
                     lang === 'ja' ? '視聴済みエピソード' :
                     lang === 'zh' ? '已观看剧集' :
                     lang === 'ru' ? 'просмотренных эпизодов' :
                     lang === 'ar' ? 'حلقات تم مشاهدتها' :
                     lang === 'hi' ? 'देखे गए एपिसोड' :
                     lang === 'ko' ? '시청한 에피소드' :
                     'episódios assistidos',
	noStreams: lang === 'en' ? 'No streaming services found for this device.' :
				lang === 'es' ? 'No se encontraron servicios de streaming para este dispositivo.' :
				lang === 'fr' ? 'Aucun service de streaming trouvé pour cet appareil.' :
				lang === 'de' ? 'Keine Streaming-Dienste für dieses Gerät gefunden.' :
				lang === 'it' ? 'Nessun servizio di streaming trovato per questo dispositivo.' :
				lang === 'ja' ? 'このデバイスに対応するストリーミングサービスが見つかりません。' :
				lang === 'zh' ? '未找到适用于此设备的流媒体服务。' :
				lang === 'ru' ? 'Стриминговые сервисы для этого устройства не найдены.' :
				lang === 'ar' ? 'لم يتم العثور على خدمات بث لهذا الجهاز.' :
				lang === 'hi' ? 'इस डिवाइस के लिए कोई स्ट्रीमिंग सेवा नहीं मिली।' :
				lang === 'ko' ? '이 기기에서 사용할 수 있는 스트리밍 서비스를 찾을 수 없습니다.' :
				'Nenhum serviço de streaming encontrado para este dispositivo.',
	available: lang === 'en' ? 'Available on:' :
			   lang === 'es' ? 'Disponible en:' :
			   lang === 'fr' ? 'Disponible sur :' :
			   lang === 'de' ? 'Verfügbar auf:' :
			   lang === 'it' ? 'Disponibile su:' :
			   lang === 'ja' ? '視聴可能:' :
			   lang === 'zh' ? '可在以下平台观看：' :
			   lang === 'ru' ? 'Доступно на:' :
			   lang === 'ar' ? 'متاح على:' :
			   lang === 'hi' ? 'उपलब्ध है:' :
			   lang === 'ko' ? '시청 가능:' :
			   'Disponível em:',
	watchWeb: lang === 'en' ? 'Watch on web' :
			  lang === 'es' ? 'Ver en web' :
			  lang === 'fr' ? 'Regarder sur le web' :
			  lang === 'de' ? 'Im Web ansehen' :
			  lang === 'it' ? 'Guarda sul web' :
			  lang === 'ja' ? 'ウェブで視聴' :
			  lang === 'zh' ? '在网页上观看' :
			  lang === 'ru' ? 'Смотреть в вебе' :
			  lang === 'ar' ? 'مشاهدة على الويب' :
			  lang === 'hi' ? 'वेब पर देखें' :
			  lang === 'ko' ? '웹에서 시청' :
			  'Assistir na Web',
	watchAndroid: lang === 'en' ? 'Watch on Android' :
				  lang === 'es' ? 'Ver en Android' :
				  lang === 'fr' ? 'Regarder sur Android' :
				  lang === 'de' ? 'Auf Android ansehen' :
				  lang === 'it' ? 'Guarda su Android' :
				  lang === 'ja' ? 'Androidで視聴' :
				  lang === 'zh' ? '在Android上观看' :
				  lang === 'ru' ? 'Смотреть на Android' :
				  lang === 'ar' ? 'مشاهدة على أندرويد' :
				  lang === 'hi' ? 'Android पर देखें' :
				  lang === 'ko' ? 'Android에서 시청' :
				  'Assistir no Android',
	watchIOS: lang === 'en' ? 'Watch on iOS' :
			  lang === 'es' ? 'Ver en iOS' :
			  lang === 'fr' ? 'Regarder sur iOS' :
			  lang === 'de' ? 'Auf iOS ansehen' :
			  lang === 'it' ? 'Guarda su iOS' :
			  lang === 'ja' ? 'iOSで視聴' :
			  lang === 'zh' ? '在iOS上观看' :
			  lang === 'ru' ? 'Смотреть на iOS' :
			  lang === 'ar' ? 'مشاهدة على iOS' :
			  lang === 'hi' ? 'iOS पर देखें' :
			  lang === 'ko' ? 'iOS에서 시청' :
			  'Assistir no iOS',
};