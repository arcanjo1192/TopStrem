// ==================== Listas ====================
async function fetchUserLists(type) {
    const url = type ? `/api/lists?type=${type}` : '/api/lists';
    const headers = await getAuthHeader();
    headers['Accept'] = 'application/json';
    const res = await fetch(url, { headers, credentials: 'same-origin' });
    if (!res.ok) return [];
    const data = await res.json();
    return data.lists || [];
}

async function createList(name, type) {
    const headers = await getAuthHeader();
    headers['Content-Type'] = 'application/json';
    const res = await fetch('/api/lists', {
        method: 'POST',
        headers,
        credentials: 'same-origin',
        body: JSON.stringify({ action: 'create', listName: name, listType: type })
    });
    return res.ok;
}

async function deleteList(name) {
    const headers = await getAuthHeader();
    headers['Content-Type'] = 'application/json';
    const res = await fetch('/api/lists', {
        method: 'POST',
        headers,
        credentials: 'same-origin',
        body: JSON.stringify({ action: 'delete', listName: name })
    });
    return res.ok;
}

async function addItemToList(listName, item) {
    const headers = await getAuthHeader();
    headers['Content-Type'] = 'application/json';
    const res = await fetch('/api/lists', {
        method: 'POST',
        headers,
        credentials: 'same-origin',
        body: JSON.stringify({ action: 'add_item', listName, item })
    });
    return res.ok;
}

// ==================== Modal personalizado ====================
function showCustomModal(options) {
    const {
        title = '',
        message = '',
        type = 'alert',
        defaultValue = '',
        confirmText = 'OK',
        cancelText = 'Cancelar'
    } = options;

    return new Promise((resolve) => {
        const overlay = document.createElement('div');
        overlay.className = 'custom-modal-overlay';

        let contentHTML = `
            <div class="custom-modal-content">
                <div class="custom-modal-title">${title}</div>
                <div class="custom-modal-message">${message}</div>
        `;

        if (type === 'prompt') {
            contentHTML += `<input class="custom-modal-input" type="text" value="${defaultValue}" placeholder="" />`;
        }

        contentHTML += `<div class="custom-modal-buttons">`;

        if (type === 'confirm' || type === 'prompt') {
            contentHTML += `<button class="custom-modal-btn secondary cancel-btn">${cancelText}</button>`;
        }
        contentHTML += `<button class="custom-modal-btn primary confirm-btn">${confirmText}</button>`;
        contentHTML += `</div></div>`;

        overlay.innerHTML = contentHTML;
        document.body.appendChild(overlay);

        const input = overlay.querySelector('.custom-modal-input');
        const confirmBtn = overlay.querySelector('.confirm-btn');
        const cancelBtn = overlay.querySelector('.cancel-btn');

        const closeModal = (result) => {
            document.body.removeChild(overlay);
            resolve(result);
        };

        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) {
                closeModal({ confirmed: false, value: '' });
            }
        });

        confirmBtn.addEventListener('click', () => {
            if (type === 'prompt') {
                closeModal({ confirmed: true, value: input.value.trim() });
            } else {
                closeModal({ confirmed: true });
            }
        });

        if (cancelBtn) {
            cancelBtn.addEventListener('click', () => {
                closeModal({ confirmed: false, value: '' });
            });
        }

        if (input) {
            input.focus();
            input.addEventListener('keydown', (e) => {
                if (e.key === 'Enter') {
                    confirmBtn.click();
                }
            });
        }
    });
}

// ==================== Validação de limite ====================
async function createListWithLimit(name, type) {
    const lists = await fetchUserLists(type);
    if (lists.length >= 5) {
        await showCustomModal({
            title: 'Limite Atingido',
            message: `Você já possui 5 listas de ${type === 'movie' ? 'filmes' : 'séries'}. Remova uma para criar outra.`,
            type: 'alert'
        });
        return false;
    }
    return await createList(name, type);
}

// ==================== Modal de seleção de lista (adicionar item) ====================
async function showListSelectionPopup(item) {
    const lists = await fetchUserLists(item.type);
    if (lists.length === 0) {
        const result = await showCustomModal({
            title: 'Sem listas',
            message: `Você ainda não tem listas para ${item.type === 'movie' ? 'filmes' : 'séries'}. Deseja criar uma agora?`,
            type: 'confirm',
            confirmText: 'Criar',
            cancelText: 'Cancelar'
        });
        if (result.confirmed) {
            const nameResult = await showCustomModal({
                title: 'Nova Lista',
                message: 'Digite o nome da nova lista:',
                type: 'prompt',
                confirmText: 'Criar',
                cancelText: 'Cancelar',
                defaultValue: ''
            });
            if (nameResult.confirmed && nameResult.value) {
                const ok = await createListWithLimit(nameResult.value, item.type);
                if (ok) {
                    await showCustomModal({
                        title: 'Sucesso',
                        message: 'Lista criada! Agora escolha a lista para adicionar o item.',
                        type: 'alert'
                    });
                    showListSelectionPopup(item);
                } else {
                    await showCustomModal({
                        title: 'Erro',
                        message: 'Não foi possível criar a lista.',
                        type: 'alert'
                    });
                }
            }
        }
        return;
    }

    let html = '<div class="lists-popup"><h3>Selecionar lista</h3><ul>';
    lists.forEach(list => {
        html += `
            <li style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                <button class="list-option" data-listname="${list.name}" style="flex: 1; text-align: left;">
                    ${list.name} (${list.items.length} itens)
                </button>
                <button class="delete-list-btn" data-listname="${list.name}" style="background: none; border: none; color: #e50914; cursor: pointer; margin-left: 10px;" title="Deletar lista">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        <line x1="10" y1="11" x2="10" y2="17"/>
                        <line x1="14" y1="11" x2="14" y2="17"/>
                    </svg>
                </button>
            </li>
        `;
    });
    html += '</ul><button id="new-list-btn">+ Nova lista</button></div>';

    const popup = document.createElement('div');
    popup.className = 'list-selection-modal';
    popup.innerHTML = html;
    document.body.appendChild(popup);

    popup.addEventListener('click', function(e) {
        if (e.target === popup) {
            document.body.removeChild(popup);
        }
    });

    popup.querySelectorAll('.list-option').forEach(btn => {
        btn.addEventListener('click', async () => {
            const listName = btn.dataset.listname;
            const added = await addItemToList(listName, item);
            const msg = added ? 'Item adicionado com sucesso!' : 'Erro ou item já existente.';
            await showCustomModal({
                title: added ? 'Sucesso' : 'Erro',
                message: msg,
                type: 'alert'
            });
            document.body.removeChild(popup);
        });
    });

    popup.querySelectorAll('.delete-list-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            e.stopPropagation();
            const listName = btn.dataset.listname;
            const result = await showCustomModal({
                title: 'Deletar lista',
                message: `Tem certeza que deseja deletar a lista "${listName}"?`,
                type: 'confirm',
                confirmText: 'Deletar',
                cancelText: 'Cancelar'
            });
            if (result.confirmed) {
                const deleted = await deleteList(listName);
                if (deleted) {
                    await showCustomModal({
                        title: 'Deletada',
                        message: 'Lista deletada com sucesso.',
                        type: 'alert'
                    });
                    document.body.removeChild(popup);
                    showListSelectionPopup(item);
                } else {
                    await showCustomModal({
                        title: 'Erro',
                        message: 'Erro ao deletar a lista.',
                        type: 'alert'
                    });
                }
            }
        });
    });

    popup.querySelector('#new-list-btn').addEventListener('click', async () => {
        const result = await showCustomModal({
            title: 'Nova Lista',
            message: 'Digite o nome da nova lista:',
            type: 'prompt',
            confirmText: 'Criar',
            cancelText: 'Cancelar',
            defaultValue: ''
        });
        if (result.confirmed && result.value) {
            const ok = await createListWithLimit(result.value, item.type);
            if (ok) {
                document.body.removeChild(popup);
                showListSelectionPopup(item);
            } else {
                await showCustomModal({
                    title: 'Erro',
                    message: 'Não foi possível criar a lista.',
                    type: 'alert'
                });
            }
        }
    });
}

// Configurar botão "Adicionar à lista" na página de detalhes
document.addEventListener('DOMContentLoaded', () => {
    const addBtn = document.getElementById('add-to-list-btn');
    if (addBtn) {
        addBtn.addEventListener('click', () => {
            const item = {
                id: addBtn.dataset.id,
                type: addBtn.dataset.type,
                name: addBtn.dataset.name,
                year: addBtn.dataset.year
            };
            showListSelectionPopup(item);
        });
    }
});

// ==================== Modal "Todas as suas listas" ====================
async function showAllListsModal() {
    const lists = await fetchUserLists();
    if (lists.length === 0) {
        await showCustomModal({
            title: 'Sem listas',
            message: 'Você ainda não possui listas.',
            type: 'alert'
        });
        return;
    }
	
	const shareToken = await getShareToken();

    const modal = document.createElement('div');
    modal.className = 'list-selection-modal';
    modal.innerHTML = `
        <h3>Todas as suas listas</h3>
        <ul>
            ${lists.map(list => `
                <li class="all-lists-item" data-type="${list.type}" data-ids="${list.items.map(i => i.id).join(',')}">
                    <div class="list-name">${list.name}</div>
                    <div class="list-meta">
                        <span>${list.type === 'movie' ? 'Filmes' : 'Séries'} · ${list.items.length} itens</span>
                        <button class="delete-list-btn" data-listname="${list.name}">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <polyline points="3 6 5 6 21 6"/>
                                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                                <line x1="10" y1="11" x2="10" y2="17"/>
                                <line x1="14" y1="11" x2="14" y2="17"/>
                            </svg>
                            Deletar
                        </button>
                    </div>
                </li>
            `).join('')}
        </ul>
    `;
    document.body.appendChild(modal);

    modal.addEventListener('click', function(e) {
        if (e.target === modal) {
            document.body.removeChild(modal);
        }
    });

	modal.querySelectorAll('.all-lists-item').forEach(item => {
		item.addEventListener('click', function(e) {
			if (e.target.closest('.delete-list-btn')) return;
			const type = item.dataset.type;
			const ids = item.dataset.ids;
			const listName = item.querySelector('.list-name').textContent.trim();
			let url = `/lists?type=${type}&ids=${ids}&list=${encodeURIComponent(listName)}`;
			if (shareToken) {
				url += `&user=${encodeURIComponent(shareToken)}`;
			}
			window.location.href = url;
		});
	});

    modal.querySelectorAll('.delete-list-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            e.stopPropagation();
            const listName = btn.dataset.listname;
            const result = await showCustomModal({
                title: 'Deletar lista',
                message: `Tem certeza que deseja deletar a lista "${listName}"?`,
                type: 'confirm',
                confirmText: 'Deletar',
                cancelText: 'Cancelar'
            });
            if (result.confirmed) {
                const deleted = await deleteList(listName);
                if (deleted) {
                    document.body.removeChild(modal);
                    showAllListsModal();
                } else {
                    await showCustomModal({
                        title: 'Erro',
                        message: 'Erro ao deletar a lista.',
                        type: 'alert'
                    });
                }
            }
        });
    });
}

// ==================== Remover item da lista na página de listas ====================
async function removeItemFromCurrentList(itemId) {
    const grid = document.getElementById('catalog-grid');
    if (!grid) return;
    const listName = grid.dataset.listName;
    if (!listName) {
        await showCustomModal({
            title: 'Erro',
            message: 'Não foi possível identificar a lista.',
            type: 'alert'
        });
        return;
    }

    const result = await showCustomModal({
        title: 'Remover item',
        message: 'Tem certeza que deseja remover este item da lista?',
        type: 'confirm',
        confirmText: 'Remover',
        cancelText: 'Cancelar'
    });
    if (!result.confirmed) return;

    const headers = await getAuthHeader();
    headers['Content-Type'] = 'application/json';
    const res = await fetch('/api/lists', {
        method: 'POST',
        headers,
        credentials: 'same-origin',
        body: JSON.stringify({ action: 'remove_item', listName, item: { id: itemId } })
    });

    if (res.ok) {
        // Remove o card do DOM
        const card = document.querySelector(`.movie-card[data-movie-id="${itemId}"]`);
        if (card) card.remove();

        // Atualiza a URL sem recarregar
        const url = new URL(window.location);
        let ids = url.searchParams.get('ids');
        if (ids) {
            ids = ids.split(',').map(id => id.trim()).filter(id => id !== itemId).join(',');
            if (ids) {
                url.searchParams.set('ids', ids);
            } else {
                url.searchParams.delete('ids');
            }
            window.history.replaceState({}, '', url);
        }

        if (grid.querySelectorAll('.movie-card').length === 0) {
            window.location.href = '/';
        }
    } else {
        await showCustomModal({
            title: 'Erro',
            message: 'Não foi possível remover o item da lista.',
            type: 'alert'
        });
    }
}

// Delegate para botões de remover (criados dinamicamente)
document.addEventListener('DOMContentLoaded', () => {
    document.body.addEventListener('click', (e) => {
        const btn = e.target.closest('.remove-from-list-btn');
        if (btn) {
            const itemId = btn.dataset.id;
            removeItemFromCurrentList(itemId);
        }
    });
});

async function getShareToken() {
    if (!isUserLoggedIn()) return '';
    try {
        const headers = await getAuthHeader();
        headers['Accept'] = 'application/json';
        const res = await fetch('/api/share-token', { headers, credentials: 'same-origin' });
        if (res.ok) {
            const data = await res.json();
            return data.token || '';
        }
    } catch(e) {}
    return '';
}