// episodes.js
document.addEventListener('DOMContentLoaded', function() {
	const sidebar = document.getElementById('episodes-sidebar');
	const openBtn = document.getElementById('episodes-button');
	if (!openBtn) return;

	const seriesId = openBtn.getAttribute('data-series-id');
	const closeBtn = document.getElementById('close-sidebar');
	const contentDiv = document.getElementById('episodes-content');
	let loaded = false;
	let fixedProgressDiv = null;

	// ========== Funções para episódios assistidos ==========
	function getWatchedEpisodesKey() {
		const email = typeof getCurrentUserEmail === 'function' ? getCurrentUserEmail() : null;
		if (!email) return null;
		return `topstrem_watched_episodes_${email}`;
	}

	function getWatchedEpisodes() {
		const key = getWatchedEpisodesKey();
		if (!key) return [];
		const stored = localStorage.getItem(key);
		return stored ? JSON.parse(stored) : [];
	}

	function saveWatchedEpisodes(episodes) {
		const key = getWatchedEpisodesKey();
		if (key) localStorage.setItem(key, JSON.stringify(episodes));
	}

	function toggleWatchedEpisode(episodeId) {
		let watched = getWatchedEpisodes();
		if (watched.includes(episodeId)) {
			watched = watched.filter(id => id !== episodeId);
		} else {
			watched.push(episodeId);
		}
		saveWatchedEpisodes(watched);
		return watched.includes(episodeId);
	}

	function isUserLoggedIn() {
		return typeof getCurrentUserEmail === 'function' && !!getCurrentUserEmail();
	}
	// =====================================================

	// Abre o episódio no Stremio
	function openEpisodeInStremio(season, episode) {
		const url = `stremio:///detail/series/${seriesId}?season=${season}&episode=${episode}`;
		window.location.href = url;
	}

	openBtn.addEventListener('click', function() {
		sidebar.classList.add('open');
		if (!loaded) {
			contentDiv.innerHTML = '<div class="loading">' + translations.loading + '</div>';
			fetch('/api/episodes/' + seriesId)
				.then(response => response.json())
				.then(seasons => {
					renderEpisodes(seasons);
					loaded = true;
				})
				.catch(err => {
					contentDiv.innerHTML = '<div class="error">' + translations.error + '</div>';
					console.error(err);
				});
		}
	});

	closeBtn.addEventListener('click', () => sidebar.classList.remove('open'));
	document.addEventListener('click', (event) => {
		if (sidebar.classList.contains('open') && !sidebar.contains(event.target) && event.target !== openBtn) {
			sidebar.classList.remove('open');
		}
	});

	function renderEpisodes(seasons) {
		if (!seasons.length) {
			contentDiv.innerHTML = '<p>' + translations.noEpisodes + '</p>';
			return;
		}

		let defaultSeasonIndex = 0;
		for (let i = 0; i < seasons.length; i++) {
			if (seasons[i].season === 1) {
				defaultSeasonIndex = i;
				break;
			}
		}

		const isLogged = isUserLoggedIn();

		// Cabeçalho com dropdown. A barra de progresso só é incluída se o usuário estiver logado
		const progressHtml = isLogged ? '<div id="fixed-progress" class="season-progress-fixed"></div>' : '';
		const headerHtml = `
			<div class="episodes-fixed-header">
				<select id="season-select" class="season-select">
					${seasons.map((season, idx) => {
						const label = season.season === 0 ? translations.special : translations.season + ' ' + season.season;
						const selected = (idx === defaultSeasonIndex) ? 'selected' : '';
						return `<option value="${idx}" ${selected}>${label}</option>`;
					}).join('')}
				</select>
				${progressHtml}
			</div>
			<div id="episodes-scrollable" class="episodes-scrollable"></div>
		`;
		contentDiv.innerHTML = headerHtml;

		const selectEl = document.getElementById('season-select');
		let currentSeasonIndex = defaultSeasonIndex;

		if (isLogged) {
			fixedProgressDiv = document.getElementById('fixed-progress');
		}

		// Função para atualizar a barra de progresso fixa (só executa se logado)
		function updateFixedProgress(seasonIndex) {
			if (!isLogged) return;
			const season = seasons[seasonIndex];
			if (!season) return;
			const watchedEpisodes = getWatchedEpisodes();
			const total = season.episodes.length;
			const watchedCount = season.episodes.filter(ep => watchedEpisodes.includes(ep.id)).length;
			const percent = total === 0 ? 0 : (watchedCount / total) * 100;
			const progressText = `${watchedCount}/${total} ${translations.watchedEpisodes}`;

			fixedProgressDiv.innerHTML = `
				<div class="season-progress-text">${progressText}</div>
				<div class="progress-bar-bg">
					<div class="progress-bar-fill" style="width: ${percent}%;"></div>
				</div>
			`;
		}

		function renderEpisodesForSeason(seasonIndex) {
			const season = seasons[seasonIndex];
			const container = document.getElementById('episodes-scrollable');
			if (!container) return;

			if (!season.episodes.length) {
				container.innerHTML = '<p>' + translations.noEpisodesInSeason + '</p>';
				updateFixedProgress(seasonIndex);
				return;
			}

			const isLogged = isUserLoggedIn();
			const watchedEpisodes = isLogged ? getWatchedEpisodes() : [];

			let episodesHtml = '<ul class="episode-list">';
			season.episodes.forEach(ep => {
				const releaseDate = ep.released ? new Date(ep.released) : null;
				const isFuture = releaseDate && releaseDate > new Date();
				let mediaHtml = '';

				const calendarSvg = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>`;

				if (isFuture) {
					mediaHtml = `<div class="ep-placeholder">${calendarSvg} ${translations.soon}</div>`;
				} else if (ep.thumbnail) {
					mediaHtml = `<img src="${ep.thumbnail}" alt="${ep.name}" class="ep-thumb check-thumb">`;
				} else {
					mediaHtml = `<div class="ep-placeholder">${calendarSvg} ${translations.soon}</div>`;
				}

				let checkButton = '';
				if (isLogged && !isFuture) {
					const isWatched = watchedEpisodes.includes(ep.id);
					const checkedSvg = `<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>`;
					const uncheckedSvg = `<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/></svg>`;
					checkButton = `<button class="watched-check" data-ep-id="${ep.id}">${isWatched ? checkedSvg : uncheckedSvg}</button>`;
				}

				// Adiciona atributos data-season e data-episode para clique
				episodesHtml += `
					<li data-season="${ep.season}" data-episode="${ep.episode}" class="episode-item-li">
						<div class="episode-item">
							${mediaHtml}
							<div class="ep-info">
								<strong>${ep.episode}. ${ep.name}</strong>
								${ep.released ? `<br><small>${new Date(ep.released).toLocaleDateString()}</small>` : ''}
							</div>
							${checkButton}
						</div>
					</li>
				`;
			});
			episodesHtml += '</ul>';
			container.innerHTML = episodesHtml;

			// Validação de imagens
			const images = container.querySelectorAll('img.check-thumb');
			images.forEach(img => {
				function validateImage() {
					if (img.naturalWidth === 0) {
						const placeholder = document.createElement('div');
						placeholder.className = 'ep-placeholder';
						placeholder.textContent = translations.soon;
						img.parentNode.replaceChild(placeholder, img);
					} else {
						img.classList.remove('check-thumb');
					}
				}
				if (img.complete) {
					validateImage();
				} else {
					img.onload = validateImage;
					img.onerror = () => {
						const placeholder = document.createElement('div');
						placeholder.className = 'ep-placeholder';
						placeholder.textContent = translations.soon;
						img.parentNode.replaceChild(placeholder, img);
					};
				}
			});

			// Evento de clique no episódio (para abrir no Stremio)
			const episodeItems = container.querySelectorAll('.episode-item-li');
			episodeItems.forEach(item => {
				item.addEventListener('click', (e) => {
					// Evita que o clique no botão de check dispare a abertura
					if (e.target.closest('.watched-check')) return;
					const season = item.getAttribute('data-season');
					const episode = item.getAttribute('data-episode');
					openEpisodeInStremio(season, episode);
				});
			});

			// Eventos para botões de assistido
			if (isLogged) {
				const checkButtons = container.querySelectorAll('.watched-check');
				checkButtons.forEach(btn => {
					btn.addEventListener('click', (e) => {
						e.stopPropagation();
						const epId = btn.dataset.epId;
						const newState = toggleWatchedEpisode(epId);
						const checkedSvg = `<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>`;
						const uncheckedSvg = `<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/></svg>`;
						btn.innerHTML = newState ? checkedSvg : uncheckedSvg;
						updateFixedProgress(seasonIndex);
					});
				});
			}

			// Atualiza a barra de progresso fixa para esta temporada (se logado)
			updateFixedProgress(seasonIndex);
		}

		selectEl.addEventListener('change', (e) => {
			currentSeasonIndex = parseInt(e.target.value, 10);
			renderEpisodesForSeason(currentSeasonIndex);
		});
		renderEpisodesForSeason(currentSeasonIndex);
	}
});