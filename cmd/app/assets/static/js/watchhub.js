document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById('watchhub-modal');
    const contentDiv = document.getElementById('watchhub-modal-content');
    const closeBtn = document.querySelector('.watchhub-modal-close');
    const watchBtn = document.getElementById('watchhub-btn');

    if (!watchBtn) return;

    // Usa as traduções já definidas globalmente (do objeto translations)
    const watchTrans = {
        loading: translations.loading,
        noStreams: translations.noStreams,
        error: translations.error,
        available: translations.available,
        watchWeb: translations.watchWeb,
        watchAndroid: translations.watchAndroid,
        watchIOS: translations.watchIOS
    };

    function getPlatform() {
        const ua = navigator.userAgent;
        if (/android/i.test(ua)) return 'android';
        if (/iPad|iPhone|iPod/.test(ua)) return 'ios';
        return 'web';
    }

    watchBtn.addEventListener('click', function() {
        const mediaType = this.getAttribute('data-type');
        const id = this.getAttribute('data-id');
        contentDiv.innerHTML = `<div class="loading">${watchTrans.loading}</div>`;
        modal.classList.add('show');

        fetch(`/api/watch/${mediaType}/${id}`)
            .then(response => {
                if (!response.ok) throw new Error(`HTTP ${response.status}`);
                return response.json();
            })
            .then(data => {
                const platform = getPlatform();
                const streams = data.streams || [];
                const filtered = streams.filter(s => {
                    if (platform === 'web') return s.externalUrl;
                    if (platform === 'android') return s.androidUrl;
                    if (platform === 'ios') return s.iosUrl;
                    return false;
                });

                if (filtered.length === 0) {
                    contentDiv.innerHTML = `<p>${watchTrans.noStreams}</p>`;
                    return;
                }

                let html = `<h3>${watchTrans.available}</h3><ul class="stream-list">`;
                for (const stream of filtered) {
                    let url = '';
                    let label = '';
                    if (platform === 'web') {
                        url = stream.externalUrl;
                        label = watchTrans.watchWeb;
                    } else if (platform === 'android') {
                        url = stream.androidUrl;
                        label = watchTrans.watchAndroid;
                    } else if (platform === 'ios') {
                        url = stream.iosUrl;
                        label = watchTrans.watchIOS;
                    }
                    html += `<li><a href="${url}" target="_blank" rel="noopener noreferrer" class="stream-btn">${label} (${stream.name})</a></li>`;
                }
                html += '</ul>';
                contentDiv.innerHTML = html;
            })
            .catch(err => {
                contentDiv.innerHTML = `<div class="error">${watchTrans.error} Detalhe: ${err.message}</div>`;
                console.error(err);
            });
    });

    if (closeBtn) {
        closeBtn.onclick = () => modal.classList.remove('show');
        window.onclick = (e) => { if (e.target === modal) modal.classList.remove('show'); };
    }
});