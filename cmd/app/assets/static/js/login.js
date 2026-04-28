// ==================== JWT / Autenticação ====================
function parseJWT(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
        return JSON.parse(jsonPayload);
    } catch (e) {
        return null;
    }
}

function getCurrentUserEmail() {
    const token = localStorage.getItem('auth_token');
    if (!token) return null;
    const payload = parseJWT(token);
    return payload && payload.email ? payload.email : null;
}

function updateUserUI() {
    const token = localStorage.getItem('auth_token');
    const userDropdown = document.getElementById('user-dropdown');
    const loginLink = document.getElementById('login-link');
    const notificationsBtn = document.getElementById('notifications-btn');
    const userNameSpan = document.getElementById('user-name-display');
    const userInitialSpan = document.getElementById('user-initial');

    const isLoggedIn = token && parseJWT(token) && parseJWT(token).name;
	
	const isMobileApp = window.isMobileApp === true;

    if (isLoggedIn) {
        const payload = parseJWT(token);
        if (userDropdown) userDropdown.style.display = 'inline-block';
        if (loginLink) loginLink.style.display = 'none';
        if (notificationsBtn) notificationsBtn.style.display = 'inline-block';
        if (userNameSpan) userNameSpan.textContent = payload.name;
        if (userInitialSpan) userInitialSpan.textContent = payload.name.charAt(0).toUpperCase();
        attachFavoritesListeners();
    } else {
        if (userDropdown) userDropdown.style.display = 'none';
        if (loginLink) loginLink.style.display = isMobileApp ? 'none' : 'inline-block';
        if (notificationsBtn) notificationsBtn.style.display = 'none';
        if (token) localStorage.removeItem('auth_token');
    }

    const watchlistBtns = document.querySelectorAll('.watchlist-btn, .watchlist-detail-btn');
    watchlistBtns.forEach(btn => {
        btn.style.display = isLoggedIn ? 'inline-flex' : 'none';
    });

    initWatchlistButtons();
}

function logout() {
    localStorage.removeItem('auth_token');
    updateUserUI();
    window.location.reload();
}

// ==================== Watchlist ====================
function getWatchlistKey() {
    const email = getCurrentUserEmail();
    if (!email) return null;
    return `topstrem_watchlist_${email}`;
}

function getWatchlist() {
    const key = getWatchlistKey();
    if (!key) return [];
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : [];
}

function saveWatchlist(list) {
    const key = getWatchlistKey();
    if (!key) return;
    localStorage.setItem(key, JSON.stringify(list));
}

function isInWatchlist(id) {
    const list = getWatchlist();
    return list.some(item => item.id === id);
}

function addToWatchlist(item) {
    const list = getWatchlist();
    if (!list.some(i => i.id === item.id)) {
        list.push(item);
        saveWatchlist(list);
        return true;
    }
    return false;
}

function removeFromWatchlist(id) {
    let list = getWatchlist();
    const newList = list.filter(item => item.id !== id);
    saveWatchlist(newList);
    return newList.length !== list.length;
}

function toggleWatchlist(item) {
    if (isInWatchlist(item.id)) {
        removeFromWatchlist(item.id);
        return false;
    } else {
        addToWatchlist(item);
        return true;
    }
}

function isUserLoggedIn() {
    return !!getCurrentUserEmail();
}

function updateWatchlistButton(btn, saved) {
    if (saved) {
        btn.classList.add('saved');
    } else {
        btn.classList.remove('saved');
    }
}

// ==================== Episódios assistidos (para filtro de notificações) ====================
function getWatchedEpisodesKey() {
    const email = getCurrentUserEmail();
    if (!email) return null;
    return `topstrem_watched_episodes_${email}`;
}

function getWatchedEpisodes() {
    const key = getWatchedEpisodesKey();
    if (!key) return [];
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : [];
}

// ==================== Notificações ====================
async function fetchSeriesEpisodes(seriesId, seriesName) {
    const url = `https://v3-cinemeta.strem.io/meta/series/${seriesId}.json`;
    try {
        const response = await fetch(url);
        if (!response.ok) return [];
        const data = await response.json();
        const videos = data.meta?.videos || [];
        const seriesLogo = data.meta?.logo || '';
        const today = new Date();
        const startDate = new Date(today);
        startDate.setDate(today.getDate() - 6);
        const endDate = new Date(today);
        endDate.setDate(today.getDate() + 7);

        return videos.filter(video => {
            if (!video.released) return false;
            const releaseDate = new Date(video.released);
            return releaseDate >= startDate && releaseDate <= endDate;
        }).map(video => ({
            id: video.id,
            name: video.name,
            episode: video.episode,
            season: video.season,
            released: video.released,
            thumbnail: video.thumbnail || seriesLogo,
            seriesId: seriesId,
            seriesName: seriesName,
            seriesLogo: seriesLogo
        }));
    } catch (err) {
        console.error(`Erro ao buscar série ${seriesId}:`, err);
        return [];
    }
}

let notificationsLoaded = false;
let cachedNotifications = [];

async function loadNotifications() {
    if (!isUserLoggedIn()) return [];

    const watchlist = getWatchlist();
    const seriesList = watchlist.filter(item => item.type === 'series');
    if (seriesList.length === 0) return [];

    const watchedEpisodes = getWatchedEpisodes(); // IDs dos episódios já assistidos

    const allEpisodes = [];
    for (const series of seriesList) {
        const episodes = await fetchSeriesEpisodes(series.id, series.name);
        // Filtra apenas os episódios que NÃO estão assistidos
        const unwatchedEpisodes = episodes.filter(ep => !watchedEpisodes.includes(ep.id));
        if (unwatchedEpisodes.length > 0) {
            allEpisodes.push(...unwatchedEpisodes);
        }
    }
    allEpisodes.sort((a, b) => new Date(a.released) - new Date(b.released));
    return allEpisodes;
}

function renderNotifications(episodes) {
    const contentDiv = document.getElementById('notifications-content');
    if (!contentDiv) return;

    if (episodes.length === 0) {
        contentDiv.innerHTML = '<p>Nenhuma notificação no momento.</p>';
        updateBadge(0);
        return;
    }

    let html = '';
    for (const ep of episodes) {
        const releaseDate = ep.released ? new Date(ep.released) : null;
        const formattedDate = releaseDate ? releaseDate.toLocaleDateString() : '';
        const imgSrc = ep.thumbnail || ep.seriesLogo;
        html += `
            <div class="notification-item">
                <div style="display: flex; gap: 12px; align-items: flex-start;">
                    <img class="notif-thumb" src="${imgSrc}" alt="${ep.name}" data-fallback="${ep.seriesLogo}"
                         style="width: 80px; height: auto; border-radius: 4px; object-fit: cover;">
                    <div>
                        <strong>${ep.seriesName}</strong><br>
                        <span>T${ep.season} E${ep.episode}: ${ep.name}</span><br>
                        <small>${formattedDate}</small>
                    </div>
                </div>
            </div>
        `;
    }
    contentDiv.innerHTML = html;

    const images = contentDiv.querySelectorAll('.notif-thumb');
    images.forEach(img => {
        function validateImage() {
            if (img.naturalWidth === 0) {
                const fallback = img.getAttribute('data-fallback');
                if (fallback) {
                    img.src = fallback;
                    img.removeAttribute('data-fallback');
                } else {
                    img.outerHTML = '<div class="ep-placeholder">📷</div>';
                }
            }
        }
        if (img.complete) {
            validateImage();
        } else {
            img.onload = validateImage;
            img.onerror = () => {
                const fallback = img.getAttribute('data-fallback');
                if (fallback) {
                    img.src = fallback;
                    img.removeAttribute('data-fallback');
                } else {
                    img.outerHTML = '<div class="ep-placeholder">📷</div>';
                }
            };
        }
    });

    updateBadge(episodes.length);
}

function updateBadge(count) {
    const badgeSpan = document.getElementById('notification-badge');
    const notificationsBtn = document.getElementById('notifications-btn');
    if (!badgeSpan) return;
    if (count > 0) {
        badgeSpan.textContent = count;
        badgeSpan.style.display = 'flex';
        if (notificationsBtn) notificationsBtn.style.color = '#e50914';
    } else {
        badgeSpan.style.display = 'none';
        if (notificationsBtn) notificationsBtn.style.color = '#fff';
    }
}

// ==================== Atualização de notificações após alteração na watchlist ====================
async function refreshNotifications() {
    if (!isUserLoggedIn()) return;
    const episodes = await loadNotifications();
    cachedNotifications = episodes;
    updateBadge(episodes.length);
    const sidebar = document.getElementById('notifications-sidebar');
    if (sidebar && sidebar.classList.contains('open')) {
        renderNotifications(episodes);
    }
}

async function handleWatchlistClick(e) {
    e.preventDefault();
    e.stopPropagation();
    const btn = e.currentTarget;
    const id = btn.dataset.id;
    const type = btn.dataset.type;
    const name = btn.dataset.name;
    const year = btn.dataset.year;

    if (!isUserLoggedIn()) {
        if (confirm('Você precisa estar logado para salvar itens. Deseja fazer login agora?')) {
            window.location.href = '/auth/login';
        }
        return;
    }

    const item = { id, type, name, year };
    const saved = toggleWatchlist(item);
    updateWatchlistButton(btn, saved);
    await refreshNotifications();
}

function initWatchlistButtons() {
    const btns = document.querySelectorAll('.watchlist-btn, .watchlist-detail-btn');
    btns.forEach(btn => {
        const id = btn.dataset.id;
        let saved = false;
        if (isUserLoggedIn()) {
            saved = isInWatchlist(id);
        }
        updateWatchlistButton(btn, saved);
        btn.removeEventListener('click', handleWatchlistClick);
        btn.addEventListener('click', handleWatchlistClick);
    });
}

// ==================== Favoritos (menu dropdown) ====================
function attachFavoritesListeners() {
    const favSeriesBtn = document.getElementById('favorites-series');
    const favMovieBtn = document.getElementById('favorites-movie');

    if (favSeriesBtn && !favSeriesBtn.hasAttribute('data-listener')) {
        favSeriesBtn.addEventListener('click', () => {
            const email = getCurrentUserEmail();
            if (!email) return;
            const key = `topstrem_watchlist_${email}`;
            const stored = localStorage.getItem(key);
            const watchlist = stored ? JSON.parse(stored) : [];
            const seriesIds = watchlist.filter(item => item.type === 'series').map(item => item.id);
            if (seriesIds.length === 0) {
                alert('Nenhuma série favorita encontrada.');
                return;
            }
            const idsParam = seriesIds.map(encodeURIComponent).join(',');
            window.location.href = `/favorites?type=series&ids=${idsParam}`;
        });
        favSeriesBtn.setAttribute('data-listener', 'true');
    }

    if (favMovieBtn && !favMovieBtn.hasAttribute('data-listener')) {
        favMovieBtn.addEventListener('click', () => {
            const email = getCurrentUserEmail();
            if (!email) return;
            const key = `topstrem_watchlist_${email}`;
            const stored = localStorage.getItem(key);
            const watchlist = stored ? JSON.parse(stored) : [];
            const movieIds = watchlist.filter(item => item.type === 'movie').map(item => item.id);
            if (movieIds.length === 0) {
                alert('Nenhum filme favorito encontrado.');
                return;
            }
            const idsParam = movieIds.map(encodeURIComponent).join(',');
            window.location.href = `/favorites?type=movie&ids=${idsParam}`;
        });
        favMovieBtn.setAttribute('data-listener', 'true');
    }
}

// ==================== Inicialização ====================
document.addEventListener('DOMContentLoaded', function() {
    updateUserUI();
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) logoutBtn.addEventListener('click', logout);

    const notificationsBtn = document.getElementById('notifications-btn');
    const sidebar = document.getElementById('notifications-sidebar');
    const closeBtn = document.getElementById('close-notifications');
    const contentDiv = document.getElementById('notifications-content');

    if (notificationsBtn && sidebar && contentDiv) {
        notificationsBtn.addEventListener('click', async function() {
            sidebar.classList.add('open');
            if (!notificationsLoaded) {
                contentDiv.innerHTML = '<div class="loading">Carregando notificações...</div>';
                const episodes = await loadNotifications();
                cachedNotifications = episodes;
                renderNotifications(episodes);
                notificationsLoaded = true;
            } else {
                renderNotifications(cachedNotifications);
            }
        });
        if (closeBtn) {
            closeBtn.onclick = () => sidebar.classList.remove('open');
            window.onclick = (e) => { if (e.target === sidebar) sidebar.classList.remove('open'); };
        }
    }

    if (isUserLoggedIn()) {
        loadNotifications().then(episodes => {
            cachedNotifications = episodes;
            updateBadge(episodes.length);
        });
    }

    initWatchlistButtons();
    attachFavoritesListeners();
});

document.addEventListener('htmx:afterSettle', function() {
    initWatchlistButtons();
    updateUserUI();
});