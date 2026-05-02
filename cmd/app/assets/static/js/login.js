// ==================== Autenticação Híbrida: Web + Mobile Nativo ====================

// Estado global do usuário
let currentUser = {
    email: null,
    name: null,
    isLoggedIn: false,
    token: null // para app nativo
};

let userFavorites = [];

// Detectar tipo de client
let clientType = {
    isNativeApp: window.isMobileApp === true || !!window.getAuthToken,
    isWebBrowser: true
};

// ========== Funções para App Nativo (com token) ==========

// Obter token armazenado (implementado no app nativo)
function getStoredToken() {
    // Para app nativo: implementar no lado nativo
    // return await secureStorage.getToken()
    
    // Fallback: verificar se foi passado via window
    return window.authToken || null;
}

// Armazenar token (implementado no app nativo)
async function storeToken(token) {
    // Para app nativo: implementar no lado nativo
    // await secureStorage.setToken(token)
    
    // Fallback: guardar em window
    window.authToken = token;
}

// Obter token para enviar em requisições (cookie ou header)
async function getAuthHeader() {
    if (clientType.isNativeApp) {
        const token = await getStoredToken();
        if (token) {
            return { 'Authorization': 'Bearer ' + token };
        }
    }
    // Para web: browser envia cookie automaticamente
    return {};
}

// ========== Funções Comuns para Web + Mobile ==========

// Função para buscar os dados do usuário
async function fetchCurrentUser() {
    try {
        const headers = await getAuthHeader();
        const response = await fetch('/api/me', {
            headers: headers,
            credentials: 'same-origin' // inclui o cookie HttpOnly (para web)
        });
        
        if (response.ok) {
            const user = await response.json();
            currentUser.email = user.email;
            currentUser.name = user.name;
            currentUser.isLoggedIn = true;
            await syncFavorites();
        } else {
            currentUser.email = null;
            currentUser.name = null;
            currentUser.isLoggedIn = false;
            userFavorites = [];
            if (clientType.isNativeApp) {
                await storeToken(null); // limpar token inválido
            }
        }
        updateUserUI();
        return currentUser.isLoggedIn;
    } catch (err) {
        console.error('Erro ao verificar autenticação:', err);
        currentUser.isLoggedIn = false;
        updateUserUI();
        return false;
    }
}

// Funções auxiliares
function getCurrentUserEmail() {
    return currentUser.email;
}

function isUserLoggedIn() {
    return currentUser.isLoggedIn;
}

// Atualiza a interface de acordo com o estado do usuário
function updateUserUI() {
    const userDropdown = document.getElementById('user-dropdown');
    const loginLink = document.getElementById('login-link');
    const notificationsBtn = document.getElementById('notifications-btn');
    const userNameSpan = document.getElementById('user-name-display');
    const userInitialSpan = document.getElementById('user-initial');
	const myListsBtn = document.getElementById('my-lists-btn');
	const addToListBtn = document.getElementById('add-to-list-btn');

    const isLoggedIn = currentUser.isLoggedIn;
    const isNativeApp = clientType.isNativeApp;

    if (isLoggedIn) {
        if (userDropdown) userDropdown.style.display = 'inline-block';
        if (loginLink) loginLink.style.display = 'none';
        if (notificationsBtn) notificationsBtn.style.display = 'inline-block';
        if (userNameSpan) userNameSpan.textContent = currentUser.name;
        if (userInitialSpan) userInitialSpan.textContent = currentUser.name.charAt(0).toUpperCase();
        attachFavoritesListeners();
		if (myListsBtn) myListsBtn.addEventListener('click', showAllListsModal);
		if (addToListBtn) addToListBtn.style.display = 'inline-flex';
    } else {
        if (userDropdown) userDropdown.style.display = 'none';
        if (loginLink) loginLink.style.display = isNativeApp ? 'none' : 'inline-block';
        if (notificationsBtn) notificationsBtn.style.display = 'none';
		if (addToListBtn) addToListBtn.style.display = 'none';
    }

    const watchlistBtns = document.querySelectorAll('.watchlist-btn, .watchlist-detail-btn');
    watchlistBtns.forEach(btn => {
        btn.style.display = isLoggedIn ? 'inline-flex' : 'none';
    });

    initWatchlistButtons();
}

// Logout
async function logout() {
    try {
        const headers = await getAuthHeader();
        await fetch('/auth/logout', { 
            method: 'POST',
            headers: headers,
            credentials: 'same-origin' 
        });
        
        currentUser.isLoggedIn = false;
        currentUser.email = null;
        currentUser.name = null;
        
        if (clientType.isNativeApp) {
            await storeToken(null); // limpar token
        }
        
        updateUserUI();
        window.location.reload();
    } catch (err) {
        console.error('Erro no logout:', err);
    }
}

// ========== Fluxo de Login (Web + Mobile) ==========

// Callback quando autenticação é bem-sucedida (app nativo chama após interceptar deepLink)
window.onAuthTokenReceived = function(token) {
    currentUser.token = token;
    storeToken(token);
    fetchCurrentUser();
    window.location.href = '/';
};

// Iniciar login
async function startLogin() {
    // Caso 1: Web browser tradicional
    if (!clientType.isNativeApp) {
        // Redireciona a página inteira – o backend então redireciona para o Google
        window.location.href = '/auth/login';
        return;
    }

    // Caso 2: App nativo – faz fetch para obter a URL de autenticação
    try {
        const response = await fetch('/auth/login', {
            headers: { 'X-Client-Type': 'native' }
        });
        if (!response.ok) throw new Error('Erro ao iniciar login');
        const data = await response.json();

        if (data.authUrl) {
            if (window.openExternalBrowser) {
                window.openExternalBrowser(data.authUrl);
            } else {
                window.open(data.authUrl, '_blank');
            }
        }
    } catch (err) {
        console.error('Erro ao iniciar login:', err);
        alert('Erro ao autenticar. Tente novamente.');
    }
}

// Setup de listeners
function setupMobileLogin() {
    const loginBtn = document.getElementById('login-link');
    if (loginBtn) {
        loginBtn.addEventListener('click', function(e) {
            e.preventDefault();
            startLogin(); // usa a função unificada
        });
    }
}

// ==================== Watchlist (server-side favorites quando logado) ====================
function getWatchlistKey() {
    const email = getCurrentUserEmail();
    if (!email) return null;
    return `topstrem_watchlist_${email}`;
}

async function fetchFavorites() {
    if (!isUserLoggedIn()) {
        return [];
    }

    try {
        const headers = await getAuthHeader();
        const response = await fetch('/api/favorites', {
            method: 'GET',
            headers,
            credentials: 'same-origin'
        });
        if (!response.ok) {
            return [];
        }
        const data = await response.json();
        return Array.isArray(data.favorites) ? data.favorites : [];
    } catch (err) {
        console.error('Erro ao buscar favoritos do servidor:', err);
        return [];
    }
}

async function syncFavorites() {
    userFavorites = await fetchFavorites();
}

function getWatchlist() {
    if (isUserLoggedIn()) {
        return userFavorites;
    }

    const key = getWatchlistKey();
    if (!key) return [];
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : [];
}

function saveWatchlist(list) {
    if (isUserLoggedIn()) {
        userFavorites = list;
        return;
    }

    const key = getWatchlistKey();
    if (!key) return;
    localStorage.setItem(key, JSON.stringify(list));
}

function isInWatchlist(id) {
    const list = getWatchlist();
    return list.some(item => item.id === id);
}

async function updateFavoritesOnServer(item, action) {
    try {
        const headers = await getAuthHeader();
        headers['Content-Type'] = 'application/json';
        const response = await fetch('/api/favorites', {
            method: 'POST',
            headers,
            credentials: 'same-origin',
            body: JSON.stringify({ action, item })
        });
        if (!response.ok) {
            return false;
        }
        const data = await response.json();
        userFavorites = Array.isArray(data.favorites) ? data.favorites : userFavorites;
        return true;
    } catch (err) {
        console.error('Erro ao atualizar favoritos no servidor:', err);
        return false;
    }
}

async function addToWatchlist(item) {
    if (isUserLoggedIn()) {
        const added = await updateFavoritesOnServer(item, 'add');
        if (added && !userFavorites.some(i => i.id === item.id)) {
            userFavorites.push(item);
        }
        return added;
    }

    const list = getWatchlist();
    if (!list.some(i => i.id === item.id)) {
        list.push(item);
        saveWatchlist(list);
        return true;
    }
    return false;
}

async function removeFromWatchlist(id) {
    if (isUserLoggedIn()) {
        const removed = await updateFavoritesOnServer({ id }, 'remove');
        if (removed) {
            userFavorites = userFavorites.filter(item => item.id !== id);
        }
        return removed;
    }

    let list = getWatchlist();
    const newList = list.filter(item => item.id !== id);
    saveWatchlist(newList);
    return newList.length !== list.length;
}

async function toggleWatchlist(item) {
    if (isInWatchlist(item.id)) {
        return !(await removeFromWatchlist(item.id));
    } else {
        return await addToWatchlist(item);
    }
}

function updateWatchlistButton(btn, saved) {
    if (saved) {
        btn.classList.add('saved');
    } else {
        btn.classList.remove('saved');
    }
}

// ==================== Episódios assistidos (servidor + localStorage fallback) ====================
let watchedEpisodesCache = null;

function getWatchedEpisodesKey() {
    const email = getCurrentUserEmail();
    if (!email) return null;
    return `topstrem_watched_episodes_${email}`;
}

async function fetchWatchedEpisodesFromServer() {
    if (!isUserLoggedIn()) {
        return [];
    }

    if (watchedEpisodesCache !== null) {
        return watchedEpisodesCache;
    }

    try {
        const headers = await getAuthHeader();
        const response = await fetch('/api/watched-episodes', {
            method: 'GET',
            headers,
            credentials: 'same-origin'
        });
        if (!response.ok) {
            return [];
        }
        const data = await response.json();
        watchedEpisodesCache = Array.isArray(data.watchedEpisodes) ? data.watchedEpisodes : [];
        return watchedEpisodesCache;
    } catch (err) {
        console.error('Erro ao buscar episódios assistidos do servidor:', err);
        return [];
    }
}

async function getWatchedEpisodes() {
    if (isUserLoggedIn()) {
        if (watchedEpisodesCache !== null) {
            return watchedEpisodesCache;
        }
        return await fetchWatchedEpisodesFromServer();
    }

    const key = getWatchedEpisodesKey();
    if (!key) return [];
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : [];
}

function saveWatchedEpisodes(episodes) {
    if (isUserLoggedIn()) {
        watchedEpisodesCache = episodes || [];
        return;
    }

    const key = getWatchedEpisodesKey();
    if (!key) return;
    localStorage.setItem(key, JSON.stringify(episodes));
}

async function toggleWatchedEpisode(episodeId) {
    if (isUserLoggedIn()) {
        const watched = watchedEpisodesCache !== null ? watchedEpisodesCache : await fetchWatchedEpisodesFromServer();
        const isWatched = watched.includes(episodeId);
        const action = isWatched ? 'remove' : 'add';
        try {
            const headers = await getAuthHeader();
            headers['Content-Type'] = 'application/json';
            const response = await fetch('/api/watched-episodes', {
                method: 'POST',
                headers,
                credentials: 'same-origin',
                body: JSON.stringify({ action, episodeId })
            });
            if (!response.ok) {
                return isWatched;
            }
            const data = await response.json();
            watchedEpisodesCache = Array.isArray(data.watchedEpisodes) ? data.watchedEpisodes : watchedEpisodesCache;
            return watchedEpisodesCache.includes(episodeId);
        } catch (err) {
            console.error('Erro ao atualizar episódio assistido no servidor:', err);
            return isWatched;
        }
    }

    let watched = await getWatchedEpisodes();
    if (watched.includes(episodeId)) {
        watched = watched.filter(id => id !== episodeId);
    } else {
        watched.push(episodeId);
    }
    saveWatchedEpisodes(watched);
    return watched.includes(episodeId);
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

    const watchedEpisodes = await getWatchedEpisodes();

    const allEpisodes = [];
    for (const series of seriesList) {
        const episodes = await fetchSeriesEpisodes(series.id, series.name);
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
            <div class="notification-item" data-series-id="${ep.seriesId}">
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

    const items = contentDiv.querySelectorAll('.notification-item');
    items.forEach(item => {
        item.style.cursor = 'pointer';
        item.addEventListener('click', function(e) {
            const seriesId = this.getAttribute('data-series-id');
            if (seriesId) {
                window.location.href = `/detail/series/${seriesId}`;
            }
        });
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
    const saved = await toggleWatchlist(item);
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
            const watchlist = getWatchlist();
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
            const watchlist = getWatchlist();
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
document.addEventListener('DOMContentLoaded', async function() {
    // Primeiro, obtém o estado do usuário via cookie
    await fetchCurrentUser();

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
    setupMobileLogin();
});

document.addEventListener('htmx:afterSettle', function() {
    initWatchlistButtons();
    updateUserUI();
});